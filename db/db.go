package localDB

import (
	"time"

	scribble "github.com/nanobox-io/golang-scribble"
)

var tableName = "downloadedItems"

type DB struct {
	storage *scribble.Driver
}

type FileItem struct {
	Name             string
	OrigianlFileName string
	FileMoved        bool
	Date             time.Time
}

func NewDb() *DB {
	s, _ := scribble.New("./db/data", nil)

	newDb := DB{
		storage: s,
	}

	return &newDb
}

func (db *DB) CheckDownloadedName(name string) bool {
	file := FileItem{}

	db.storage.Read(tableName, name, &file)

	return file.Name != ""

}

func (db *DB) SaveFile(filename string, origianlFileName string) {
	db.storage.Write(tableName, filename, FileItem{
		Name:             filename,
		OrigianlFileName: origianlFileName,
		Date:             time.Now(),
	})
}
