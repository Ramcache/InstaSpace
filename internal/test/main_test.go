package test

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"InstaSpace/internal/handlers"
	"InstaSpace/internal/repositories"
	"InstaSpace/internal/services"
	"InstaSpace/pkg/config"
	"InstaSpace/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	db         *pgxpool.Pool
	zapLogger  *zap.Logger
	testServer *httptest.Server
)

func TestMain(m *testing.M) {
	var err error

	cfg := config.LoadConfig()

	zapLogger, err = logger.NewLogger()
	if err != nil {
		log.Fatalf("Не удалось инициализировать логгер: %v", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("Инициализация тестов")

	if err := godotenv.Load("../../.env"); err != nil {
		zapLogger.Warn("Не удалось загрузить файл .env", zap.Error(err))
	}

	db, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		zapLogger.Fatal("Не удалось подключиться к базе данных", zap.Error(err))
	}
	defer db.Close()

	zapLogger.Info("Успешное подключение к базе данных для тестов")

	r := mux.NewRouter()

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService, zapLogger)

	photoRepo := repositories.NewPhotoRepository(db)
	photoService = services.NewPhotoService(photoRepo)
	photoHandler := handlers.NewPhotoHandler(photoService, zapLogger)

	commentRepo := repositories.NewCommentRepository(db)
	commentService = services.NewCommentService(commentRepo)
	commentHandler := handlers.NewCommentHandler(commentService, zapLogger)

	likeRepo := repositories.NewLikeRepository(db)
	likeService := services.NewLikeService(likeRepo)
	likeHandler := handlers.NewLikeHandler(likeService, zapLogger)

	messageRepo := repositories.NewMessageRepository(db)
	messageService := services.NewMessageService(messageRepo)
	messageHandler := handlers.NewMessageHandler(messageService, zapLogger)

	wsHandler = handlers.NewWebSocketHandler(zapLogger, messageService)
	r.HandleFunc("/ws", wsHandler.HandleWS).Methods("GET")

	r.HandleFunc("/api/messages", messageHandler.SendMessage).Methods("POST")
	r.HandleFunc("/api/messages/{conversationID}", messageHandler.GetMessages).Methods("GET")
	r.HandleFunc("/api/messages/{messageID}", messageHandler.DeleteMessageHandler).Methods("DELETE")

	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	r.HandleFunc("/api/photos", photoHandler.UploadPhoto).Methods("POST")

	r.HandleFunc("/api/comments", commentHandler.CreateComment).Methods("POST")
	r.HandleFunc("/api/comments/{photoID}", commentHandler.GetCommentsByPhotoID).Methods("GET")
	r.HandleFunc("/api/comments/{id}/edit", commentHandler.UpdateComment).Methods("PUT")
	r.HandleFunc("/api/comments/{id}/delete", commentHandler.DeleteComment).Methods("DELETE")

	r.HandleFunc("/api/likes", likeHandler.AddLikeHandler).Methods("POST")
	r.HandleFunc("/api/likes", likeHandler.RemoveLikeHandler).Methods("DELETE")
	r.HandleFunc("/api/likes", likeHandler.GetLikesHandler).Methods("GET")
	r.HandleFunc("/api/likes/count", likeHandler.GetLikeCountHandler).Methods("GET")

	testServer = httptest.NewServer(r)
	defer testServer.Close()

	m.Run()
}
