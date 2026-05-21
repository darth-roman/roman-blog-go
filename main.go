package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoClient struct {
	Client *mongo.Client
	Err error
}

var mongoClient MongoClient

func connectToDatabase(dbName string) *MongoClient{
	database_uri := "mongodb://localhost:27017/"+dbName
	if dbName == "" {
		log.Fatal("No Database provided")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(database_uri))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to %s\n", dbName)

	return &MongoClient{Client: client, Err: err}

}

func main(){
	godotenv.Load()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
    	// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
    	AllowedOrigins:   []string{"https://*", "http://*"},
    	// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
  }))
	dabataseName := os.Getenv("DATABASE_NAME")
	mongoClient = *connectToDatabase(dabataseName)
	r.Route("/blogs", func(r chi.Router){
		r.Post("/create", CreatePost(mongoClient.Client, dabataseName))
		r.Get("/{id}", GetOnePostByUUID(mongoClient.Client, dabataseName))
	})
	
	var port string = fmt.Sprintf(":%s", os.Getenv("PORT"))
	fmt.Println(port)
	http.ListenAndServe(port, r)

	// defer func(){
	// 	if err := mongoClient.Client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()


}