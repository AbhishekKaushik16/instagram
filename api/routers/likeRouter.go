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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getLikesById(w http.ResponseWriter, r *http.Request) {
	pathVariables := mux.Vars(r)
	id := pathVariables["id"]
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	filter := bson.M{"_id": objectId}
	var like models.Likes
	err = database.Collection("likes").FindOne(context.Background(), filter).Decode(&like)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Like not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	likeJSON, err := json.Marshal(like)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(likeJSON)
}

func getLikesOnParent(w http.ResponseWriter, r *http.Request) {
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
	var like models.Likes

	err = database.Collection("likes").FindOne(context.Background(), filter).Decode(&like)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(like)
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

func addLikes(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a User struct
	var like models.Likes
	if err := json.NewDecoder(r.Body).Decode(&like); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	parentId := like.ParentID
	parentType := like.ParentType
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"parentId": parentId, "parentType": parentType}
	update := bson.M{"$inc": bson.M{"likes": 1}}
	// Increment the like of the parentId and parentType
	result, err := database.Collection("likes").UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		log.Printf("Error updating likes: %s\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	response, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error updating likes: %s\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	// Set the response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func removeLikes(w http.ResponseWriter, r *http.Request) {
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

func AddLikesRouter(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/likes").Subrouter()
	s.HandleFunc("/{parentId}", getLikesOnParent).Methods("GET")
	s.HandleFunc("/{id}", getLikesById).Methods("GET")
	s.HandleFunc("", addLikes).Methods("POST")
	s.HandleFunc("/{id}", removeLikes).Methods("POST")

	return s
}
