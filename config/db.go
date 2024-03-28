package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB сущность базы данных
var DB *gorm.DB

// InitPgSQL Инициализация базы данных PgSQL
func InitPgSQL() error {
	var db *gorm.DB

	// создание строки подключения она всегда статична и имеет такое количество и порядок аргументов
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		Env.DbUser,
		Env.DbPassword,
		Env.DbHost,
		Env.DbPort,
		Env.DbName,
	)

	// подключение к бд
	db, err := gorm.Open("postgres", connectionString)

	if err != nil {
		return err
	}

	DB = db

	return nil
}
