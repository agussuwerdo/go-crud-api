// main_test.go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"goroutine/handler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testClient *mongo.Client
var testCollection *mongo.Collection

func setup() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	mongoUser := os.Getenv("MONGO_USER")
	mongoPass := os.Getenv("MONGO_PASS")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")
	mongoDB := os.Getenv("MONGO_DB")

	var mongoURI = "mongodb://" + mongoHost + ":" + mongoPort

	clientOptions := options.Client().ApplyURI(mongoURI).SetAuth(options.Credential{
		Username: mongoUser,
		Password: mongoPass,
	})
	testClient, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}

	testCollection = testClient.Database(mongoDB).Collection("items")
}

func teardown() {
	testClient.Disconnect(context.Background())
}

func TestGetItems(t *testing.T) {
	setup()
	defer teardown()

	r := gin.Default()
	r.GET("/items", handler.GetItems)

	req, _ := http.NewRequest("GET", "/items", nil)

	fmt.Println(req)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d. Response body: %s", http.StatusOK, resp.Code, resp.Body.String())
	} else {
		assert.Equal(t, http.StatusOK, resp.Code)
	}
}

func TestCreateItem(t *testing.T) {
	setup()
	defer teardown()

	r := gin.Default()
	r.POST("/items", handler.CreateItem)

	item := handler.Item{Name: "Test Item", Price: 100}
	jsonData, _ := json.Marshal(item)
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestUpdateItem(t *testing.T) {
	setup()
	defer teardown()

	r := gin.Default()
	r.PUT("/items/:id", handler.UpdateItem)

	item := handler.Item{Name: "Test Item", Price: 100}
	insertResult, _ := testCollection.InsertOne(context.Background(), item)

	id := insertResult.InsertedID.(primitive.ObjectID).Hex()
	updatedItem := handler.Item{Name: "Updated Item", Price: 150}
	jsonData, _ := json.Marshal(updatedItem)
	req, _ := http.NewRequest("PUT", "/items/"+id, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDeleteItem(t *testing.T) {
	setup()
	defer teardown()

	r := gin.Default()
	r.DELETE("/items/:id", handler.DeleteItem)

	item := handler.Item{Name: "Test Item", Price: 100}
	insertResult, _ := testCollection.InsertOne(context.Background(), item)

	id := insertResult.InsertedID.(primitive.ObjectID).Hex()
	req, _ := http.NewRequest("DELETE", "/items/"+id, nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
