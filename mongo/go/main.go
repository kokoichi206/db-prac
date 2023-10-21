package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	Title  string   `bson:"title"`
	Author string   `bson:"author"`
	Pages  int      `bson:"pages"`
	Genres []string `bson:"genres"`
	Rating float64  `bson:"rating"`
}

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	colection := client.Database("bookstore").Collection("books")

	curs, err := colection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	for curs.Next(ctx) {
		// var book Book
		// if err := curs.Decode(&book); err != nil {
		// 	log.Fatal(err)
		// }
		// log.Println(book)

		var res bson.M
		if err := curs.Decode(&res); err != nil {
			log.Fatal(err)
		}
		log.Println(res)
	}

	if err := curs.Err(); err != nil {
		log.Fatal(err)
	}
}
