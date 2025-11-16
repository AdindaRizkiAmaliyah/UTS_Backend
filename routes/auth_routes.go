package routes

import (
	"clean-archi/app/service"
	repo "clean-archi/app/repository/MongoRepo"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

// AuthRoutes — endpoint login dan register
func AuthRoutes(router fiber.Router, mongoClient *mongo.Client) {
	mongoDatabase := mongoClient.Database(os.Getenv("MONGODB_DATABASE"))

	mongoRepo := repo.NewAlumniMongoRepository(mongoDatabase, "alumni")
	authService := service.NewAuthService(mongoRepo)

	auth := router.Group("/api/v1/auth")

	auth.Post("/register", authService.Register)
	auth.Post("/login", authService.Login)
}
