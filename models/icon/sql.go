package icon

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

//AddIcon добавляет нового аватарку пользователя
func (_ *Icon) AddIcon(icon Icon, ok *bool) error {
	var errs []error

	*ok, errs = icon.checkFields("user_id", "user_icon")
	if !*ok {
		return errs[0]
	}
	*ok = false

	query := "INSERT INTO icon(user_id,user_icon) VALUES (?,?)"
	err := queryExecutionHandler(query, *icon.UserId, *icon.UserIcon)
	if err != nil {
		return err
	}
	//}
	*ok = true
	return nil
}

//GetIconById получает аватарку пользователя по его id
func (_ *Icon) GetIconByUserId(userId int64, icon Icon) error {
	con, err := sql.Open("mysql", "root:1234@tcp(localhost:32772)/msg")
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}
	if con == nil {
		log.Fatalf("User: Connection is nil")
	}
	row, err := con.Query(fmt.Sprintf("SELECT user_id, user_icon FROM icon WHERE id = ?"), userId)
	if err != nil {
		return err
	}
	for row.Next() {
		err = row.Scan(&icon.UserId, &icon.UserIcon)
	}
	return nil
}

//GetIcons получает аватарки пользователей
func (_ *Icon) GetIcons(id int64, resp []Icon) error {

	con, err := sql.Open("mysql", "root:1234@tcp(localhost:32772)/msg")
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}
	if con == nil {
		log.Fatalf("User: Connection is nil")
	}
	rows, err := con.Query(fmt.Sprintf("SELECT * FROM icon"))
	if err != nil {
		return err
	}
	defer rows.Close()
	icons := make([]Icon, 0)
	for rows.Next() {
		icon, err := scanAllFields(rows)
		if err != nil {
			return err
		}
		icons = append(icons, *icon)
	}
	resp = icons
	return nil
}

//ChangeIcon получает аватарки пользователей
func (_ *Icon) ChangeIcon(icon Icon) error {

	con, err := sql.Open("mysql", "root:1234@tcp(localhost:32772)/msg")
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}
	if con == nil {
		log.Fatalf("User: Connection is nil")
	}
	stmt, err := con.Prepare("update icon set user_icon=? where user_id=?")
	_, err = stmt.Exec(*icon.UserIcon, *icon.UserId)
	if err != nil {
		return err
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
