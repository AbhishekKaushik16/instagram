package routers

import (
	"net/http"

	"github.com/AbhishekKaushik16/instagram/api/db"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	client   *mongo.Client
	database *mongo.Database
	err      error
)

func SetupDbClient() {
	// Set up MongoDB connection
	client, err = db.GetDBClient()
	if err != nil {
		log.Fatal(err)
	}
	// Access the "users" collection
	database = client.Database("users")
}

func Routers() *mux.Router {
	//StrictSlash defines the trailing slash behavior for new routes. The initial value is false.
	//When true, if the route path is "/path/", accessing "/path" will perform a redirect to the former and vice versa.
	SetupDbClient()
	r := mux.NewRouter().StrictSlash(true)

	//PathPrefix /api adds a matcher for the URL path prefix.
	s := r.PathPrefix("/api").Subrouter()
	AddUsersRouter(s)
	AddPostsRouter(s)
	AddCommentsRouter(s)
	AddLikesRouter(s)
	r.Use(loggingMiddleware)
	return r
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
