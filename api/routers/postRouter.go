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

func getPostById(w http.ResponseWriter, r *http.Request) {
	pathVariables := mux.Vars(r)
	id := pathVariables["id"]
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	filter := bson.M{"_id": objectId}
	var post models.Posts
	err = database.Collection("posts").FindOne(context.Background(), filter).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	userJSON, err := json.Marshal(post)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	// Retrieve all users from the MongoDB collection
	cursor, err := database.Collection("posts").Find(context.Background(), bson.M{})
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var posts []models.Posts
	for cursor.Next(context.Background()) {
		var post models.Posts
		if err := cursor.Decode(&post); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	response, err := json.Marshal(posts)
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

func createPost(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a User struct
	var post models.Posts
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Insert the new user into the MongoDB collection

	result, err := database.Collection("posts").InsertOne(context.Background(), post)
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

func deletePost(w http.ResponseWriter, r *http.Request) {
	pathVariables := mux.Vars(r)
	userId := pathVariables["id"]
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	filter := bson.M{"_id": objectId}
	result, err := database.Collection("posts").DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("Error deleting post: %s\n", err)
		http.Error(w, "Error deleting post", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		log.Printf("Post not found with ID: %s\n", userId)
		http.Error(w, "Post not found", http.StatusNotFound)
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

func AddPostsRouter(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/posts").Subrouter()
	s.HandleFunc("", getPosts).Methods("GET")
	s.HandleFunc("/{id}", getPostById).Methods("GET")
	s.HandleFunc("", createPost).Methods("POST")
	s.HandleFunc("/{id}", deleteUser).Methods("DELETE")
	s.HandleFunc("", deletePost).Methods("GET")

	return s
}
