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

var dbUserName = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
var dbUserPassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
var mongoUrl = fmt.Sprintf("mongodb://%s:%s@mongo", dbUserName, dbUserPassword)

type DB struct {
	Storage  *mongo.Client
	episodes *mongo.Collection
	CTX      context.Context
}

func NewDb() *DB {

	if dbUserName == "" {
		mongoUrl = "mongodb://127.0.0.1:27017"
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUrl))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	if err != nil {
		log.Fatal(err)
	}

	torrentDb := createDb(client)
	episodesCollection := createCollection(torrentDb, "episodes")

	newDb := DB{
		Storage:  client,
		episodes: episodesCollection,
		CTX:      ctx,
	}

	return &newDb
}

func createCollection(db *mongo.Database, collectionName string) *mongo.Collection {
	return db.Collection(collectionName)
}
func createDb(client *mongo.Client) *mongo.Database {
	return client.Database("torrent")
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

func (db *DB) SaveFile(origianlFileName string, seasonAndEpisode string) {
	newName := fmt.Sprintf("%s_%s", origianlFileName, seasonAndEpisode)
	_, err := db.episodes.InsertOne(context.TODO(), bson.M{
		"fileName": newName,
		"Date":     time.Now(),
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Saving new file: %s", newName)

}
