package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword генерирует хэш для заданного пароля.
// Используется для безопасного хранения паролей пользователей.
//
// Возвращает строку, содержащую хэш пароля, и возможную ошибку.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash проверяет соответствие заданного пароля и хэша.
// Используется для аутентификации пользователя при вводе пароля.
//
// Принимает на вход строку пароля и строку хэша.
// Возвращает true, если пароль соответствует хэшу, иначе возвращает false.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
