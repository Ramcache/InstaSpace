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

	r.HandleFunc("/api/photos", photoHandler.UploadPhoto).Methods("POST")
	testServer = httptest.NewServer(r)
	defer testServer.Close()
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	testServer = httptest.NewServer(r)
	defer testServer.Close()

	m.Run()
}
