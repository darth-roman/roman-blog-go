package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BlogPost struct {
	BlogID string			`bson:"blog_id"`
	Title string			`bson:"title"`
	Techs string			`bson:"techs"`
	Status string			`bson:"status"`
	Clickme string			`bson:"clickme"`
	Tryme string			`bson:"tryme,omitempty"`
	Readme string			`bson:"readme,omitempty"`
	CreatedAt time.Time		`bson:"create_at"`
	UpdatedAt time.Time		`bson:"updated_at"`
}

type BlogPostResponse struct {
	ID any		    `bson:"_id"`	
	Title string	`bson:"title"`
}

var coll *mongo.Collection

func CreatePost(client *mongo.Client, dbName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		coll = client.Database(dbName).Collection("posts")

		title := r.FormValue("title")
		techs := r.FormValue("techs")
		status := r.FormValue("status")
		clickme := r.FormValue("clickme")
		tryme := r.FormValue("tryme")
		readme := r.FormValue("readme")

		blogId, _ := uuid.NewRandom()

		fmt.Println(blogId.String())

		blogPost := &BlogPost{
			BlogID: blogId.String(),
			Title: title,
			Techs: techs,
			Status: status,
			Clickme: clickme,
			Tryme: tryme,
			Readme: readme,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		result, err := coll.InsertOne(context.TODO(), blogPost)
		if err != nil {
			log.Fatal(err)
		}
		
		blogPostResponse := &BlogPostResponse{
			ID: result.InsertedID ,
			Title: blogPost.Title,
		}
		returnResponseJSON(w, http.StatusCreated, blogPostResponse)
	}
}

func GetOnePostByUUID(client *mongo.Client, dbName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		coll = client.Database(dbName).Collection("posts")

		uuid := strings.TrimPrefix(r.URL.Path, "/blogs/")

		if uuid == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		filter := bson.D{{"blog_id", uuid}}
		var result bson.M

		err := coll.FindOne(context.TODO(), filter).Decode(&result)
		
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				http.Error(w, "document not found", http.StatusNotFound)
				return
			}
			log.Printf("failed to find document: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		returnResponseJSON(w, http.StatusOK, result)
	}
}