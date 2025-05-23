package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func scan1(dir string) {
	db, err := sql.Open("sqlite", "./files.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS file_info (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT,
		"size" INTEGER,
		"mode" TEXT,
		"mod_time" DATETIME,
		"path" TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	var files []FileInfo
	var wg sync.WaitGroup

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			print(err.Error())
		}
		if !info.IsDir() {
			files = append(files, FileInfo{info.Name(), info.Size(), info.Mode(), info.ModTime().Format("2006-01-02 15:04:05"), path, ""})
			if len(files) == 1000 {
				wg.Add(1)
				go insertFiles(db, files, &wg)
				files = nil // reset slice
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(files) > 0 { // Insert remaining files
		wg.Add(1)
		go insertFiles(db, files, &wg)
	}

	wg.Wait() // Wait for all goroutines to finish
}

// CalculateSelectiveMD5 calculates the MD5 hash of the first 1%, middle, and last 1% of the file in a streaming manner.
func CalculateSelectiveMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Get file size
	info, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSize := info.Size()

	// Calculate the size of the parts to read (1% of the file)
	partSize := fileSize / 1000
	if partSize == 0 {
		partSize = 1 // Ensure at least 1 byte is read if file is too small
	}

	hash := md5.New()

	// Function to hash a part
	hashPart := func(offset int64, size int64) error {
		if _, err := file.Seek(offset, io.SeekStart); err != nil {
			return err
		}

		// Create a limited reader to read only the partSize or remaining bytes
		limitedReader := io.LimitReader(file, size)
		if _, err := io.Copy(hash, limitedReader); err != nil {
			return err
		}
		print("!")
		return nil
	}

	// Hash the first part
	if err := hashPart(0, partSize); err != nil {
		return "", err
	}

	// Hash the middle part
	middleOffset := fileSize/2 - partSize/2
	if middleOffset+partSize > fileSize {
		middleOffset = fileSize - partSize // Adjust if calculation exceeds file size
	}
	if err := hashPart(middleOffset, partSize); err != nil {
		return "", err
	}

	// Hash the last part
	lastOffset := fileSize - partSize
	if lastOffset < 0 {
		lastOffset = 0 // Ensure positive offset
	}
	if err := hashPart(lastOffset, partSize); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
