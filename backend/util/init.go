package util

import (
	"database/sql"

	"github.com/jesee-kuya/forum/backend/database"
)

var DB *sql.DB

func Init() {
	DB = database.CreateConnection()
}
