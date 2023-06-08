package main

import (
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	joonix "github.com/joonix/log"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"

	"stark/database"
	"stark/services"
	"stark/services/auth"
	"stark/services/client"
	"stark/services/email_verification"
	"stark/services/profile"
	"stark/services/user"
	"stark/services/user_detail"
	"stark/services/user_location"
	"stark/utils"
	"stark/utils/activity"
	"stark/utils/log"
	"stark/utils/middleware"
)

func main() {
	// Define log
	configureLogging()
	ctx := activity.NewContext("init_stark")
	ctx = activity.WithClientID(ctx, "stark_system")

	// Database repository for service
	mysqlDB, err := database.NewMySQL()
	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "mysql connection error"))
		return
	}

	redisDB, err := database.NewRedis()
	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "redis connection error"))
		return
	}

	_, err = database.NewMongo()
	if err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "mongo connection error"))
		return
	}

	// Define repository, service, handler
	clientRepo := client.NewSQLRepository(mysqlDB)
	clientService := client.NewService(clientRepo)
	clientHandler := client.NewHandler(clientService)
	userRepo := user.NewSQLRepository(mysqlDB)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	userDetailRepo := user_detail.NewSQLRepository(mysqlDB)
	userDetailService := user_detail.NewService(userDetailRepo)
	userDetailHandler := user_detail.NewHandler(userDetailService)
	userLocationRepo := user_location.NewSQLRepository(mysqlDB)
	userLocationService := user_location.NewService(userLocationRepo)
	userLocationHandler := user_location.NewHandler(userLocationService)
	emailVerificationRepo := email_verification.NewSQLRepository(mysqlDB)
	emailVerificationService := email_verification.NewService(emailVerificationRepo)
	authService := auth.NewService(redisDB, userService, emailVerificationService)
	authHandler := auth.NewHandler(authService)
	profileService := profile.NewService(userService, userDetailService, userLocationService)
	profileHandler := profile.NewHandler(profileService)

	// Set application mode
	mode := os.Getenv("APP_MODE")
	gin.SetMode(mode)

	// Define application
	app := gin.Default()
	app.Use(
		gin.Recovery(),
		gin.Logger(),
		middleware.CORSMiddleware(),
		middleware.JSONMiddleware(),
	)

	// Init route
	services.InitRoute(
		ctx,
		app,
		redisDB,
		clientService,
		clientHandler,
		userHandler,
		userDetailHandler,
		userLocationHandler,
		authHandler,
		profileHandler,
	)

	// Let's get started!
	if err := app.Run(":" + os.Getenv("SERVER_PORT")); err != nil {
		log.WithContext(ctx).Error(stacktrace.Propagate(err, "stark running error"))
		return
	}
}

func configureLogging() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.AddHook(utils.LogrusSourceContextHook{})

	if gin.Mode() != "release" {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	} else {
		logrus.SetFormatter(&joonix.FluentdFormatter{})
	}
}
