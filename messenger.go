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
	var resp []icon.Icon
	reply := true

	user := user.User{
		Login:    newStr("vasya"),
		Name:     newStr("Vasya"),
		LastName: newStr("Terskikh"),
		Password: newStr("1234"),
	}

	icon := icon.Icon{
		UserId:   newIn64(19),
		UserIcon: newStr("˙ ͜ʟ˙"),
	}
	err = user.GetUserByLogin(*user.Login, user)
	fmt.Println("err AddUser =", err)

	if user.Login == nil {
		err = user.AddUser(user, &reply)
		fmt.Println("err AddUser =", err)
	}
	err = user.GetUserById(2, user)
	fmt.Println("err GetUserById =", err)
	fmt.Println("resp GetUserById=", *user.Name, *user.LastName)

	err = icon.AddIcon(icon, &reply)
	fmt.Println("err AddIcon =", err)

	err = icon.GetIconByUserId(19, icon)
	fmt.Println("err GetIconByUserId =", err)
	fmt.Println("resp GetIconByUserId=", *icon.UserIcon, *icon.UserId)

	err = icon.GetIcons(19, resp)
	fmt.Println("err GetIcons =", err)
	fmt.Println("resp GetIcons=", resp)

}
