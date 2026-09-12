package api

import (
	"go-fiber-digilib/api/handlers"
	"go-fiber-digilib/api/middlewares"
	"go-fiber-digilib/auth"
	"go-fiber-digilib/books"
	"go-fiber-digilib/carts"
	"go-fiber-digilib/categories"
	"go-fiber-digilib/pkg/ai"
	"go-fiber-digilib/pkg/rajaongkir"
	"go-fiber-digilib/rents"
	"go-fiber-digilib/settings"
	"go-fiber-digilib/users"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gofiber/fiber/v3"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewApp(repository *gorm.DB, cld *cloudinary.Cloudinary, jwtConfig middlewares.JWTConfig) *fiber.App {
	var (
		app = fiber.New(fiber.Config{
			StructValidator: &middlewares.CustomValidator{
				Validator: middlewares.InitValidator(),
			},
		})
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

	app.Use(logMiddleware.Init())

	authRoutes := app.Group("/api/v1/auth")

	authRoutes.Post("/register", authHandler.Register)
	authRoutes.Post("/login", authHandler.Login)

	userRoutes := app.Group("/api/v1", echojwt.WithConfig(jwtMiddleware), middlewares.VerifyToken)

	userRoutes.Get("/profile", usersHandler.GetProfile)
	userRoutes.Put("/profile/edit", usersHandler.EditProfile)

	categoryRoutes := app.Group("/api/v1", echojwt.WithConfig(jwtMiddleware), middlewares.VerifyToken)
	categoryRoutes.Get("/categories", categoriesHandler.GetAll)
	categoryRoutes.Get("/categories/:id", categoriesHandler.GetByID)
	categoryRoutes.Post("/categories", categoriesHandler.Create, middlewares.VerifyAdmin, middlewares.ValidateBody(&categories.CategoryRequest{}))
	categoryRoutes.Put("/categories/:id", categoriesHandler.Update, middlewares.VerifyAdmin, middlewares.ValidateBody(&categories.CategoryRequest{}))
	categoryRoutes.Delete("/categories/:id", categoriesHandler.Delete, middlewares.VerifyAdmin)

	bookRoutes := app.Group("/api/v1", echojwt.WithConfig(jwtMiddleware), middlewares.VerifyToken)
	bookRoutes.Get("/books", booksHandler.GetAll)
	bookRoutes.Get("/books/:id", booksHandler.GetByID)
	bookRoutes.Get("/books/category/:id", booksHandler.GetByCategory)
	bookRoutes.Post("/books", booksHandler.Create, middlewares.VerifyAdmin, middlewares.ValidateBody(&books.BookRequest{}))
	bookRoutes.Put("/books/:id", booksHandler.Update, middlewares.VerifyAdmin, middlewares.ValidateBody(&books.BookRequest{}))
	bookRoutes.Delete("/books/:id", booksHandler.Delete, middlewares.VerifyAdmin)
	bookRoutes.Post("/books/recommendations", booksHandler.GetBookRecommendation, middlewares.ValidateBody(&ai.BookRecommendationRequest{}))

	cartRoutes := app.Group("/api/v1", echojwt.WithConfig(jwtMiddleware), middlewares.VerifyToken)
	cartRoutes.Get("/carts/user", cartsHandler.GetByUser)
	cartRoutes.Post("/carts", cartsHandler.Create, middlewares.ValidateBody(&carts.CartRequest{}))
	cartRoutes.Put("/carts/:id", cartsHandler.Update, middlewares.ValidateBody(&carts.CartRequest{}))
	cartRoutes.Delete("/carts/:id", cartsHandler.Delete)

	rentRoutes := app.Group("/api/v1", echojwt.WithConfig(jwtMiddleware), middlewares.VerifyToken)
	rentRoutes.Get("/rents", rentsHandler.GetAll, middlewares.VerifyAdmin)
	rentRoutes.Get("/rents/user", rentsHandler.GetByUser)
	rentRoutes.Post("/rents", rentsHandler.Create, middlewares.ValidateBody(&rents.RentRequest{}))
	rentRoutes.Put("/rents/:id", rentsHandler.Update, middlewares.VerifyAdmin, middlewares.ValidateBody(&rents.RentUpdateRequest{}))
	rentRoutes.Delete("/rents/:id", rentsHandler.Delete, middlewares.VerifyAdmin)

	settingRoutes := app.Group("/api/v1", echojwt.WithConfig(jwtMiddleware), middlewares.VerifyToken, middlewares.VerifyAdmin)
	settingRoutes.Get("/settings", settingsHandler.GetAll)
	settingRoutes.Get("/settings/:id", settingsHandler.GetByID)
	settingRoutes.Get("/settings/key/:key", settingsHandler.GetByKey)
	settingRoutes.Post("/settings", settingsHandler.Create, middlewares.ValidateBody(&settings.SettingRequest{}))
	settingRoutes.Put("/settings/:id", settingsHandler.Update, middlewares.ValidateBody(&settings.SettingRequest{}))
	settingRoutes.Delete("/settings/:id", settingsHandler.Delete)

	return app
}
