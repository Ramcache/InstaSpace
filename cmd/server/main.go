package main

import (
	"InstaSpace/pkg/logger"
	"go.uber.org/zap"
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
	cfg := config.LoadConfig()

	zapLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Не удалось инициализировать логгер: %v", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("Запуск приложения")

	db, err := config.ConnectDB(cfg)
	if err != nil {
		zapLogger.Fatal("Ошибка подключения к базе данных", zap.Error(err))
	}
	defer db.Close()

	zapLogger.Info("Подключение к базе данных успешно установлено")

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService, zapLogger)

	photoRepo := repositories.NewPhotoRepository(db)
	photoService := services.NewPhotoService(photoRepo)
	photoHandler := handlers.NewPhotoHandler(photoService, zapLogger)

	r := mux.NewRouter()
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	secure := r.PathPrefix("/api").Subrouter()
	secure.Use(middleware.JWTMiddleware(cfg.JWTSecret, zapLogger))
	secure.HandleFunc("/photos", photoHandler.UploadPhoto).Methods("POST")

	zapLogger.Info("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
