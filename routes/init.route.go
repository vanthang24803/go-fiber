package routes

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type RouteFunc func(router fiber.Router, db *mongo.Database)

type RouteGroup struct {
	Router fiber.Router
	DB     *mongo.Database
}

func NewRouteGroup(router fiber.Router, db *mongo.Database) *RouteGroup {
	return &RouteGroup{
		Router: router,
		DB:     db,
	}
}

func (rg *RouteGroup) Apply(routeFuncs ...RouteFunc) {
	for _, rf := range routeFuncs {
		rf(rg.Router, rg.DB)
	}
}

func InitRoutes(router fiber.Router, db *mongo.Database) {
	rg := NewRouteGroup(router, db)
	rg.Apply(
		AuthRoutes,
		MeRoutes,
	)
}
