package main

import (
	"log"
	"net/http"

	"InstaSpace/internal/handlers"
	"InstaSpace/internal/repositories"
	"InstaSpace/internal/services"
	"InstaSpace/pkg/config"
	"InstaSpace/pkg/middleware"
	"github.com/gorilla/mux"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Подключаемся к базе данных
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()

	log.Println("Database connection established successfully!")

	// Инициализация слоев
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	// Настройка роутов
	r := mux.NewRouter()
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Защищенные маршруты
	secure := r.PathPrefix("/secure").Subrouter()
	secure.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	secure.HandleFunc("/profile", authHandler.Profile).Methods("GET")

	// Запуск сервера
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
