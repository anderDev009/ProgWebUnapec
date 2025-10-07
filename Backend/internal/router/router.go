package router

import (
	"petmatch/internal/config"
	"petmatch/internal/handlers"
	"petmatch/internal/middleware"
	"petmatch/internal/models"
	"petmatch/internal/repositories"
	"petmatch/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func New(db *gorm.DB, cfg config.Config) (*gin.Engine, error) {
	userRepo := repositories.NewUserRepository(db)
	petRepo := repositories.NewPetRepository(db)
	adoptionRepo := repositories.NewAdoptionRepository(db)

	authService, err := services.NewAuthService(userRepo, cfg)
	if err != nil {
		return nil, err
	}

	petService := services.NewPetService(petRepo)
	adoptionService := services.NewAdoptionService(adoptionRepo, petRepo)

	authHandler := handlers.NewAuthHandler(authService)
	petHandler := handlers.NewPetHandler(petService)
	adoptionHandler := handlers.NewAdoptionHandler(adoptionService)
	adminHandler := handlers.NewAdminHandler(userRepo, authService)

	r := gin.Default()

	v1 := r.Group("/api/v1")

	authRoutes := v1.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.GET("/me", middleware.Authentication(authService), handlers.CurrentUserHandler)
	}

	v1.GET("/pets", petHandler.List)
	v1.GET("/pets/:id", petHandler.Get)

	authMiddleware := middleware.Authentication(authService)

    shelterGroup := v1.Group("")
    shelterGroup.Use(authMiddleware, middleware.RequireRoles(models.RoleShelter))
    {
        shelterPets := shelterGroup.Group("/pets")
        shelterPets.POST("", petHandler.Create)
        shelterPets.PUT("/:id", petHandler.Update)
        shelterPets.DELETE("/:id", petHandler.Delete)

    }

    adopterGroup := v1.Group("")
    adopterGroup.Use(authMiddleware, middleware.RequireRoles(models.RoleAdopter))
    {
        adopterGroup.POST("/pets/:petId/adoption-requests", adoptionHandler.Create)
    }

    // Shared route for listing adoption requests based on role
    requestsGroup := v1.Group("")
    requestsGroup.Use(authMiddleware, middleware.RequireRoles(models.RoleShelter, models.RoleAdopter))
    {
        requestsGroup.GET("/adoption-requests", adoptionHandler.List)
    }

    shelterGroup.PATCH("/adoption-requests/:id", adoptionHandler.UpdateStatus)

	adminGroup := v1.Group("/admin")
	adminGroup.Use(authMiddleware, middleware.RequireRoles(models.RoleAdmin))
	{
		adminGroup.GET("/users", adminHandler.ListUsers)
		adminGroup.POST("/shelters/:id/approve", adminHandler.ApproveShelter)
	}

	return r, nil
}
