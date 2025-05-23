package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime string
	Path    string
	VXID    string
}

const VS_USER = "<REPLACE>"
const VS_PASS = "<REPLACE>"

func findVXbyFilePath(storage, file string) *FileSearchResult {
	file = url.PathEscape(file)
	url := fmt.Sprintf("http://10.12.128.15:8080/API/storage/%s/file?includeItem=true&includeShapes=true&path=%s", storage, file)
	username := VS_USER
	password := VS_PASS

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	req.Header.Add("Accept", "application/json")
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	req.Header.Add("Authorization", "Basic "+auth)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}
	result := &FileSearchResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}
	return result
}

type DupedFileRow struct {
	ID      int
	Md5Hash string
	Path    string
	Size    int
	Cnt     int
	Total   int
	FileVX  string
	AssetVX string
}

func getDupesListForPath(path string) []DupedFileRow {
	db, err := sql.Open("sqlite3", "./files.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	likePath := path + "%" // Append wildcard for LIKE query

	query := `
SELECT fi.id, fi.size, fi.path, fi.md5_hash, h.cnt, h.total
FROM file_info fi
JOIN (
    SELECT md5_hash, count(*) as cnt, sum(size) as total
    FROM file_info
    WHERE md5_hash IS NOT NULL
    AND path NOT LIKE '%temp%'
    GROUP BY md5_hash
    HAVING cnt > 1
) AS h ON h.md5_hash = fi.md5_hash AND fi.path NOT LIKE '%temp%'
WHERE fi.path LIKE ?
ORDER BY h.total;`

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(likePath)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var files []DupedFileRow

	for rows.Next() {
		var f DupedFileRow
		err = rows.Scan(&f.ID, &f.Size, &f.Path, &f.Md5Hash, &f.Cnt, &f.Total)
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, f)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return files
}

func UpdateVXIDS(db *sql.DB, md5Hash, path, fileVX, assetVX string) (int64, error) {
	query := `
UPDATE file_info
SET fileVX = ?, assetVX = ?
WHERE md5_hash = ? AND path = ?
`

	// Prepare the query for execution
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err // Return the error to the caller
	}
	defer stmt.Close()

	// Execute the query with the provided parameters
	res, err := stmt.Exec(fileVX, assetVX, md5Hash, path)
	if err != nil {
		return 0, err // Return the error to the caller
	}

	// Check how many rows were affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err // Return the error to the caller
	}

	return rowsAffected, nil // Return the number of rows affected and no error
}

// Scan files on a file system and connect them to vidispine IDs
// Store the results in a database

func main() {
	prefix := "/mnt/isilon/system/tempingest/"
	files := getDupesListForPath(prefix)

	for _, f := range files {
		path, _ := strings.CutPrefix(f.Path, prefix)
		fsRes := findVXbyFilePath("VX-47", path)
		if fsRes.Hits == 0 {
			continue
		}

		if fsRes.Hits > 1 {
			continue
		}

		f.FileVX = fsRes.File[0].ID

		if len(fsRes.File[0].Item) > 0 {
			f.AssetVX = fsRes.File[0].Item[0].ID
		}
		print(".")
	}

	db, err := sql.Open("sqlite3", "./files.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	for _, f := range files {
		cnt, err := UpdateVXIDS(db, f.Md5Hash, f.Path, f.FileVX, f.AssetVX)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(cnt)
	}

}

type FileSearchResult struct {
	Hits int    `json:"hits,omitempty"`
	File []File `json:"file,omitempty"`
}

type Component struct {
	ID string `json:"id,omitempty"`
}

type Shape struct {
	ID        string      `json:"id,omitempty"`
	Component []Component `json:"component,omitempty"`
}

type Item struct {
	ID    string  `json:"id,omitempty"`
	Shape []Shape `json:"shape,omitempty"`
}

type Field struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type Metadata struct {
	Field []Field `json:"field,omitempty"`
}

type File struct {
	ID          string   `json:"id,omitempty"`
	Path        string   `json:"path,omitempty"`
	URI         []string `json:"uri,omitempty"`
	State       string   `json:"state,omitempty"`
	Size        int      `json:"size,omitempty"`
	Hash        string   `json:"hash,omitempty"`
	Timestamp   string   `json:"timestamp,omitempty"`
	RefreshFlag int      `json:"refreshFlag,omitempty"`
	Storage     string   `json:"storage,omitempty"`
	Item        []Item   `json:"item,omitempty"`
	Metadata    Metadata `json:"metadata,omitempty"`
}
