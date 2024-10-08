# Go Project Setup and Swagger Documentation Guide

## Overview

This README provides instructions for setting up a Go project, running the application, generating Swagger documentation, and deploying it with ReDocly.

## Prerequisites

- Go (1.16+)
- MongoDB
- Git
- NPM

## Project Setup

### 1. Clone the Repository

Clone your Go project repository:

```bash
git clone https://github.com/# Go Project Setup and Swagger Documentation Guide
```

## Overview

This README provides instructions for setting up a Go project, running the application, generating Swagger documentation, and deploying it with ReDocly.

## Prerequisites

- Go (1.16+)
- MongoDB
- Git
- Docker (for ReDocly deployment)

## Project Setup

### 1. Clone the Repository

Clone your Go project repository:

```bash
git clone https://github.com/yourusername/yourproject.git
cd yourprojectyourusername/yourproject.git
cd yourproject
```

### 2. Install Dependencies

Make sure you have Go installed and your GOPATH is set up. Then, install the required Go dependencies:

```bash
go mod tidy
```

### 3. Set Up Environment Variables

Create a .env file in the root of your project and configure the necessary environment variables:

```env
MONGO_USER=your_mongo_user
MONGO_PASS=your_mongo_password
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_DB=your_database_name
PORT=8080
```

### 4. Run the Application

Start the application:

```bash
go run main.go
```

The application will start and listen on port 8080 (or the port specified in your .env file).

## Generate Swagger Documentation

### 1. Install Swagger Tools

You need to install the swaggo tools for generating Swagger documentation:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 2. Annotate Your Code

Ensure your Go code is properly annotated with Swagger comments. For example:

```go
// @title My API
// @version 1.0
// @description This is a sample server.

// @host localhost:8080
// @BasePath /
```

Annotate your routes with Swagger comments:

```go
// getItems retrieves items from MongoDB with pagination support
// @Summary Get items
// @Description Get all items from the database with pagination
// @Tags items
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} PaginatedItemsResponse "Successfully retrieved items"
// @Failure 400 {object} Response "Invalid request"
// @Failure 500 {object} Response "Internal server error"
// @Router /items [get]
func getItems(c *gin.Context) {
    // Function implementation
}

```

### 3. Generate Swagger Files

Run the swag command to generate the Swagger documentation:

```bash
swag init
```

### 4. Serve Swagger UI

To view the Swagger UI, you need to add a route to serve it in your main.go:

```go
import (
    "github.com/swaggo/files"
    "github.com/swaggo/gin-swagger"
)

func main() {
    // Other setup code...

    r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

    // Run the server
    r.Run(":" + port)
}
```

## Deploy Swagger Documentation with ReDocly

### 1. Create a ReDocly Account

Sign up for a ReDocly account at ReDocly.

### 2. Install ReDocly CLI

Install the ReDocly CLI tool:

```bash
npm install -g @redocly/cli
```

### 3. Initialize ReDocly

Navigate to your project directory and initialize ReDocly:

```bash
redocly init
```

Follow the prompts to set up your project.

### 4. Deploy Swagger Documentation

Deploy your Swagger documentation to ReDocly:

```bash
redocly preview-docs
```

This command will preview your API documentation locally. To deploy it, use:

```bash
redocly push --file docs/swagger.yaml
```

This uploads your Swagger documentation to ReDocly, where you can manage and publish it.

### Troubleshooting

<ul>
  <li>Swagger Annotations Not Showing: Ensure that your Swagger annotations are correct and re-run swag init.</li>
  <li>Environment Variables: Double-check that all required environment variables are set correctly in your .env file.</li>
</ul>
Additional Resources
<ul>
  <li>[Go Documentation](https://golang.org/doc/)
  </li>
  <li>[Swaggo GitHub Repository](https://github.com/swaggo/swag)
  </li>
  <li>[ReDocly Documentation](https://redocly.com/docs)
  </li>
</ul>

### License
This project is licensed under the MIT License. See the LICENSE file for details.


### Explanation:

- **Setup Instructions**: Instructions to clone the repository, install dependencies, and run the Go application.
- **Swagger Documentation**: Steps for installing Swagger tools, annotating code, generating Swagger files, and serving Swagger UI.
- **Deploying with ReDocly**: Steps to sign up, install CLI tools, initialize, and deploy Swagger documentation with ReDocly.
- **Troubleshooting**: Common issues and their solutions.
- **Additional Resources**: Links to Go documentation, Swaggo, and ReDocly resources.
- **License**: Licensing information for the project.