package database

import (
	"fmt"
	"os"
	"strings"

	"MVC-golang/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pelletier/go-toml"
)

var (
	DB = New()
)

/**
 * Set up a database connection
 * @param diver string
 */

func New() *gorm.DB {

	// if getAppEnv() == "test" {

	// } else {
	driver := config.Conf.Get("database.driver").(string)
	configTree := config.Conf.Get(driver).(*toml.Tree)

	connect := configTree.Get("connect").(string)

	DB, err := gorm.Open(driver, connect)

	if err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to  database, but got err=%+v", err))
	}

	return DB
	// }
}

// Get the program running environment
// According to the program running path suffix
// If it is test, it is the test environment.
func getAppEnv() string {
	file := os.Args[0]
	s := strings.Split(file, ".")
	return s[len(s)-1]
}
