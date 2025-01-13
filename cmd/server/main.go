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
	photoRepo := repositories.NewPhotoRepository(db)
	commentRepo := repositories.NewCommentRepository(db)

	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	photoService := services.NewPhotoService(photoRepo)
	commentService := services.NewCommentService(commentRepo)

	authHandler := handlers.NewAuthHandler(authService, zapLogger)
	photoHandler := handlers.NewPhotoHandler(photoService, zapLogger)
	commentHandler := handlers.NewCommentHandler(commentService, zapLogger)

	r := mux.NewRouter()
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	secure := r.PathPrefix("/api").Subrouter()
	secure.Use(middleware.JWTMiddleware(cfg.JWTSecret, zapLogger))

	secure.HandleFunc("/photos", photoHandler.UploadPhoto).Methods("POST")

	secure.HandleFunc("/comments", commentHandler.CreateComment).Methods("POST")
	secure.HandleFunc("/comments/{photoID}", commentHandler.GetCommentsByPhotoID).Methods("GET")
	secure.HandleFunc("/comments/{id}/edit", commentHandler.UpdateComment).Methods("PUT")      // Обновление
	secure.HandleFunc("/comments/{id}/delete", commentHandler.DeleteComment).Methods("DELETE") // Удаление

	port := cfg.ServerPort
	zapLogger.Info("Сервер запущен", zap.String("порт", port))
	log.Fatal(http.ListenAndServe(":"+port, r))
}
