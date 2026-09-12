package main

import (
	"fmt"
	"go-fiber-digilib/api"
	"go-fiber-digilib/api/middlewares"
	"go-fiber-digilib/db/drivers"
	"go-fiber-digilib/pkg/constant"
	"go-fiber-digilib/pkg/fileupload"
	"go-fiber-digilib/pkg/utils"
	"log"
	"strconv"
)

func main() {
	dbConfig := drivers.DBConfig{
		Username: utils.GetConfig(constant.DB_USERNAME),
		Password: utils.GetConfig(constant.DB_PASSWORD),
		Database: utils.GetConfig(constant.DB_NAME),
		Host:     utils.GetConfig(constant.DB_HOST),
		Port:     utils.GetConfig(constant.DB_PORT),
	}

	cloudinaryConfig := fileupload.CloudinaryConfig{
		CloudinaryURL: utils.GetConfig(constant.CLOUDINARY_URL),
	}

	expireDuration, err := strconv.Atoi(utils.GetConfig("JWT_EXPIRE_DURATION"))

	if err != nil {
		log.Fatalf("error when parsing expire duration: %v\n", err)
	}

	jwtConfig := middlewares.JWTConfig{
		SecretKey:       utils.GetConfig("JWT_SECRET_KEY"),
		ExpiresDuration: expireDuration,
	}

	var (
		repository = dbConfig.InitDB()
		cloudinary = cloudinaryConfig.InitCloudinary()
		app        = api.NewApp(repository, cloudinary, jwtConfig)
	)

	drivers.MigrateDB(repository)

	appPort := fmt.Sprintf(":%s", utils.GetConfig(constant.PORT))

	if err := app.Listen(appPort); err != nil {
		log.Fatalf("failed to start server %v\n", err)
	}
}
