package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/siddharthPriyadarshi/mongoapi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
)

// const connectionString = MONGO_DB_CON_STR
const dbName = "netflix"
const colName = "watchlist"

// MOST IMPORTANT

var collection *mongo.Collection

// connect with mongoDB
// first time when application execute
// run once

func init() {

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var connectionString = os.Getenv("MONGO_DB_CON_STR")
	fmt.Println("Connection String: ", connectionString)

	//	client options
	clientOptions := options.Client().ApplyURI(connectionString)

	//	connect to mongodb
	//whenever u req with req with machines
	//outside your machine then you need to make sure
	//upto what duration, deadline of connection
	//TODO() is type of context
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		//panic
		log.Fatal(err)
	}

	fmt.Println("MONGODB connection success")

	collection = client.Database(dbName).Collection(colName)

	//	collection instance/referecnce
	fmt.Println("Collection instance is ready!")
}

// Mongo Helpers - file
// insert one record
func insertOneMovie(movie model.Netflix) {
	fmt.Println(movie)

	inserted, err := collection.InsertOne(context.Background(), movie)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Movie inserted with id: ", inserted.InsertedID)
}

// update 1 record
func updateOneMovie(movieId string) {
	id, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Modified count: ", result.ModifiedCount)
}

func deleteOneMovie(movieId string) {
	id, err := primitive.ObjectIDFromHex(movieId)

	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"_id": id}

	delete, err := collection.DeleteOne(context.Background(), filter, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Movie Deleted Count: ", delete.DeletedCount)
}

// delete all records from mongo Db
func deleteAllMovies() int64 {
	//{} indicates to select all record in db
	filter := bson.D{{}}
	deleteRes, err := collection.DeleteMany(context.Background(), filter, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("All movies deleted count: ", deleteRes.DeletedCount)
	return deleteRes.DeletedCount
}

// get all movies from database
func getAllMovies() []primitive.M {
	//return a cursor which is type of mongo db object
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var movies []primitive.M
	for cur.Next(context.Background()) {
		var movie bson.M
		err := cur.Decode(&movie)
		if err != nil {
			log.Fatal(err)

		}
		movies = append(movies, movie)
	}
	defer cur.Close(context.Background())
	return movies
}

//Actual controllers - file

func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie model.Netflix
	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)
}

func MarkedAsWatcher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}

func DeleteOneMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)

	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	count := deleteAllMovies()

	json.NewEncoder(w).Encode(count)
}
