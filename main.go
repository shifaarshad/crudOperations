package main

import (
	"context"
	"encoding/json"

	//"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"

	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define a struct for the data you want to store in MongoDB
type User struct {
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}

var user1 = User{"ali", "ali.ali@yahoo.com", "abcdef"}
var user2 = User{"arooba", "arooba.ali@yahoo.com", "ahhjjh"}

// MongoDB Atlas configuration
const connectionString = "mongodb+srv://shifahajiarshad:hajiarshad99@basicapi.nkkkgjl.mongodb.net/?retryWrites=true&w=majority"
const dbName = "mongogorilla"
const collectionName = "gorilla"

// Create a MongoDB client and connect to the database
var client *mongo.Client

// it is run only for the first time
func init() {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb+srv://shifahajiarshad:hajiarshad99@basicapi.nkkkgjl.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

}

// Create a new user in the MongoDB database
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//its is struct variable
	//var user User

	//decode the values of user struct and than after decoding ,
	// you can send it to the encoder to encode in mongodb
	_ = json.NewDecoder(r.Body).Decode(&user1)
	collection := client.Database("mongogorilla").Collection("gorilla")
	// bson.M unordered BSOn document (MAP)

	result, err := collection.InsertOne(context.Background(), user1)
	if err != nil {
		log.Fatal(err)
	}
	//return the response to the user.
	json.NewEncoder(w).Encode(result)
}

// Get a list of all users in the MongoDB database
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []User
	collection := client.Database("mongogorilla").Collection("gorilla")

	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	//return the response to the user.
	json.NewEncoder(w).Encode(users)
}

// Get a single user by Name from the MongoDB database
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// grab id from the request, vars(all variable vlaues will be extratcted from the request which is being passed by the user)
	//mux.Vars actually map of  key value pair
	params := mux.Vars(r)
	localName := params["name"]

	collection := client.Database("mongogorilla").Collection("gorilla")
	err := collection.FindOne(context.Background(), bson.M{"name": localName}).Decode(&user1)
	if err != nil {
		log.Fatal(err)
	}

	//return the response to the user.
	json.NewEncoder(w).Encode(user1)
}

// Update a user by Name in the MongoDB database
func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User

	//parse request body into a person struct
	_ = json.NewDecoder(r.Body).Decode(&user)

	collection := client.Database("mongogorilla").Collection("gorilla")

	filter := bson.M{"name": "ali"}
	update := bson.M{"$set": bson.M{
		"name":     "shifa",
		"email":    "assdfhh@yahoo.com",
		"password": "assdfg",
	}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	//return the response to the user.
	json.NewEncoder(w).Encode(result)
}

func putsomedata(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	collection := client.Database("mongogorilla").Collection("gorilla")

	// Create filter and update document
	filter := bson.M{"name": "ali"}
	update := bson.M{"$set": bson.M{
		"email": "aghanoor@yahoo.com",
	}}

	// Perform update operation
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user was found and updated
	if result.MatchedCount == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Write updated user to response

	json.NewEncoder(w).Encode(result)
}

func main() {

	// create router with NewRouter method and assign it to instance r
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.go run main.go

	//after declaring a new router instance, you can use the HandleFunc method of your router
	//instance to assign routes to handler functions along with the request type that the
	//handler function handles.

	r.HandleFunc("/create", createUser).Methods("POST")
	r.HandleFunc("/getuser/{name}", getUser).Methods("GET")
	r.HandleFunc("/getusers", getUsers).Methods("GET")
	r.HandleFunc("/update/{name}", updateUser).Methods("PUT")
	r.HandleFunc("/patch/{name}", putsomedata).Methods("PATCH")

	//http.ListenAndServe() function to start the server and tell it to listen for
	//new HTTP requests and
	//then serve them using the handler functions you set up
	//You can set up a server using the ListenAndServe method of the http package.
	//The ListenAndServe method takes as arguments the port you want the server to run on
	//and a router instance
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":7000", r))
}
