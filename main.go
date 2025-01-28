package main

import (
	"github.com/jesee-kuya/forum/backend/database"
	"github.com/jesee-kuya/forum/backend/repositories"
)

func main() {
	db := database.CreateConnection()
	repositories.RemoveFile(db)
	//repositories.AddFile(db)
}
