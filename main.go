package main

import (
	"github.com/Biubiubiuuuu/goDoutu/db/mysql"
	"github.com/Biubiubiuuuu/goDoutu/models"
)

func main() {
	mysql.DB.InitConn()
	db := mysql.GetMysqlDB()
	db.AutoMigrate(
		&models.Emoticons{},
		&models.User{},
		&models.EmoticonsType{},
		&models.EmoticonsGrouping{},
	)
}
