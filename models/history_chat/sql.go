package history_chat

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

//AddMessage добавляет сообщение
func (_ *HistoryChat) AddMessage(historyChat *HistoryChat, ok *bool) error {
	*ok = false

	query := "INSERT INTO history_chat(user_id, chat_id, text) VALUES (?,?,?)"
	err := queryExecutionHandler(query, *historyChat.UserId, *historyChat.ChatId, *historyChat.Text)
	if err != nil {
		return err
	}
	*ok = true
	return nil
}

func (_ *HistoryChat) GetMessageByChat(historyChat *HistoryChat, resp *[]HistoryChat) error {
	con, err := sql.Open("mysql", "root:1234@tcp(localhost:32772)/msg")
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}
	if con == nil {
		log.Fatalf("Auth: Connection is nil")
	}
	rows, err := con.Query(fmt.Sprintf("SELECT id, user_id, time, text FROM history_auth where user_id = ? and chat_id=? ORDER BY time"), *historyChat.UserId, *historyChat.ChatId)
	if err != nil {
		return err
	}
	defer rows.Close()
	hChats := make([]HistoryChat, 0)
	for rows.Next() {
		historyChat2 := new(HistoryChat)
		historyChat2.allocateMem()
		err = rows.Scan(historyChat2.Id, historyChat2.UserId, historyChat2.Time, historyChat2.Text)
		if err != nil {
			return err
		}
		hChats = append(hChats, *historyChat2)
	}
	*resp = hChats
	return nil
}

// Запрос в базу
func queryExecutionHandler(query string, args ...interface{}) error {
	con, err := sql.Open("mysql", "root:1234@tcp(localhost:32772)/msg")
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}
	if con == nil {
		log.Fatalf("User: Connection is nil")
	}
	row, err := con.Exec(query, args...)
	if err != nil {
		return err
	}
	err = rowNumbersHandler(row)
	return err
}

// Проверяет колличество обработаных записей, если не было обработано ни одной - возвращает ошибку noRowsProcessedError, иначе nil.
func rowNumbersHandler(row sql.Result) error {
	noRowsProcessedError := errors.New("Failed to update/create the user. Maybe there is no user with such ID in the database")
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return noRowsProcessedError
	}
	return nil
}
