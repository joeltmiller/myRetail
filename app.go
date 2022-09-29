package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

type App struct {
	Router *mux.Router
}

func (a *App) Initialize(mongoUser string, mongoPass string, baseUrl string) {
	//Initialize Mongo Database Client
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	mongoUri := fmt.Sprintf("mongodb+srv://%s:%s@%s", mongoUser, mongoPass, baseUrl)
	clientOptions := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, mongoError := mongo.Connect(ctx, clientOptions)

	if mongoError != nil {
		log.Fatal(mongoError)
	}

	// Initialize Product Repository with Mongo Client
	rep := productRepository{client: mongoClient}

	// Initialize Red Sky HTTP Client
	client := NewClient(os.Getenv("INTERNAL_BASE_URL"), os.Getenv("INTERNAL_KEY"))

	// Create Router
	a.Router = mux.NewRouter()

	// Initialize Product Controller with Client and Router
	NewProductController(rep, client).SetRoutes(a.Router)
}

func (a *App) Listen(port string) {
	address := fmt.Sprintf(":%s", port)
	fmt.Printf("Listening on port %s\n", address)
	httpErr := http.ListenAndServe(address, a.Router)

	if httpErr != nil {
		fmt.Println("Failed to start web server")
		log.Fatal(httpErr)
	}
}
