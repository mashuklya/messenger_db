package user_chat

import (
	"database/sql"
	"log"
)

type (
	UserChat struct {
		Id     *int64
		UserId *int64
		ChatId *int64
	}
)

const AllFields = "id, user_id, time"

func (userChat *UserChat) allocateMem() {
	userChat.Id = new(int64)
	userChat.UserId = new(int64)
	userChat.ChatId = new(int64)
}

var connection *sql.DB

func SetConnection(conn *sql.DB) {
	if conn == nil {
		log.Fatalf("UserChat: Connection is nil")
	}
}

func scanAllFields(rows *sql.Rows) (*UserChat, error) {
	userChat := new(UserChat)
	userChat.allocateMem()
	err := rows.Scan(
		userChat.Id,
		userChat.UserId,
		userChat.ChatId,
	)
	if err != nil {
		return nil, err
	}
	return userChat, nil
}
