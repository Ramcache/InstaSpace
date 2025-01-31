package services

import (
	"InstaSpace/internal/models"
	"InstaSpace/internal/repositories"
	"context"
	"errors"
)

type MessageServiceInterface interface {
	GetOrCreateConversation(ctx context.Context, user1ID, user2ID int) (int, error)
	SendMessage(ctx context.Context, conversationID, senderID int, content string) (int, error)
	GetMessages(ctx context.Context, conversationID int) ([]models.Message, error)
	DeleteMessage(ctx context.Context, messageID int) error
}

type MessageService struct {
	Repo repositories.MessageRepositoryInterface
}

func NewMessageService(repo repositories.MessageRepositoryInterface) *MessageService {
	return &MessageService{Repo: repo}
}

func (s *MessageService) GetOrCreateConversation(ctx context.Context, user1ID, user2ID int) (int, error) {
	return s.Repo.CreateConversation(ctx, user1ID, user2ID)
}

var ErrConversationNotFound = errors.New("conversation not found")

func (s *MessageService) SendMessage(ctx context.Context, conversationID, senderID int, content string) (int, error) {
	var exists bool
	err := s.Repo.ConversationExists(ctx, conversationID, &exists)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, ErrConversationNotFound
	}

	return s.Repo.SendMessage(ctx, conversationID, senderID, content)
}

func (s *MessageService) GetMessages(ctx context.Context, conversationID int) ([]models.Message, error) {
	var exists bool
	err := s.Repo.ConversationExists(ctx, conversationID, &exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrConversationNotFound
	}

	return s.Repo.GetMessages(ctx, conversationID)
}

func (s *MessageService) DeleteMessage(ctx context.Context, messageID int) error {
	return s.Repo.DeleteMessage(ctx, messageID)
}
