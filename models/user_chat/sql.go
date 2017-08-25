package user_chat

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

//AddChat добавляет чат
func (_ *UserChat) AddChat(userChat *UserChat, ok *bool) error {
	*ok = false

	query := "INSERT INTO chats(notes) VALUES (?)"
	err := queryExecutionHandler(query, "new chat")
	if err != nil {
		return err
	}
	query = "INSERT INTO user_chat(user_id, chat_id) VALUES (?,?,?)"
	err = queryExecutionHandler(query, *userChat.UserId, *userChat.ChatId)
	if err != nil {
		return err
	}
	*ok = true

	return nil
}

//вывод всех чатов юзера
func (_ *UserChat) GetChats(user_id string, resp *[]UserChat) error {
	con, err := sql.Open("mysql", "root:1234@tcp(localhost:32772)/msg")
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}
	if con == nil {
		log.Fatalf("Auth: Connection is nil")
	}
	rows, err := con.Query(fmt.Sprintf("SELECT chat_id from user_chat where user_id = ?"), user_id)
	if err != nil {
		return err
	}
	defer rows.Close()
	chatsId := make([]int64, 0)
	for rows.Next() {
		var chat int64
		err = rows.Scan(&chat)
		if err != nil {
			return err
		}
		chatsId = append(chatsId, chat)
	}

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
