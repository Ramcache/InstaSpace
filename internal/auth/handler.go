package auth

import (
	_ "InstaSpace/docs"
	"InstaSpace/internal/jwt"
	"encoding/json"
	_ "github.com/swaggo/http-swagger"
	"net/http"
)

// AuthHandler представляет обработчик для авторизации
type AuthHandler struct {
	service Service
}

// NewAuthHandler создает новый обработчик для авторизации
func NewAuthHandler(service Service) *AuthHandler {
	return &AuthHandler{service: service}
}

// RegisterRequest представляет данные для регистрации пользователя
type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Username string `json:"username" example:"john_doe"`
	Password string `json:"password" example:"password123"`
}

// LoginRequest представляет структуру данных для входа
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

// ErrorResponse используется для возвращения сообщений об ошибках
type ErrorResponse struct {
	Message string `json:"message"`
}
type ResendConfirmationRequest struct {
	Email string `json:"email"`
}

// ResendConfirmation godoc
// @Summary Повторная отправка письма с подтверждением
// @Description Этот эндпоинт позволяет повторно отправить письмо для подтверждения почты.
// @Tags auth
// @Accept json
// @Produce json
// @Param data body ResendConfirmationRequest true "Email для повторной отправки подтверждения"
// @Success 200 {string} string "Письмо с подтверждением отправлено повторно"
// @Failure 400 {string} string "Некорректные данные"
// @Failure 404 {string} string "Пользователь не найден"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /resend-confirmation [post]
func (h *AuthHandler) ResendConfirmation(w http.ResponseWriter, r *http.Request) {
	var req ResendConfirmationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	// Попытка повторной отправки письма
	err := h.service.ResendConfirmationEmail(req.Email)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
		} else if err.Error() == "already active" {
			http.Error(w, "Пользователь уже активирован", http.StatusBadRequest)
		} else {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Письмо с подтверждением отправлено повторно"})
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Этот эндпоинт позволяет зарегистрировать нового пользователя, используя email, имя пользователя и пароль.
// @Tags auth
// @Accept json
// @Produce json
// @Param data body RegisterRequest true "Данные для регистрации пользователя"
// @Success 200 {object} map[string]string "message" "Registration successful! Please check your email to confirm."
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 500 {object} ErrorResponse "Ошибка регистрации"
// @Router /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid input"})
		return
	}

	_, err := h.service.RegisterUser(req.Email, req.Username, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Failed to register"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful! Please check your email to confirm."})
}

// ConfirmEmail godoc
// @Summary Подтверждение email пользователя
// @Description Подтверждает email пользователя на основе предоставленного токена.
// @Tags auth
// @Accept  json
// @Produce json
// @Param   token query string true "Токен подтверждения"
// @Success 200 {string} string "Email успешно подтверждён!"
// @Failure 400 {object} ErrorResponse "Token is missing"
// @Failure 401 {object} ErrorResponse "Invalid or expired token"
// @Failure 500 {object} ErrorResponse "Failed to activate user"
// @Router /confirm [get]
func (h *AuthHandler) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is missing", http.StatusBadRequest)
		return
	}

	// Проверяем токен и получаем email
	email, err := jwt.ValidateConfirmationToken(token) // Функция для проверки и извлечения email из токена
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Активируем пользователя (обновляем статус в базе данных)
	err = h.service.ActivateUser(email)
	if err != nil {
		http.Error(w, "Failed to activate user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Email успешно подтверждён!"))
}

// Login godoc
// @Summary Авторизация пользователя
// @Description Этот эндпоинт позволяет пользователю авторизоваться, используя email и пароль.
// @Tags auth
// @Accept json
// @Produce json
// @Param data body LoginRequest true "Данные для входа пользователя"
// @Success 200 {object} map[string]string "token" "Токен успешно создан"
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 401 {object} ErrorResponse "Неверные учетные данные"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid input"})
		return
	}

	token, err := h.service.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid credentials"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
