package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"restapirian/config"
	delivery "restapirian/internal/delivery/http"
	"restapirian/internal/domain"
	"restapirian/internal/repository"
	"restapirian/internal/service"
	jwtauth "restapirian/pkg/auth"

	"github.com/gin-gonic/gin"
)

func InitializeRouter() error {
	db, err := config.ConnectDb()
	if err != nil {
		return err
	}

	productRepository := repository.NewProductRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	userRepository := repository.NewUserRepository(db)
	roleRepository := repository.NewRoleRepository(db)

	if err := service.BootstrapRBAC(context.Background(), userRepository, roleRepository, service.BootstrapAdminInput{
		Username: os.Getenv("ADMIN_USERNAME"), Password: os.Getenv("ADMIN_PASSWORD"),
		Email: os.Getenv("ADMIN_EMAIL"), Name: os.Getenv("ADMIN_NAME"),
	}); err != nil {
		return fmt.Errorf("bootstrap RBAC: %w", err)
	}

	tokenDuration := time.Hour
	if value := os.Getenv("JWT_TTL_MINUTES"); value != "" {
		minutes, parseErr := strconv.Atoi(value)
		if parseErr != nil || minutes < 1 {
			return fmt.Errorf("JWT_TTL_MINUTES must be a positive integer")
		}
		tokenDuration = time.Duration(minutes) * time.Minute
	}
	tokens, err := jwtauth.NewJWTManager(os.Getenv("JWT_SECRET"), "restapirian-api", tokenDuration)
	if err != nil {
		return err
	}

	productController := delivery.NewProductController(service.NewProductService(productRepository, categoryRepository))
	categoryController := delivery.NewCategoryController(service.NewCategoryService(categoryRepository))
	userController := delivery.NewUserController(service.NewUserService(userRepository, roleRepository))
	roleController := delivery.NewRoleController(service.NewRoleService(roleRepository))
	authController := delivery.NewAuthController(service.NewAuthService(userRepository, roleRepository, tokens))

	engine := gin.Default()
	if err := engine.SetTrustedProxies(nil); err != nil {
		return fmt.Errorf("configure trusted proxies: %w", err)
	}
	engine.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	authRoutes := engine.Group("/api/v1/auth")
	authRoutes.POST("/register", authController.Register)
	authRoutes.POST("/login", authController.Login)

	protected := engine.Group("/api/v1", JWTAuth(tokens))
	readAccess := RequireRoles(domain.RoleAdmin, domain.RoleUser)
	adminOnly := RequireRoles(domain.RoleAdmin)

	products := protected.Group("/product")
	products.GET("/", readAccess, productController.Index)
	products.GET("/show/:id", readAccess, productController.Show)
	products.POST("/create", adminOnly, productController.Create)
	products.PUT("/update/:id", adminOnly, productController.Update)
	products.DELETE("/delete/:id", adminOnly, productController.Delete)

	categories := protected.Group("/category")
	categories.GET("/", readAccess, categoryController.Index)
	categories.GET("/show/:id", readAccess, categoryController.Show)
	categories.POST("/create", adminOnly, categoryController.Create)
	categories.PUT("/update/:id", adminOnly, categoryController.Update)
	categories.DELETE("/delete/:id", adminOnly, categoryController.Delete)

	users := protected.Group("/user", adminOnly)
	users.GET("/", userController.Index)
	users.GET("/show/:id", userController.Show)
	users.POST("/create", userController.Create)
	users.PUT("/update/:id", userController.Update)
	users.DELETE("/delete/:id", userController.Delete)

	roles := protected.Group("/role", adminOnly)
	roles.GET("/", roleController.Index)

	return engine.Run(":9090")
}
