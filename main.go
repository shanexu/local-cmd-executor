package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Host     string               `toml:"host"`
	Port     int                  `toml:"port"`
	Commands map[string]CmdConfig `toml:"commands"`
}

type CmdConfig struct {
	Cmds []string `toml:"cmds"`
}

func main() {
	viper.SetConfigName("config")                            // name of config file (without extension)
	viper.SetConfigType("toml")                              // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/local-cmd-executor/")          // path to look for the config file in
	viper.AddConfigPath("$HOME/.config/local-cmd-executor/") // call multiple times to add many search paths
	viper.AddConfigPath(".")                                 // optionally look for config in the working directory
	err := viper.ReadInConfig()                              // Find and read the config file
	if err != nil {                                          // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	var config Config
	viper.Unmarshal(&config)
	fmt.Printf("%+v\n", config)
	fmt.Println("hello world")
}
