package config

import (
	"fmt"

	"github.com/pelletier/go-toml"
)

var (
	Conf = New()
)

/**
 * Return to a singleton instance
 * @method New
 */
func New() *toml.Tree {
	config, err := toml.LoadFile("./config/config.toml")

	if err != nil {
		fmt.Println("TomlErr ", err.Error())
	}

	return config
}
