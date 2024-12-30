package routes

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gofiber/fiber/v2"
	"github.com/vanthang24803/go-api/internal/database"
	"github.com/vanthang24803/go-api/internal/schema"
	"github.com/vanthang24803/go-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func AuthRoutes(router fiber.Router) {
	router.Get("/register", func(c *fiber.Ctx) error {

		// Generate fake user data
		user := &schema.User{
			Email:     gofakeit.Email(),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Age:       gofakeit.Number(18, 99),
		}

		// Hash the password using bcrypt
		hash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			// If there is an error hashing the password, return an error response
			return c.JSON(utils.NewAppError(500, "Error hashing password"))
		}

		// Set the hashed password on the user object
		user.HashPassword = string(hash)

		// Get the "users" collection from the database
		db := database.GetDatabase().Collection("users")

		// Insert the user into the MongoDB collection
		result, err := db.InsertOne(c.Context(), user)
		if err != nil {
			// If there is an error inserting the user into the database, return an error response
			return c.JSON(utils.NewAppError(500, err.Error()))
		}

		// Assign the generated ID to the user object
		user.ID = result.InsertedID.(primitive.ObjectID)

		// Return the user object with its ID as a response
		return c.JSON(utils.NewResponse(201, user))
	})
}
