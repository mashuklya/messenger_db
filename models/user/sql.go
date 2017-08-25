package user

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
)

//Create добавляет нового пользователя в БД
func (_ *User) AddUser(user User, ok *bool) error {
	var errs []error

	*ok, errs = user.checkFields("password", "login", "email", "name")
	if !*ok {
		return errs[0]
	}
	*ok = false
	salt := GenSalt(*user.Login, *user.Name)
	byteKey := GetByteKey(*user.Password, salt)
	key := base64.StdEncoding.EncodeToString(byteKey)

	query := "INSERT INTO user(login,pwd_key,salt,name,last_name) VALUES (?,?,?,?,?)"
	err := queryExecutionHandler(query, *user.Login, key, salt, *user.Name, *user.LastName)
	if err != nil {
		return err
	}

	*ok = true
	return nil
}

func (_ *User) GetUserById(id int64, resp User) error {
	//resp.allocateMem()
	//TODO: получить иконку
	con, err := sql.Open("mysql", "root:1234@tcp(localhost:32772)/msg")
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}
	if con == nil {
		log.Fatalf("User: Connection is nil")
	}
	row, err := con.Query(fmt.Sprintf("SELECT %s FROM user WHERE id = ?", AllFields), id)
	if err != nil {
		return err
	}
	user, err := scanAllFields(row)
	if err != nil {
		return err
	}
	resp = *user
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
