package main

import (
	_ "InstaSpace/docs"
	InstaHandlers "InstaSpace/internal/handlers"
	_ "InstaSpace/internal/models"
	"InstaSpace/internal/repositories"
	"InstaSpace/internal/services"
	"InstaSpace/pkg/config"
	"InstaSpace/pkg/logger"
	"InstaSpace/pkg/middleware"
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	sugaredLogger, zapLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Не удалось инициализировать логгер: %v", err)
	}
	defer func(sugaredLogger *zap.Logger) {
		err = sugaredLogger.Sync()
		if err != nil {

		}
	}(sugaredLogger)

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
	likeRepo := repositories.NewLikeRepository(db)
	messageRepo := repositories.NewMessageRepository(db)

	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	photoService := services.NewPhotoService(photoRepo)
	commentService := services.NewCommentService(commentRepo)
	likeService := services.NewLikeService(likeRepo)
	messageService := services.NewMessageService(messageRepo)

	authHandler := InstaHandlers.NewAuthHandler(authService, sugaredLogger)
	photoHandler := InstaHandlers.NewPhotoHandler(photoService, sugaredLogger)
	commentHandler := InstaHandlers.NewCommentHandler(commentService, sugaredLogger)
	likeHandler := InstaHandlers.NewLikeHandler(likeService, sugaredLogger)
	messageHandler := InstaHandlers.NewMessageHandler(messageService, sugaredLogger)
	wsHandler := InstaHandlers.NewWebSocketHandler(sugaredLogger, messageService)

	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	secure := r.PathPrefix("/api").Subrouter()
	secure.Use(middleware.JWTMiddleware(cfg.JWTSecret, sugaredLogger))

	secure.HandleFunc("/photos", photoHandler.UploadPhoto).Methods("POST")

	secure.HandleFunc("/comments", commentHandler.CreateComment).Methods("POST")
	secure.HandleFunc("/comments/{photoID}", commentHandler.GetCommentsByPhotoID).Methods("GET")
	secure.HandleFunc("/comments/{id}/edit", commentHandler.UpdateComment).Methods("PUT")
	secure.HandleFunc("/comments/{id}/delete", commentHandler.DeleteComment).Methods("DELETE")

	secure.HandleFunc("/likes", likeHandler.AddLikeHandler).Methods("POST")
	secure.HandleFunc("/likes", likeHandler.RemoveLikeHandler).Methods("DELETE")
	secure.HandleFunc("/likes/users", likeHandler.GetLikesHandler).Methods("GET")
	secure.HandleFunc("/likes/count", likeHandler.GetLikeCountHandler).Methods("GET")

	secure.HandleFunc("/messages", messageHandler.SendMessage).Methods("POST")
	secure.HandleFunc("/messages/{conversationID}", messageHandler.GetMessages).Methods("GET")
	secure.HandleFunc("/messages/{messageID}", messageHandler.DeleteMessageHandler).Methods("DELETE")

	r.HandleFunc("/ws", wsHandler.HandleWS).Methods("GET")

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Разрешаем все источники
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
	)

	port := cfg.ServerPort
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      corsMiddleware(r),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		zapLogger.Info("Сервер запущен", zap.String("порт", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Ошибка сервера", zap.Error(err))
		}
	}()

	<-stop
	zapLogger.Info("Остановка сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		zapLogger.Fatal("Ошибка завершения работы сервера", zap.Error(err))
	}

	zapLogger.Info("Сервер успешно остановлен")
}
