package api

import (
	"go-digilib/api/handlers"
	"go-digilib/api/middlewares"
	"go-digilib/auth"
	"go-digilib/books"
	"go-digilib/carts"
	"go-digilib/categories"
	"go-digilib/pkg/ai"
	"go-digilib/pkg/rajaongkir"
	"go-digilib/rents"
	"go-digilib/settings"
	"go-digilib/users"

	"github.com/cloudinary/cloudinary-go/v2"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewEcho(repository *gorm.DB, cld *cloudinary.Cloudinary, jwtConfig middlewares.JWTConfig) *echo.Echo {
	var (
		e                     = echo.New()
		categoryService       = categories.New(repository)
		bookService           = books.New(repository)
		authService           = auth.New(repository)
		userService           = users.New(repository)
		settingService        = settings.New(repository)
		cartService           = carts.New(repository)
		rentService           = rents.New(repository)
		roService             = rajaongkir.InitService()
		recommendationService = ai.InitService()
		authHandler           = handlers.NewAuth(authService, jwtConfig)
		usersHandler          = handlers.NewUsers(userService, cld)
		categoriesHandler     = handlers.NewCategories(categoryService)
		booksHandler          = handlers.NewBooks(bookService, cld, recommendationService)
		settingsHandler       = handlers.NewSettings(settingService)
		cartsHandler          = handlers.NewCarts(cartService)
		rentsHandler          = handlers.NewRents(
			rentService, settingService, userService, roService,
		)
	)

	e.Validator = &middlewares.CustomValidator{
		Validator: middlewares.InitValidator(),
	}

	logger, _ := zap.NewProduction()
	loggerConfig := middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}

	logMiddleware := middlewares.LoggerConfig{Config: loggerConfig}
	jwtMiddleware := jwtConfig.Init()

	e.Use(logMiddleware.Init())

	authRoutes := e.Group("/api/v1/auth")

	authRoutes.POST("/register", authHandler.Register)
	authRoutes.POST("/login", authHandler.Login)

	userRoutes := e.Group("/api/v1", echojwt.WithConfig(jwtMiddleware), middlewares.VerifyToken)

	userRoutes.GET("/profile", usersHandler.GetProfile)
	userRoutes.PATCH("/profile/edit", usersHandler.EditProfile)

	categoryRoutes := e.Group("/api/v1", echojwt.WithConfig(jwtMiddleware), middlewares.VerifyToken)
	categoryRoutes.GET("/categories", categoriesHandler.GetAll)
	categoryRoutes.GET("/categories/:id", categoriesHandler.GetByID)
	categoryRoutes.POST("/categories", categoriesHandler.Create, middlewares.VerifyAdmin, middlewares.ValidateBody(&categories.CategoryRequest{}))
	categoryRoutes.PATCH("/categories/:id", categoriesHandler.Update, middlewares.VerifyAdmin, middlewares.ValidateBody(&categories.CategoryRequest{}))
	categoryRoutes.DELETE("/categories/:id", categoriesHandler.Delete, middlewares.VerifyAdmin)

	bookRoutes := e.Group("/api/v1", echojwt.WithConfig(jwtMiddleware), middlewares.VerifyToken)
	bookRoutes.GET("/books", booksHandler.GetAll)
	bookRoutes.GET("/books/:id", booksHandler.GetByID)
	bookRoutes.GET("/books/category/:id", booksHandler.GetByCategory)
	bookRoutes.POST("/books", booksHandler.Create, middlewares.VerifyAdmin, middlewares.ValidateBody(&books.BookRequest{}))
	bookRoutes.PATCH("/books/:id", booksHandler.Update, middlewares.VerifyAdmin, middlewares.ValidateBody(&books.BookRequest{}))
	bookRoutes.DELETE("/books/:id", booksHandler.Delete, middlewares.VerifyAdmin)
	bookRoutes.POST("/books/recommendations", booksHandler.GetBookRecommendation, middlewares.ValidateBody(&ai.BookRecommendationRequest{}))

	cartRoutes := e.Group("/api/v1", echojwt.WithConfig(jwtMiddleware), middlewares.VerifyToken)
	cartRoutes.GET("/carts/user", cartsHandler.GetByUser)
	cartRoutes.POST("/carts", cartsHandler.Create, middlewares.ValidateBody(&carts.CartRequest{}))
	cartRoutes.PATCH("/carts/:id", cartsHandler.Update, middlewares.ValidateBody(&carts.CartRequest{}))
	cartRoutes.DELETE("/carts/:id", cartsHandler.Delete)

	rentRoutes := e.Group("/api/v1", echojwt.WithConfig(jwtMiddleware), middlewares.VerifyToken)
	rentRoutes.GET("/rents", rentsHandler.GetAll, middlewares.VerifyAdmin)
	rentRoutes.GET("/rents/user", rentsHandler.GetByUser)
	rentRoutes.POST("/rents", rentsHandler.Create, middlewares.ValidateBody(&rents.RentRequest{}))
	rentRoutes.PATCH("/rents/:id", rentsHandler.Update, middlewares.VerifyAdmin, middlewares.ValidateBody(&rents.RentUpdateRequest{}))
	rentRoutes.DELETE("/rents/:id", rentsHandler.Delete, middlewares.VerifyAdmin)

	settingRoutes := e.Group("/api/v1", echojwt.WithConfig(jwtMiddleware), middlewares.VerifyToken, middlewares.VerifyAdmin)
	settingRoutes.GET("/settings", settingsHandler.GetAll)
	settingRoutes.GET("/settings/:id", settingsHandler.GetByID)
	settingRoutes.GET("/settings/key/:key", settingsHandler.GetByKey)
	settingRoutes.POST("/settings", settingsHandler.Create, middlewares.ValidateBody(&settings.SettingRequest{}))
	settingRoutes.PATCH("/settings/:id", settingsHandler.Update, middlewares.ValidateBody(&settings.SettingRequest{}))
	settingRoutes.DELETE("/settings/:id", settingsHandler.Delete)

	return e
}
