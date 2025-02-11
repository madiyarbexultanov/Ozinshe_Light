package main

import (
	"context"
	"goozinshe/config"
	"goozinshe/docs"
	"goozinshe/handlers"
	"goozinshe/logger"
	"goozinshe/middlewares"
	"goozinshe/repositories"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"

	ginzap "github.com/gin-contrib/zap"
	swaggerfiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

// @title           Ozinshe API
// @version         1.0
// @description     Personal online platform providing information about films
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   Madiyar Bexultanov
// @contact.url    https://www.linkedin.com/in/madiyar-bexultanov-b21902258/
// @contact.email  bexultanovmadiyar@gmail.com
//
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host      localhost:8081
// @BasePath  /
//
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
//
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	r := gin.New()

    logger := logger.GetLogger()
    r.Use(
        ginzap.Ginzap(logger, time.RFC3339, true),
        ginzap.RecoveryWithZap(logger, true),
    )
    

    corsConfig := cors.Config{
        AllowAllOrigins:    true,
        AllowHeaders:       []string{"*"},
        AllowMethods:       []string{"*"},
    }

    r.Use(cors.New(corsConfig))
    gin.SetMode(gin.ReleaseMode)

    err := loadConfig()
    if err != nil {
        panic(err)
    }

    conn, err := connectToDb()
    if err!=nil {
        panic(err)
    }

    moviesRepository := repositories.NewMoviesRepository(conn)
    genresRepository := repositories.NewGenresRepository(conn)
    watchlistRepository := repositories.NewWatchlistRepository(conn)
    usersRepository := repositories.NewUsersRepository(conn)

    moviesHandler := handlers.NewMoviesHandler(genresRepository, moviesRepository)
    genresHandler := handlers.NewGenresHandler(genresRepository)
    watchlistHandler := handlers.NewWatchlistHandler(watchlistRepository)
    usersHandler := handlers.NewUsersHandler(usersRepository)
    authHandler := handlers.NewAuthHandlers(usersRepository)

    imageHandler := handlers.NewImageHandlers()

    authorized := r.Group("")
    authorized.Use(middlewares.AuthMiddleware)

    authorized.GET("/movies", moviesHandler.FindAll)     
    authorized.GET("/movies/:id", moviesHandler.FindById)
    authorized.POST("/movies", moviesHandler.Create)
    authorized.PUT("/movies/:id", moviesHandler.Update)
    authorized.DELETE("/movies/:id", moviesHandler.Delete)
    authorized.PATCH("/movies/:movieId/rate", moviesHandler.SetRating)
    authorized.PATCH("/movies/:movieId/setWatched", moviesHandler.SetWatched)

    authorized.GET("/genres", genresHandler.FindAll)     
    authorized.GET("/genres/:id", genresHandler.FindById)
    authorized.POST("/genres", genresHandler.Create)
    authorized.PUT("/genres/:id", genresHandler.Update)
    authorized.DELETE("/genres/:id", genresHandler.Delete)

    

    authorized.GET("/watchlist", watchlistHandler.FindAll)
    authorized.POST("/watchlist/:movieId", watchlistHandler.AddToWatchlist)
    authorized.DELETE("/watchlist/:movieId", watchlistHandler.Delete)

    authorized.GET("/users", usersHandler.FindAll)
    authorized.GET("/users/:id", usersHandler.FindById)
    authorized.POST("/users", usersHandler.Create)
    authorized.PUT("/users/:id", usersHandler.Update)
    authorized.PATCH("/users/:id/changePassword", usersHandler.ChangePasswordHash)
    authorized.DELETE("/users/:id", usersHandler.Delete)

    authorized.POST("/auth/signOut", authHandler.SignOut)
    authorized.GET("/auth/userInfo", authHandler.GetUserInfo)

    unauthorized := r.Group("")
    unauthorized.POST("/auth/signIn", authHandler.SignIn)

    unauthorized.GET("/images/:imageId", imageHandler.HandleGetImageById)

    docs.SwaggerInfo.BasePath = "/"
    unauthorized.GET("/swagger/*any", swagger.WrapHandler(swaggerfiles.Handler))

    logger.Info("Application starting...")
    
    r.Run(config.Config.AppHost)
}

func loadConfig() error {
    viper.SetConfigFile(".env")
    err := viper.ReadInConfig()
    if err != nil {
        return err
    }

    var mapConfig config.MapConfig
    err = viper.Unmarshal(&mapConfig)
    if err != nil {
        return err
    }

    config.Config = &mapConfig
    return nil
}

func connectToDb() (*pgxpool.Pool, error) {
    conn, err := pgxpool.New(context.Background(), config.Config.DbConnectionString)
    if err != nil{
        return nil, err
    } 
    
    err = conn.Ping(context.Background())
    if err != nil {
        return nil, err
    }
    return conn, nil
}