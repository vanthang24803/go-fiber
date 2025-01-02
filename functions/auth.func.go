package functions

import (
	"context"
	"time"

	"github.com/vanthang24803/go-api/domain"
	"github.com/vanthang24803/go-api/internal/schema"
	"github.com/vanthang24803/go-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *mongo.Database
}

func NewAuthService(db *mongo.Database) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) RegisterHandler(request *domain.RegisterRequest) (*schema.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &schema.User{
		Email:        request.Email,
		Avatar:       "",
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		Username:     request.Username,
		Roles:        []string{"user"},
		HashPassword: string(hash),
	}

	usersCollection := s.db.Collection("users")

	result, err := usersCollection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return user, nil
}

func (s *AuthService) LoginHandler(request *domain.LoginRequest) (*utils.TokenResponse, *utils.AppError) {
	var user schema.User

	// Find user
	err := s.db.Collection("users").FindOne(context.Background(), bson.M{"username": request.Username}).Decode(&user)

	// Check if user exists
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.NewAppError(404, "User not found")
		}
		return nil, utils.NewAppError(500, err.Error())
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(request.Password))
	if err != nil {
		return nil, utils.NewAppError(401, "Username or password is incorrect")
	}

	// Check for valid refresh token
	var refreshToken string
	var needNewRefreshToken bool = true

	for _, token := range user.Tokens {
		if token.Name == "refresh_token" && token.ExpiredAt.After(time.Now()) {
			refreshToken = token.Value
			needNewRefreshToken = false
			break
		}
	}

	// Generate JWT
	var token *utils.TokenResponse

	token, err = utils.GenerateJWT(&user)

	if err != nil {
		return nil, utils.NewAppError(400, err.Error())
	}

	if needNewRefreshToken {
		// Add new refresh token to user's tokens
		user.Tokens = append(user.Tokens, schema.Token{
			Name:      "refresh_token",
			Value:     token.RefreshToken,
			ExpiredAt: time.Now().Add(time.Hour * 24 * 7), // 7 days expiration
		})

		// Update user in database
		_, err = s.db.Collection("users").UpdateOne(
			context.Background(),
			bson.M{"_id": user.ID},
			bson.M{"$set": bson.M{"tokens": user.Tokens}},
		)
		if err != nil {
			return nil, utils.NewAppError(500, "Failed to update user tokens")
		}
	} else {
		// Generate new access token with existing refresh token
		token.RefreshToken = refreshToken
	}

	// Return
	return token, nil
}
