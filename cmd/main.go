package main

import (
	_ "InstaSpace/docs"
	"InstaSpace/internal/auth"
	"InstaSpace/internal/db"
	"InstaSpace/internal/photo"
	"InstaSpace/pkg/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load("R:/ProjectsGo/InstaSpace/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbPool := db.ConnectionDB()
	defer dbPool.Close()

	authRepo := auth.NewAuthRepository(dbPool)
	authService := auth.NewAuthService(authRepo)
	authHandler := auth.NewAuthHandler(authService)

	photoRepo := photo.NewRepository(dbPool)
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}
	photoService := photo.NewService(photoRepo, uploadDir)

	baseURL := "http://localhost:8080"
	photoHandler := photo.NewHandler(photoService, baseURL)

	router := mux.NewRouter()
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/confirm", authHandler.ConfirmEmail).Methods("GET")
	router.HandleFunc("/resend-confirmation", authHandler.ResendConfirmation).Methods("POST")
	router.HandleFunc("/photos/{id}", photoHandler.GetPhoto).Methods("GET")

	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	//Защищенные маршруты
	router.Handle("/upload", middleware.JWTAuth(http.HandlerFunc(photoHandler.UploadPhoto))).Methods("POST")
	//

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
