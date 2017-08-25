package main

import (
	"database/sql"
	"fmt"
	"log"
	"messenger_db/models/icon"
	"messenger_db/models/user"

	_ "github.com/go-sql-driver/mysql"
)

func newStr(s string) *string {
	return &s
}

func newIn64(i int64) *int64 {
	return &i
}

func main() {
	//"root:1234@tcp(localhost:3306)/msg?charset=utf8"
	con, err := sql.Open("mysql", "root:1234@tcp(localhost:32772)/msg")
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}
	if con == nil {
		log.Fatalf("User: Connection is nil")
	}
	reply := true

	user := user.User{
		Login:    newStr("maria"),
		Name:     newStr("Maria"),
		LastName: newStr("Terskikh"),
		Password: newStr("strenght"),
	}

	icon := icon.Icon{
		UserId:   newIn64(4),
		UserIcon: newStr("ಠ_ಠ"),
	}

	err = user.AddUser(user, &reply)
	fmt.Println("err =", err)

	err = user.GetUserById(2, user)
	fmt.Println("resp =", *user.Name, *user.LastName)

	err = icon.AddIcon(icon, &reply)
	fmt.Println("err =", err)

}
