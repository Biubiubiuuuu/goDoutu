package main

import (
	"github.com/Biubiubiuuuu/goDoutu/db/mysql"
	"github.com/Biubiubiuuuu/goDoutu/helper/config"
	"github.com/Biubiubiuuuu/goDoutu/models"
	"github.com/Biubiubiuuuu/goDoutu/router"
)

func main() {
	mysql.DB.InitConn()
	db := mysql.GetMysqlDB()
	db.AutoMigrate(
		&models.Emoticons{},
		&models.User{},
		&models.EmoticonsType{},
		&models.EmoticonsGrouping{},
		&models.UserFans{},
		&models.UserFollows{},
	)
	router := router.Init()
	router.Run(config.HTTPPort)
}
