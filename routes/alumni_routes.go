package routes

import (
	"clean-archi/app/service"
	repo "clean-archi/app/repository/MongoRepo"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

// AlumniRoutes — rute CRUD untuk data alumni (MongoDB)
func AlumniRoutes(router fiber.Router, mongoClient *mongo.Client) {
	mongoDatabase := mongoClient.Database(os.Getenv("MONGODB_DATABASE"))

	// Repository & Service
	mongoRepo := repo.NewAlumniMongoRepository(mongoDatabase, "alumni")
	alumniService := service.NewAlumniService(mongoRepo)

	// Base group
	alumni := router.Group("/api/v1/alumni")

	// Endpoints
	alumni.Get("/", alumniService.GetAll)
	alumni.Get("/:id", alumniService.GetByID)
	alumni.Post("/", alumniService.Create)
	alumni.Put("/:id", alumniService.Update)
	alumni.Delete("/:id", alumniService.Delete)
}
