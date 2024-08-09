// handler/handler.go
package handler

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "goroutine/docs"

	"goroutine/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	files "github.com/swaggo/files" // Swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title           My API
// @version         1.0
// @description     This is a sample server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

var client *mongo.Client
var itemsCollection = "items"

type Item struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Price int                `json:"price" bson:"price"`
}

// PaginatedItemsResponse represents a paginated response for items.
// @Description Represents a paginated response for items
// @Example {"page": 1, "limit": 10, "totalPages": 5, "items": [{"id": "607c191e810c19729de860ea", "name": "Sample Item", "price": 100}]}
type PaginatedItemsResponse struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalPages int    `json:"totalPages"` // Add this field
	TotalCount int    `json:"totalCount"` // Add this field
	Items      []Item `json:"items"`
}

// Response represents a standard API response.
// @Description Represents a standard response for errors or success messages
// @Example {"message": "Item updated", "id": "607c191e810c19729de860ea"}
type Response struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

func InitMongoDB() {
	// Load environment variables from .env file
	if strings.ToLower(os.Getenv("ENVIRONMENT")) != "production" {
		fail := godotenv.Load()
		if fail != nil {
			log.Println("Error loading .env file:", fail)
		}
	}
	mongoUser := os.Getenv("MONGO_USER")
	mongoPass := os.Getenv("MONGO_PASS")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")
	mongoDB := os.Getenv("MONGO_DB")
	if mongoDB == "" {
		log.Fatal("MONGO_DB not set in environment")
	}

	// Build the MongoDB URI
	var mongoURI = fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)

	if strings.ToLower(mongoHost) != "localhost" {
		mongoURI = fmt.Sprintf("mongodb+srv://%s/?retryWrites=true&w=majority", mongoHost)
	}

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	credential := options.Credential{
		Username: mongoUser,
		Password: mongoPass,
	}
	clientOptions = clientOptions.SetAuth(credential)
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v\n", err)
	}

	fmt.Println("Connecting to MongoDB...")
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v\n", err)
	}

	fmt.Println("Connected to MongoDB!")

	// Defer disconnecting from MongoDB
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v\n", err)
		}
	}()
}

func SetupRouter() *gin.Engine {
	// Set up Gin
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"http://localhost:3000"}, // Allow specific origin
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	// Public route
	r.POST("/login", func(c *gin.Context) {
		var loginData struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
			return
		}

		// Authenticate user (e.g., check username and password)
		// For demo purposes, assume authentication is successful
		token, err := middleware.GenerateToken(loginData.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	})
	// Example route to get data from MongoDB
	r.GET("/items", middleware.AuthMiddleware(), GetItems)
	r.POST("/items", middleware.AuthMiddleware(), CreateItem)
	r.PUT("/items/:id", middleware.AuthMiddleware(), UpdateItem)    // Update route to use ID
	r.DELETE("/items/:id", middleware.AuthMiddleware(), DeleteItem) // Add this line

	return r
}

// Handler is the exported function for Vercel to use as an entry point.
func Handler(w http.ResponseWriter, r *http.Request) {
	InitMongoDB()
	router := SetupRouter()
	router.ServeHTTP(w, r)
}

// GetItems retrieves items from MongoDB with pagination support
// @Summary Get items
// @Description Get all items from the database with pagination
// @Tags Items
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} PaginatedItemsResponse
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /items [get]
func GetItems(c *gin.Context) {
	if client == nil {
		c.JSON(http.StatusInternalServerError, Response{Message: "MongoDB client is not initialized"})
		return
	}

	collection := client.Database(os.Getenv("MONGO_DB")).Collection(itemsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse query parameters for pagination
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, Response{Message: "Invalid page number"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, Response{Message: "Invalid limit"})
		return
	}

	// Calculate the number of documents to skip
	skip := (page - 1) * limit

	// Retrieve the total count of documents
	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
		return
	}

	// Calculate the total number of pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	// Retrieve items from MongoDB with pagination
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var items []Item
	for cursor.Next(ctx) {
		var item Item
		if err := cursor.Decode(&item); err != nil {
			c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
			return
		}
		items = append(items, item)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedItemsResponse{
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,      // Include total pages in the response
		TotalCount: int(totalCount), // Include total count in the response
		Items:      items,
	})
}

// CreateItem adds a new item to MongoDB
// @Summary Create item
// @Description Add a new item to the database
// @Tags Items
// @Accept json
// @Produce json
// @Param item body Item true "Item to add"
// @Success 201 {object} Item
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /items [post]
func CreateItem(c *gin.Context) {
	var newItem Item
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}

	collection := client.Database(os.Getenv("MONGO_DB")).Collection(itemsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if newItem.ID.IsZero() {
		newItem.ID = primitive.NewObjectID()
	}

	if _, err := collection.InsertOne(ctx, newItem); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newItem)
}

// UpdateItem updates an existing item in MongoDB by ID
// @Summary Update item
// @Description Update an existing item in the database by ID
// @Tags Items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param item body Item true "Updated item data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /items/{id} [put]
func UpdateItem(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: "Invalid item ID"})
		return
	}

	var updatedItem Item
	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}

	collection := client.Database(os.Getenv("MONGO_DB")).Collection(itemsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": updatedItem}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, Response{Message: "Item not found"})
		return
	}

	c.JSON(http.StatusOK, Response{Message: "Item updated", ID: id})
}

// DeleteItem removes an item from MongoDB by ID
// @Summary Delete item
// @Description Delete an item from the database by ID
// @Tags Items
// @Param id path string true "Item ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /items/{id} [delete]
func DeleteItem(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: "Invalid item ID"})
		return
	}

	collection := client.Database(os.Getenv("MONGO_DB")).Collection(itemsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": objectID}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, Response{Message: "Item not found"})
		return
	}

	c.JSON(http.StatusOK, Response{Message: "Item deleted", ID: id})
}
