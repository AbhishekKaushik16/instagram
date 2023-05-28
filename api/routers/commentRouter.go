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

func getCommentById(w http.ResponseWriter, r *http.Request) {
	pathVariables := mux.Vars(r)
	id := pathVariables["id"]
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	filter := bson.M{"_id": objectId}
	var comment models.Comments
	err = database.Collection("comments").FindOne(context.Background(), filter).Decode(&comment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Comment not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	commentJSON, err := json.Marshal(comment)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(commentJSON)
}

func getCommentsOnParent(w http.ResponseWriter, r *http.Request) {
	// Retrieve all users from the MongoDB collection
	pathVariables := mux.Vars(r)
	id := pathVariables["parentId"]
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	filter := bson.M{"parenId": objectId}
	cursor, err := database.Collection("comments").Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var comments []models.Comments
	for cursor.Next(context.Background()) {
		var comment models.Comments
		if err := cursor.Decode(&comment); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	response, err := json.Marshal(comments)
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

func createComment(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a User struct
	var comment models.Comments
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Insert the new user into the MongoDB collection
	result, err := database.Collection("comments").InsertOne(context.Background(), comment)
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

func deleteComment(w http.ResponseWriter, r *http.Request) {
	pathVariables := mux.Vars(r)
	userId := pathVariables["id"]
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	filter := bson.M{"_id": objectId}
	result, err := database.Collection("comments").DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("Error deleting Comment: %s\n", err)
		http.Error(w, "Error deleting Comment", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		log.Printf("Comment not found with ID: %s\n", userId)
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	// Convert the result to JSON
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Comment deleted successfully",
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

func AddCommentsRouter(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/comments").Subrouter()
	s.HandleFunc("/{parentId}", getCommentsOnParent).Methods("GET")
	s.HandleFunc("/{id}", getCommentById).Methods("GET")
	s.HandleFunc("", createComment).Methods("POST")
	s.HandleFunc("/{id}", deleteComment).Methods("DELETE")

	return s
}
