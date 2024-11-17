package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	UseHTTPS        bool   `json:"use_https"`
	Port            int    `json:"port"`
	CertFile        string `json:"cert_file"`
	KeyFile         string `json:"key_file"`
	MongoURI        string `json:"mongo_uri"`
	MongoDatabase   string `json:"mongo_database"`
	MongoCollection string `json:"mongo_collection"`
}

var config Config
var mongoClient *mongo.Client
var mongoDB *mongo.Database

func server() {
	configPath := flag.String("config", "config.json", "Path to configuration file")
	var version bool
	flag.BoolVar(&version, "version", false, "Show version")
	flag.Parse()
	if version {
		fmt.Println(info())
		os.Exit(0)

	}

	// Parse command-line flags
	//     configPath := flag.String("config", "config.json", "Path to configuration file")
	//     flag.Parse()

	// Load the configuration
	err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to MongoDB
	err = connectMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Set up the HTTP/HTTPS server
	if config.UseHTTPS {
		log.Printf("Starting HTTPS server on port %d", config.Port)
		if err := fasthttp.ListenAndServeTLS(fmt.Sprintf(":%d", config.Port), config.CertFile, config.KeyFile, requestHandler); err != nil {
			log.Fatalf("Error in ListenAndServeTLS: %v", err)
		}
	} else {
		log.Printf("Starting HTTP server on port %d", config.Port)
		if err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", config.Port), requestHandler); err != nil {
			log.Fatalf("Error in ListenAndServe: %v", err)
		}
	}
}

func loadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}

	return nil
}

func connectMongoDB() error {
	clientOptions := options.Client().ApplyURI(config.MongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	mongoClient = client
	mongoDB = client.Database(config.MongoDatabase)
	return nil
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/insert":
		insertHandler(ctx)
	case "/search":
		searchHandler(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}
