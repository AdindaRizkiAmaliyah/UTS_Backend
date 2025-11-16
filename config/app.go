package config

import (
	"clean-archi/app/repository/MongoRepo"
	"clean-archi/app/service"
	"clean-archi/middleware"
	"clean-archi/routes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// NewMongoClient membuat koneksi ke MongoDB
func NewMongoClient() *mongo.Client {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // default jika .env belum diatur
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("❌ Gagal terhubung ke MongoDB:", err)
	}

	// Tes koneksi
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ Tidak dapat ping MongoDB:", err)
	}

	fmt.Println("✅ Koneksi ke MongoDB berhasil!")
	return client
}

// NewApp membuat instance Fiber dengan konfigurasi MongoDB & routes
func NewApp() *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	
	}))

	// Middleware global
	app.Use(middleware.LoggerMiddleware)

	// Koneksi MongoDB
	mongoClient := NewMongoClient()
	dbName := os.Getenv("MONGODB_DATABASE")

	// Inisialisasi repository dan service
	fileRepo := MongoRepo.NewFileRepository(mongoClient, dbName)
	fileService := service.NewFileService(fileRepo, "./uploads")

	// Register semua routes
	routes.AuthRoutes(app, mongoClient)
	routes.AlumniRoutes(app, mongoClient)
	routes.JobRoutes(app, mongoClient)
	routes.FileRoutes(app, fileService)

	// Endpoint root untuk memastikan server aktif
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Server aktif dan MongoDB terhubung ✅",
		})
	})

	return app
}
