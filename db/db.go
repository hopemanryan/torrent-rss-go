package localDB

import (
	"fmt"
	"log"
	"strings"
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
	SeasonAndEpisode string
}

func NewDb() *DB {
	s, _ := scribble.New("./db/data", nil)

	newDb := DB{
		storage: s,
	}

	return &newDb
}

func (db *DB) CheckDownloadedName(name string, seasonAndEpisode string) bool {
	file := FileItem{}
	parsedFileName := strings.ReplaceAll(name, " ", ".")

	db.storage.Read(tableName, fmt.Sprintf("%s_%s", parsedFileName, seasonAndEpisode), &file)

	return file.Name != ""

}

func (db *DB) SaveFile(filename string, origianlFileName string, seasonAndEpisode string) {
	parsedFileName := strings.ReplaceAll(filename, " ", ".")
	newName := fmt.Sprintf("%s_%s", parsedFileName, seasonAndEpisode)
	db.storage.Write(tableName, newName, FileItem{
		Name:             newName,
		OrigianlFileName: origianlFileName,
		Date:             time.Now(),
	})
	log.Printf("Saving new file: %s", filename)

}
