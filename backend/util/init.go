package util

import "github.com/jesee-kuya/forum/backend/database"

func Init() {
	database.CreateConnection()
}
