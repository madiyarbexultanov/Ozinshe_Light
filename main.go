package main

import (
	"context"
	"goozinshe/config"
	"goozinshe/handlers"
	"goozinshe/repositories"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func main() {
	r := gin.Default()

    corsConfig := cors.Config{
        AllowAllOrigins:    true,
        AllowHeaders:       []string{"*"},
        AllowMethods:       []string{"*"},
    }

    r.Use(cors.New(corsConfig))

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

    imageHandler := handlers.NewImageHandlers()

    r.GET("/movies", moviesHandler.FindAll)     
    r.GET("/movies/:id", moviesHandler.FindById)
    r.POST("/movies", moviesHandler.Create)
    r.PUT("/movies/:id", moviesHandler.Update)
    r.DELETE("/movies/:id", moviesHandler.Delete)
    r.PATCH("/movies/:movieId/rate", moviesHandler.SetRating)
    r.PATCH("/movies/:movieId/setWatched", moviesHandler.SetWatched)

    r.GET("/genres", genresHandler.FindAll)     
    r.GET("/genres/:id", genresHandler.FindById)
    r.POST("/genres", genresHandler.Create)
    r.PUT("/genres/:id", genresHandler.Update)
    r.DELETE("/genres/:id", genresHandler.Delete)

    r.GET("/images/:imageId", imageHandler.HandleGetImageById)

    r.GET("/watchlist", watchlistHandler.FindAll)
    r.POST("/watchlist/:movieId", watchlistHandler.AddToWatchlist)
    r.DELETE("/watchlist/:movieId", watchlistHandler.Delete)

    r.GET("/users", usersHandler.FindAll)
    r.GET("/users/:id", usersHandler.FindById)
    r.POST("/users", usersHandler.Create)
    r.PUT("/users/:id", usersHandler.Update)
    r.PATCH("/users/:id/changePassword", usersHandler.ChangePasswordHash)
    r.DELETE("/users/:id", usersHandler.Delete)
    
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