package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

type Repository struct {
	Client *mongo.Client
}

func (r *Repository) insert() {
	collection := r.Client.Database("mydb").Collection("trainers")
	_, err := collection.InsertOne(context.Background(), Trainer{"Ash", 10, "Pallet Town"})
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://mongoadmin:secret@localhost:27017")
	clientOptions.SetSocketTimeout(20 * time.Second)
	clientOptions.SetServerSelectionTimeout(time.Second)
	clientOptions.SetConnectTimeout(20 * time.Second)
	// clientOptions.SetMaxConnIdleTime(20 * time.Second)
	clientOptions.SetMinPoolSize(100)
	clientOptions.SetMaxPoolSize(100)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.TODO())
	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		// c.Response().Header().Set("Connection", "close")
		// defer c.Request().Body.Close()
		r := &Repository{
			Client: client,
		}
		r.insert()
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
