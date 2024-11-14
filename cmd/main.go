package main

import (
	_ "InstaSpace/docs"
	"github.com/joho/godotenv"
	"log"
	"net/http"

	"InstaSpace/internal/auth"
	"InstaSpace/internal/db"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
)

// @title Insta Space API
// @version 1.0
// @description Документация для API InstaSpace, включающая функционал авторизации и регистрации.

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @contact.name Техническая поддержка
// @contact.email ramaro@internet.ru

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbPool := db.ConnectionDB()
	defer dbPool.Close()

	authRepo := auth.NewAuthRepository(dbPool)
	authService := auth.NewAuthService(authRepo)
	authHandler := auth.NewAuthHandler(authService)

	router := mux.NewRouter()
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/confirm", authHandler.ConfirmEmail).Methods("GET")
	router.HandleFunc("/resend-confirmation", authHandler.ResendConfirmation).Methods("POST")
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Группа защищённых маршрутов
	//protected := router.PathPrefix("/protected").Subrouter()
	//protected.Use(middleware.JWTAuth)
	//protected.HandleFunc("/test", test).Methods("POST")
	//protected.HandleFunc("/test", test).Methods("GET")

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
