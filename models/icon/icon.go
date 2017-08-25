package icon

import (
	"database/sql"
	"fmt"
	"log"
)

type (
	//Icon описывает модель аватарки
	Icon struct {
		Id       *int64
		UserId   *int64
		UserIcon *string
	}
)

const AllFields = "id, user_id, user_icon"

func (icon *Icon) allocateMem() {
	icon.Id = new(int64)
	icon.UserId = new(int64)
	icon.UserIcon = new(string)
}

var connection *sql.DB

func SetConnection(conn *sql.DB) {
	if conn == nil {
		log.Fatalf("Icon: Connection is nil")
	}
}
func (icon *Icon) checkFields(fields ...string) (bool, []error) {
	errs := make([]error, 0)
	for _, field := range fields {
		switch field {
		case "user_id":
			if icon.UserId == nil {
				errs = append(errs, fmt.Errorf("Set user id field"))
			}
		case "user_icon":
			if icon.UserIcon == nil {
				errs = append(errs, fmt.Errorf("Set user icon field"))
			}
		}
	}
	if len(errs) != 0 {
		return false, errs
	}
	return true, nil
}

func scanAllFields(rows *sql.Rows) (*Icon, error) {
	icon := new(Icon)
	icon.allocateMem()
	err := rows.Scan(
		icon.Id,
		icon.UserId,
		icon.UserIcon,
	)
	if err != nil {
		return nil, err
	}
	return icon, nil
}
