package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BlogPost struct {
	PostId string	`bson:"_id"`
	Title string	`bson:"title"`
	Techs string	`bson:"techs"`
	Status string	`bson:"status"`
	Clickme string	`bson:"clickme"`
	Tryme string	`bson:"tryme,omitempty"`
	Readme string	`bson:"readme,omitempty"`
}

func CreatePost(client *mongo.Client, dbName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		coll := client.Database(dbName).Collection("posts")

		title := r.FormValue("title")
		techs := r.FormValue("techs")
		status := r.FormValue("status")
		clickme := r.FormValue("clickme")
		tryme := r.FormValue("tryme")
		readme := r.FormValue("readme")

		blogPost := &BlogPost{
			Title: title,
			Techs: techs,
			Status: status,
			Clickme: clickme,
			Tryme: tryme,
			Readme: readme,
		}

		result, err := coll.InsertOne(context.TODO(), blogPost)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Document Inserted with success, _id:%v", result.InsertedID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}