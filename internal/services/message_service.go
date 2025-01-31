package services

import (
	"InstaSpace/internal/models"
	"InstaSpace/internal/repositories"
	"context"
	"errors"
)

type MessageService struct {
	Repo *repositories.MessageRepository
}

func NewMessageService(repo *repositories.MessageRepository) *MessageService {
	return &MessageService{Repo: repo}
}

// Создать или получить ID переписки
func (s *MessageService) GetOrCreateConversation(ctx context.Context, user1ID, user2ID int) (int, error) {
	return s.Repo.CreateConversation(ctx, user1ID, user2ID)
}

// Определяем ошибку, чтобы использовать её в обработчике
var ErrConversationNotFound = errors.New("conversation not found")

func (s *MessageService) SendMessage(ctx context.Context, conversationID, senderID int, content string) (int, error) {
	// Проверяем, существует ли переписка
	var exists bool
	err := s.Repo.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM conversations WHERE id=$1)", conversationID).Scan(&exists)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, ErrConversationNotFound
	}

	// Добавляем сообщение
	return s.Repo.SendMessage(ctx, conversationID, senderID, content)
}

func (s *MessageService) GetMessages(ctx context.Context, conversationID int) ([]models.Message, error) {
	// Проверяем, существует ли переписка
	var exists bool
	err := s.Repo.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM conversations WHERE id=$1)", conversationID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrConversationNotFound
	}

	// Если переписка найдена, получаем её сообщения
	return s.Repo.GetMessages(ctx, conversationID)
}

func (s *MessageService) DeleteMessage(ctx context.Context, messageID int) error {
	return s.Repo.DeleteMessage(ctx, messageID)
}
