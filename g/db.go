package g

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open("mysql", Config.DB.Dsn)

	if err != nil {
		fmt.Printf("error connect mysql, error is %v", err)
		os.Exit(-1)
	}
}
