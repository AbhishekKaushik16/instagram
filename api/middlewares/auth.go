package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/AbhishekKaushik16/instagram/api/db"
	"github.com/AbhishekKaushik16/instagram/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticated := false
		headers := struct {
			userId string
		}{
			userId: r.Header.Get("user-id"),
		}
		if headers.userId != "" {
			_, database, err := db.SetupDbClient()
			fmt.Println("Authenticating request")
			fmt.Println(headers.userId)
			if err != nil {
				log.Fatal("Error setting up Databse")
			}
			objectId, err := primitive.ObjectIDFromHex(headers.userId)
			var user models.User
			database.Collection("user").FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&user)
			fmt.Println(user)
			if user.ID != headers.userId {
				authenticated = false
			} else {
				authenticated = true
			}
		} else {
			fmt.Println("userId not set in headers")
		}
		if !authenticated {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
