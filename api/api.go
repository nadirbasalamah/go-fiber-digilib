package api

import (
	"go-fiber-digilib/api/handlers"
	"go-fiber-digilib/api/middlewares"
	"go-fiber-digilib/auth"
	"go-fiber-digilib/books"
	"go-fiber-digilib/carts"
	"go-fiber-digilib/categories"
	"go-fiber-digilib/pkg/ai"
	"go-fiber-digilib/pkg/fileupload"
	"go-fiber-digilib/pkg/rajaongkir"
	"go-fiber-digilib/rents"
	"go-fiber-digilib/settings"
	"go-fiber-digilib/users"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/helmet"

	zapMiddleware "github.com/gofiber/contrib/v3/zap"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewApp(repository *gorm.DB, cld *cloudinary.Cloudinary, jwtConfig middlewares.JWTConfig) *fiber.App {
	var (
		app                   = fiber.New()
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
		uploader              = &fileupload.CloudinaryUploader{Cld: cld}
		usersHandler          = handlers.NewUsers(userService, uploader)
		categoriesHandler     = handlers.NewCategories(categoryService)
		booksHandler          = handlers.NewBooks(bookService, uploader, recommendationService)
		settingsHandler       = handlers.NewSettings(settingService)
		cartsHandler          = handlers.NewCarts(cartService)
		rentsHandler          = handlers.NewRents(
			rentService, settingService, userService, roService,
		)
	)

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logMiddleware := middlewares.LoggerConfig{Config: zapMiddleware.Config{
		Logger: logger,
	}}

	jwtMiddleware := jwtConfig.Init()

	app.Use(logMiddleware.Init())
	app.Use(helmet.New())

	authRoutes := app.Group("/api/v1/auth")

	authRoutes.Post("/register", middlewares.ValidateBody(&auth.RegisterRequest{}), authHandler.Register)
	authRoutes.Post("/login", middlewares.ValidateBody(&auth.LoginRequest{}), authHandler.Login)

	userRoutes := app.Group("/api/v1", jwtMiddleware, middlewares.VerifyToken())

	userRoutes.Get("/profile", usersHandler.GetProfile)
	userRoutes.Patch("/profile/edit", middlewares.ValidateBody(&users.EditProfileRequest{}), usersHandler.EditProfile)

	categoryRoutes := app.Group("/api/v1", jwtMiddleware, middlewares.VerifyToken())
	categoryRoutes.Get("/categories", categoriesHandler.GetAll)
	categoryRoutes.Get("/categories/:id", categoriesHandler.GetByID)
	categoryRoutes.Post("/categories", middlewares.VerifyAdmin(), middlewares.ValidateBody(&categories.CategoryRequest{}), categoriesHandler.Create)
	categoryRoutes.Patch("/categories/:id", middlewares.VerifyAdmin(), middlewares.ValidateBody(&categories.CategoryRequest{}), categoriesHandler.Update)
	categoryRoutes.Delete("/categories/:id", middlewares.VerifyAdmin(), categoriesHandler.Delete)

	bookRoutes := app.Group("/api/v1", jwtMiddleware, middlewares.VerifyToken())
	bookRoutes.Get("/books", booksHandler.GetAll)
	bookRoutes.Get("/books/:id", booksHandler.GetByID)
	bookRoutes.Get("/books/category/:id", booksHandler.GetByCategory)
	bookRoutes.Post("/books", middlewares.VerifyAdmin(), middlewares.ValidateBody(&books.BookRequest{}), booksHandler.Create)
	bookRoutes.Patch("/books/:id", middlewares.VerifyAdmin(), middlewares.ValidateBody(&books.BookRequest{}), booksHandler.Update)
	bookRoutes.Delete("/books/:id", middlewares.VerifyAdmin(), booksHandler.Delete)
	bookRoutes.Post("/books/recommendations", middlewares.ValidateBody(&ai.BookRecommendationRequest{}), booksHandler.GetBookRecommendation)

	cartRoutes := app.Group("/api/v1", jwtMiddleware, middlewares.VerifyToken())
	cartRoutes.Get("/carts/user", cartsHandler.GetByUser)
	cartRoutes.Post("/carts", middlewares.ValidateBody(&carts.CartRequest{}), cartsHandler.Create)
	cartRoutes.Patch("/carts/:id", middlewares.ValidateBody(&carts.CartRequest{}), cartsHandler.Update)
	cartRoutes.Delete("/carts/:id", cartsHandler.Delete)

	rentRoutes := app.Group("/api/v1", jwtMiddleware, middlewares.VerifyToken())
	rentRoutes.Get("/rents", middlewares.VerifyAdmin(), rentsHandler.GetAll)
	rentRoutes.Get("/rents/user", rentsHandler.GetByUser)
	rentRoutes.Post("/rents", middlewares.ValidateBody(&rents.RentRequest{}), rentsHandler.Create)
	rentRoutes.Patch("/rents/:id", middlewares.VerifyAdmin(), middlewares.ValidateBody(&rents.RentUpdateRequest{}), rentsHandler.Update)
	rentRoutes.Delete("/rents/:id", middlewares.VerifyAdmin(), rentsHandler.Delete)

	settingRoutes := app.Group("/api/v1", jwtMiddleware, middlewares.VerifyToken(), middlewares.VerifyAdmin())
	settingRoutes.Get("/settings", settingsHandler.GetAll)
	settingRoutes.Get("/settings/:id", settingsHandler.GetByID)
	settingRoutes.Get("/settings/key/:key", settingsHandler.GetByKey)
	settingRoutes.Post("/settings", middlewares.ValidateBody(&settings.SettingRequest{}), settingsHandler.Create)
	settingRoutes.Patch("/settings/:id", middlewares.ValidateBody(&settings.SettingRequest{}), settingsHandler.Update)
	settingRoutes.Delete("/settings/:id", settingsHandler.Delete)

	return app
}
