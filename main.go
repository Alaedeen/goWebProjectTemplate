package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Alaedeen/goWebProjectTemplate/config"
	handlers "github.com/Alaedeen/goWebProjectTemplate/handlers"
	models "github.com/Alaedeen/goWebProjectTemplate/models"
	"github.com/Alaedeen/goWebProjectTemplate/repository"
	router "github.com/Alaedeen/goWebProjectTemplate/router"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"accept", "Authorization", "Content-Type"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "HEAD"},
		// Enable Debugging for testing, consider disablin in production
		Debug: true,
	})
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	var configuration config.Configuration

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	//connect to the data base
	UserName := configuration.Database.UserName
	Password := configuration.Database.Password
	DataBase := configuration.Database.DataBase
	Charset := configuration.Database.Charset
	ParseTime := configuration.Database.ParseTime
	dsn := UserName + ":" + Password + "@/" + DataBase + "?charset=" + Charset + "&parseTime=" + ParseTime + "&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	db.AutoMigrate(&models.User{})

	userRepo := repository.UserRepo{db}
	userHandler := handlers.UserHandler{&userRepo}
	// Init Router
	r := mux.NewRouter()
	// serve static files
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	UserRouterHandler := router.UserRouterHandler{Router: r, Handler: userHandler}
	UserRouterHandler.HandleFunctions()
	// start server
	port := ":" + strconv.Itoa(configuration.Server.Port)
	handler := c.Handler(r)

	log.Fatal(http.ListenAndServe(port, handler))
}
