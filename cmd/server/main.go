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

	likeRepo := repositories.NewLikeRepository(db)
	likeService := services.NewLikeService(likeRepo)
	likeHandler := handlers.NewLikeHandler(likeService, zapLogger)
	messageRepo := repositories.NewMessageRepository(db)
	messageService := services.NewMessageService(messageRepo)
	messageHandler := handlers.NewMessageHandler(messageService, zapLogger)

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

	secure.HandleFunc("/api/likes", likeHandler.AddLikeHandler).Methods("POST")
	secure.HandleFunc("/api/likes", likeHandler.RemoveLikeHandler).Methods("DELETE")
	secure.HandleFunc("/api/likes", likeHandler.GetLikesHandler).Methods("GET")
	secure.HandleFunc("/api/likes/count", likeHandler.GetLikeCountHandler).Methods("GET")

	secure.HandleFunc("/api/messages", messageHandler.SendMessage).Methods("POST")
	secure.HandleFunc("/api/messages", messageHandler.GetMessages).Methods("GET")
	port := cfg.ServerPort
	zapLogger.Info("Сервер запущен", zap.String("порт", port))
	log.Fatal(http.ListenAndServe(":"+port, r))
}
