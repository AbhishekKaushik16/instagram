package routers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/AbhishekKaushik16/instagram/api/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getUserById(w http.ResponseWriter, r *http.Request) {
	pathVariables := mux.Vars(r)
	id := pathVariables["id"]
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	filter := bson.M{"_id": objectId}
	var user models.Users
	err = database.Collection("user").FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	userJSON, err := json.Marshal(user)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	// Retrieve all users from the MongoDB collection
	cursor, err := database.Collection("user").Find(context.Background(), bson.M{})
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var users []models.Users
	for cursor.Next(context.Background()) {
		var user models.Users
		if err := cursor.Decode(&user); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Serialize the users slice to JSON
	response, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a User struct
	var user models.Users
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Insert the new user into the MongoDB collection
	result, err := database.Collection("user").InsertOne(context.Background(), user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Serialize the inserted document ID to JSON
	response, err := json.Marshal(result.InsertedID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	pathVariables := mux.Vars(r)
	userId := pathVariables["id"]
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	filter := bson.M{"_id": objectId}
	result, err := database.Collection("user").DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("Error deleting user: %s\n", err)
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		log.Printf("User not found with ID: %s\n", userId)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Convert the result to JSON
	response := struct {
		Message string `json:"message"`
	}{
		Message: "User deleted successfully",
	}

	// Convert the response struct to JSON byte slice
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling JSON: %s\n", err)
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the HTTP response writer
	w.Write(responseJSON)
}

func AddUsersRouter(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/users").Subrouter()
	s.HandleFunc("", getUsers).Methods("GET")
	s.HandleFunc("/{id}", getUserById).Methods("GET")
	s.HandleFunc("", createUser).Methods("POST")
	s.HandleFunc("/{id}", deleteUser).Methods("DELETE")
	s.HandleFunc("", getUsers).Methods("GET")

	return s
}
