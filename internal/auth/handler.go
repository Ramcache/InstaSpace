package auth

import (
	_ "InstaSpace/docs" // Пакет сгенерированной документации
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

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Этот эндпоинт позволяет зарегистрировать нового пользователя, используя email, имя пользователя и пароль.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Токен успешно создан"
// @Failure 400 {object} ErrorResponse "Некорректные данные"
// @Failure 500 {object} ErrorResponse "Ошибка регистрации"
// @Router /register [post]

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	token, err := h.service.RegisterUser(req.Email, req.Username, req.Password)
	if err != nil {
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// Login godoc
// @Summary Авторизация пользователя
// @Description Этот эндпоинт позволяет пользователю авторизоваться, используя email и пароль.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Токен успешно создан"
// @Failure 400 {string} string "Некорректные данные"
// @Failure 401 {string} string "Неверные учетные данные"
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	token, err := h.service.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
