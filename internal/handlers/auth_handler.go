package handlers

import (
	"encoding/json"
	"net/http"

	"InstaSpace/internal/models"
	"InstaSpace/internal/services"

	"go.uber.org/zap"
)

type AuthHandler struct {
	Service services.AuthServiceInterface
	Logger  *zap.Logger
}

func NewAuthHandler(service services.AuthServiceInterface, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		Service: service,
		Logger:  logger,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Начало регистрации пользователя")

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.Logger.Warn("Некорректный ввод при регистрации", zap.Error(err))
		http.Error(w, "Некорректный ввод", http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		h.Logger.Warn("Username отсутствует при регистрации")
		http.Error(w, "Имя пользователя обязательно", http.StatusBadRequest)
		return
	}

	h.Logger.Info("Попытка регистрации пользователя", zap.String("email", user.Email))
	if err := h.Service.RegisterUser(&user); err != nil {
		h.Logger.Error("Ошибка при регистрации пользователя", zap.String("email", user.Email), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Регистрация успешна", zap.String("email", user.Email), zap.String("username", user.Username))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Успешная регистрация. Пожалуйста подтвердите email"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Начало аутентификации пользователя")

	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		h.Logger.Warn("Некорректный ввод при логине", zap.Error(err))
		http.Error(w, "Некорректный ввод", http.StatusBadRequest)
		return
	}

	h.Logger.Info("Попытка аутентификации пользователя", zap.String("email", creds.Email))
	user, err := h.Service.Authenticate(creds.Email, creds.Password)
	if creds.Password == "" {
		h.Logger.Warn("Пустой пароль при логине", zap.String("email", creds.Email))
		http.Error(w, "Пароль не может быть пустым", http.StatusBadRequest)
		return
	}

	if err != nil {
		h.Logger.Warn("Ошибка аутентификации", zap.String("email", creds.Email), zap.Error(err))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := h.Service.GenerateToken(user)
	if err != nil {
		h.Logger.Error("Ошибка генерации токена", zap.String("email", user.Email), zap.Error(err))
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Аутентификация успешна", zap.String("email", user.Email), zap.String("username", user.Username))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token, "username": user.Username})
}

func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Запрос к защищенному маршруту (Profile)")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "This is a protected route!"})
}
