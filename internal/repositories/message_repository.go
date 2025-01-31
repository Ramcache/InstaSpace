package repositories

import (
	"InstaSpace/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	DB *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{DB: db}
}

func (r *MessageRepository) CreateConversation(ctx context.Context, user1ID, user2ID int) (int, error) {
	var conversationID int
	err := r.DB.QueryRow(ctx, `
		INSERT INTO conversations (user1_id, user2_id) 
		VALUES ($1, $2) 
		ON CONFLICT (user1_id, user2_id) 
		DO UPDATE SET user1_id = EXCLUDED.user1_id
		RETURNING id`, user1ID, user2ID).Scan(&conversationID)

	if err != nil {
		return 0, err
	}
	return conversationID, nil
}

func (r *MessageRepository) SendMessage(ctx context.Context, conversationID, senderID int, content string) (int, error) {
	var messageID int
	err := r.DB.QueryRow(ctx, `
		INSERT INTO messages (conversation_id, sender_id, content) 
		VALUES ($1, $2, $3) RETURNING id`, conversationID, senderID, content).Scan(&messageID)

	if err != nil {
		return 0, err
	}
	return messageID, nil
}

func (r *MessageRepository) GetMessages(ctx context.Context, conversationID int) ([]models.Message, error) {
	var exists bool
	err := r.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM conversations WHERE id=$1)", conversationID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("conversation not found")
	}

	rows, err := r.DB.Query(ctx, `
		SELECT id, conversation_id, sender_id, content, created_at 
		FROM messages WHERE conversation_id = $1 ORDER BY created_at ASC`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *MessageRepository) DeleteMessage(ctx context.Context, messageID int) error {
	var exists bool
	err := r.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM messages WHERE id=$1)", messageID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("message not found")
	}

	// Удаляем сообщение
	_, err = r.DB.Exec(ctx, "DELETE FROM messages WHERE id = $1", messageID)
	return err
}
