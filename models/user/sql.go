package user

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

//Create добавляет нового пользователя в БД
func (_ *User) AddUser(user User, ok *bool) error {
	var errs []error

	*ok, errs = user.checkFields("password", "login", "email", "name")
	if !*ok {
		return errs[0]
	}
	*ok = false
	con, err := sql.Open("mysql", "root:1234@tcp(localhost:32772)/msg")
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}
	if con == nil {
		log.Fatalf("User: Connection is nil")
	}
	row, err := con.Query(fmt.Sprintf("SELECT id,login,name,last_name FROM user WHERE login = ?"), *user.Login)
	if err != nil {
		return err
	}
	for row.Next() {
		err = row.Scan(&user.Id, &user.Login, &user.Name, &user.LastName)
	}
	//смотрим, есть ли в базе юзер с таким логином
	if user.Id != nil {
		return nil
	} else {
		salt := GenSalt(*user.Login, *user.Name)
		byteKey := GetByteKey(*user.Password, salt)
		key := base64.StdEncoding.EncodeToString(byteKey)

		query := "INSERT INTO user(login,pwd_key,salt,name,last_name) VALUES (?,?,?,?,?)"
		err = queryExecutionHandler(query, *user.Login, key, salt, *user.Name, *user.LastName)
		if err != nil {
			return err
		}

		*ok = true
		return nil
	}
}

func (_ *User) GetUserById(id int64, user User) error {
	//resp.allocateMem()
	//TODO: получить иконку
	con, err := sql.Open("mysql", "root:1234@tcp(localhost:32772)/msg")
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}
	if con == nil {
		log.Fatalf("User: Connection is nil")
	}
	row, err := con.Query(fmt.Sprintf("SELECT id, login, name, last_name FROM user WHERE id = ?"), id)
	if err != nil {
		return err
	}
	for row.Next() {
		err = row.Scan(&user.Id, &user.Login, &user.Name, &user.LastName)
	}
	row, err = con.Query(fmt.Sprintf("SELECT user_icon FROM user WHERE user_id = ?"), id)
	if err != nil {
		return err
	}
	for row.Next() {
		err = row.Scan(&user.Icon)
	}
	return nil
}

func (_ *User) GetUserByLogin(login string, user User) error {
	//resp.allocateMem()
	//TODO: получить иконку
	con, err := sql.Open("mysql", "root:1234@tcp(localhost:32772)/msg")
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}
	if con == nil {
		log.Fatalf("User: Connection is nil")
	}
	row, err := con.Query(fmt.Sprintf("SELECT id, login, name, last_name FROM user WHERE login = ?"), login)
	if err != nil {
		return err
	}
	for row.Next() {
		err = row.Scan(&user.Id, &user.Login, &user.Name, &user.LastName)
	}
	//resp = *user
	return nil
}
func (_ *User) GetUsersByChatId(chats []int64, users *[]int64) error {

	var (
		buffer bytes.Buffer
		str    []string
	)
	con, err := sql.Open("mysql", "root:1234@tcp(localhost:32772)/msg")
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}
	if con == nil {
		log.Fatalf("User: Connection is nil")
	}

	for _, id := range chats {
		str = strconv.FormatInt(id, 10)
	}

	ids := strings.Join(str, ", ")

	rows, err := con.Query(fmt.Sprintf("SELECT user_id from history_chat where chat_id = ?"), ids)
	if err != nil {
		return err
	}
	defer rows.Close()
	usersId := make([]int64, 0)
	for rows.Next() {
		var userid int64
		err = rows.Scan(&userid)
		if err != nil {
			return err
		}
		usersId = append(usersId, userid)
	}
	*users = usersId
	//resp = *user
	return nil
}

func (_ *User) CheckPassword(user User, ok *bool) error {
	isPwd, errs := user.checkFields("password")
	if !isPwd {
		return errs[0]
	}
	isId, _ := user.checkFields("id")
	isLogin, _ := user.checkFields("login")
	var where string
	if isId {
		where = fmt.Sprintf("WHERE id='%s'", *user.Id)
	} else if isLogin {
		where = fmt.Sprintf("WHERE login='%s'", *user.Login)
	} else {
		return fmt.Errorf("Set Id or login")
	}
	var key, salt string

	con, err := sql.Open("mysql", "root:1234@tcp(localhost:32772)/msg")
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}
	if con == nil {
		log.Fatalf("User: Connection is nil")
	}
	err = con.QueryRow(fmt.Sprintf("SELECT pwd_key, salt FROM user %s", where)).Scan(&key, &salt)
	if err != nil {
		return err
	}
	if !CheckPassword(*user.Password, salt, key) {
		return fmt.Errorf("Wrong password!")
	}
	*ok = true
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
