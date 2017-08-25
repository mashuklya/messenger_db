package user

import (
	"database/sql"
	"fmt"
	"log"
)

type (
	//User описывает модель пользователя системы
	User struct {
		Id       *int64
		Login    *string
		Password *string //There is no in db. Use only for check/set password
		Name     *string
		LastName *string
		Icon     *string
	}
)

const AllFields = "id, login, name, last_name"

func (u *User) allocateMem() {
	u.Id = new(int64)
	u.Login = new(string)
	u.Password = new(string)
	u.Name = new(string)
	u.LastName = new(string)
}

var connection *sql.DB

func SetConnection(conn *sql.DB) {
	if conn == nil {
		log.Fatalf("User: Connection is nil")
	}
}

func (user *User) checkFields(fields ...string) (bool, []error) {
	errs := make([]error, 0)
	for _, field := range fields {
		switch field {
		case "password":
			if user.Password == nil {
				errs = append(errs, fmt.Errorf("Set password field"))
			}
		case "login":
			if user.Login == nil {
				errs = append(errs, fmt.Errorf("Set login field"))
			}
		case "name":
			if user.Name == nil {
				errs = append(errs, fmt.Errorf("Set name field"))
			}
		case "last_name":
			if user.LastName == nil {
				errs = append(errs, fmt.Errorf("Set last name field"))
			}
		case "id":
			if user.Id == nil {
				errs = append(errs, fmt.Errorf("Set id field"))
			}
		}
	}
	if len(errs) != 0 {
		return false, errs
	}
	return true, nil
}

func scanAllFields(rows *sql.Rows) (*User, error) {
	user := new(User)
	user.allocateMem()
	err := rows.Scan(user.Id,
		user.Id,
		user.Login,
		user.Name,
		user.LastName,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
