package history_auth

import (
	"database/sql"
	"log"
)

type (
	//Auth описывает модель аватарки
	Auth struct {
		Id     *int64
		UserId *int64
		Time   *string
	}
)

const AllFields = "id, user_id, time"

func (auth *Auth) allocateMem() {
	auth.Id = new(int64)
	auth.UserId = new(int64)
}

var connection *sql.DB

func SetConnection(conn *sql.DB) {
	if conn == nil {
		log.Fatalf("Auth: Connection is nil")
	}
}

func scanAllFields(rows *sql.Rows) (*Auth, error) {
	auth := new(Auth)
	auth.allocateMem()
	err := rows.Scan(
		auth.Id,
		auth.UserId,
		auth.Time,
	)
	if err != nil {
		return nil, err
	}
	return auth, nil
}
