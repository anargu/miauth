package miauth

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB
var dialect = "postgres"
// host=myhost port=myport user=gorm dbname=gorm password=mypassword
var postgresArgs string

func InitDB() {
	postgresArgs = Config.DB.Postgres
	var err error
	DB, err = gorm.Open(dialect, postgresArgs)
	if err != nil {
		panic("failed to connect database")
	}
	RunMigration()
}

func CloseDB() {
	err := DB.Close()
	if err != nil {
		panic(err)
	}
}


var Tables = []interface{}{
	User{},
	Session{},
	Role{},
	LoginCredential{},
	MiauthLoginCredential{},
	GoogleLoginCredential{},
	FacebookLoginCredential{},
}

var roles = []string{
	"user",
}

func RunMigration() {
	if err := DB.AutoMigrate(Tables...).Error; err != nil {
		panic(fmt.Sprintf("No error should happen when create table, but got %+v", err))
	}
	for _, roleName := range roles {
		role := Role{Name: roleName}
		DB.Where(role).FirstOrCreate(&role)
	}
}
