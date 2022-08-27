package localDB

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var tableName = "downloadedItems"
var dbUserName = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
var dbUserPassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
var dbName = os.Getenv("MONGO_INITDB_DATABASE")
var mongoUrl = fmt.Sprintf("mongodb://%s:%s@127.0.0.1:27017/%s", dbUserName, dbUserPassword, dbName)

type DB struct {
	Storage  *mongo.Client
	episodes *mongo.Collection
	CTX      context.Context
}

type FileItem struct {
	Name             string
	OrigianlFileName string
	FileMoved        bool
	Date             time.Time
	SeasonAndEpisode string
}

func NewDb() *DB {

	if dbUserName == "" {
		mongoUrl = "mongodb://127.0.0.1:27017"
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	if err != nil {
		log.Fatal(err)
	}

	quickstartDatabase := client.Database("torrents")
	episodesCollection := quickstartDatabase.Collection("episodes")

	newDb := DB{
		Storage:  client,
		episodes: episodesCollection,
		CTX:      ctx,
	}

	return &newDb
}

func (db *DB) CheckDownloadedName(name string, seasonAndEpisode string) bool {

	parsedFileName := strings.ReplaceAll(name, " ", ".")
	newName := fmt.Sprintf("%s_%s", parsedFileName, seasonAndEpisode)

	var episode bson.M
	var err = db.episodes.FindOne(context.TODO(), bson.M{"fileName": newName}).Decode(&episode)

	if err != nil {
		return false
	}
	fmt.Println(episode)

	return episode != nil
}

func (db *DB) SaveFile(filename string, origianlFileName string, seasonAndEpisode string) {
	parsedFileName := strings.ReplaceAll(filename, " ", ".")
	newName := fmt.Sprintf("%s_%s", parsedFileName, seasonAndEpisode)
	_, err := db.episodes.InsertOne(context.TODO(), bson.M{
		"fileName":         newName,
		"OrigianlFileName": origianlFileName,
		"Date":             time.Now(),
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Saving new file: %s", filename)

}

// todo test run
