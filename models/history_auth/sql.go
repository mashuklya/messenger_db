package history_auth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

//AddAuth добавляет нового аватарку пользователя
func (_ *Auth) AddAuth(auth Auth, ok *bool) error {
	*ok = false

	query := "INSERT INTO auth(user_id) VALUES (?)"
	err := queryExecutionHandler(query, *auth.UserId)
	if err != nil {
		return err
	}
	*ok = true
	return nil
}

//AddAuth добавляет нового аватарку пользователя
func (_ *Auth) GetAuth(auth Auth, resp *[]Auth) error {
	con, err := sql.Open("mysql", "root:1234@tcp(localhost:32772)/msg")
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}
	if con == nil {
		log.Fatalf("Auth: Connection is nil")
	}
	rows, err := con.Query(fmt.Sprintf("SELECT id, user_id, time FROM history_auth where user_id = ?"), *auth.UserId)
	if err != nil {
		return err
	}
	defer rows.Close()
	auths := make([]Auth, 0)
	for rows.Next() {
		err = rows.Scan(auth.Id, auth.UserId, auth.Time)
		if err != nil {
			return err
		}
		auths = append(auths, auth)
	}
	*resp = auths
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
