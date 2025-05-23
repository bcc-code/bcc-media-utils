package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
)

func insertFiles(db *sql.DB, files []FileInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO file_info(name, size, mode, mod_time, path) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, file := range files {
		_, err := stmt.Exec(file.Name, file.Size, file.Mode.String(), file.ModTime, file.Path)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	print(".")
}

func hashForQuery(offset int64) {
	dbPath := "./files.db"
	queryX := `SELECT path
FROM file_info fi
         INNER JOIN (
    SELECT size, count(*) as cnt
    FROM file_info WHERE size > 1024*1024*1024
                   AND name NOT LIKE "%%.wav"
    GROUP BY size
    HAVING COUNT(*) > 1
) b ON fi.size = b.size
WHERE b.cnt > 2
ORDER by fi.size
LIMIT 10 OFFSET %d;`

	queryX = fmt.Sprintf(queryX, offset)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	rows, err := db.Query(queryX)
	if err != nil {
		log.Fatal("Error executing query:", err)
	}
	defer rows.Close()

	hashes := make(map[string]string)

	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			log.Fatal("Error scanning row:", err)
		}

		md5Hash, err := CalculateSelectiveMD5(path)
		if err != nil {
			log.Printf("Error calculating MD5 for file %s: %v\n", path, err)
			continue
		}

		hashes[path] = md5Hash
		print(".")
	}

	rows.Close()

	for path, md5Hash := range hashes {
		if err := UpdateFileMD5(db, path, md5Hash); err != nil {
			log.Printf("Error updating MD5 for file %s: %v\n", path, err)
		}
	}

}

func UpdateFileMD5(db *sql.DB, path, md5Hash string) error {
	_, err := db.Exec("UPDATE file_info SET md5_hash = ? WHERE path = ?", md5Hash, path)
	return err
}
