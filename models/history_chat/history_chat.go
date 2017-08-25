package history_chat

import (
	"database/sql"
	"log"
)

type (
	//HistoryChat описывает модель аватарки
	HistoryChat struct {
		Id     *int64
		UserId *int64
		ChatId *int64
		Text   *string
		Time   *string
	}
)

const AllFields = "id, user_id, time"

func (historyChat *HistoryChat) allocateMem() {
	historyChat.Id = new(int64)
	historyChat.UserId = new(int64)
	historyChat.ChatId = new(int64)
	historyChat.Text = new(string)
	historyChat.Time = new(string)
}

var connection *sql.DB

func SetConnection(conn *sql.DB) {
	if conn == nil {
		log.Fatalf("HistoryChat: Connection is nil")
	}
}

func scanAllFields(rows *sql.Rows) (*HistoryChat, error) {
	historyChat := new(HistoryChat)
	historyChat.allocateMem()
	err := rows.Scan(
		historyChat.Id,
		historyChat.UserId,
		historyChat.ChatId,
		historyChat.Time,
		historyChat.Text,
	)
	if err != nil {
		return nil, err
	}
	return historyChat, nil
}
