package main

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Config struct {
	Host     string               `toml:"host"`
	Port     int                  `toml:"port"`
	Commands map[string]CmdConfig `toml:"commands"`
}

type CmdConfig struct {
	Cmd []string `toml:"cmds"`
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
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("fatal error unmarshal config: %w", err))
	}
	route := gin.Default()
	route.GET("/:name", func(c *gin.Context) {
		name := c.Param("name")
		cl, ok := config.Commands[name]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "not found",
			})
			return
		}
		stdout, stderr, err := executeCmd(c, cl.Cmd)
		slog.Info("executeCmd", "cl", cl, "stdout", stdout, "stderr", stderr, "err", err)
		var errStr string
		if err != nil {
			errStr = err.Error()
		}
		c.JSON(http.StatusOK, gin.H{
			"stdout": stdout,
			"stderr": stderr,
			"error":  errStr,
		})
	})
	route.Run(fmt.Sprintf("%s:%d", config.Host, config.Port))
}

func executeCmd(ctx context.Context, cl []string) (stdout string, stderr string, err error) {
	cmd := exec.CommandContext(ctx, cl[0], cl[1:]...)
	var o, e bytes.Buffer
	cmd.Stdout = &o
	cmd.Stderr = &e
	err = cmd.Run()
	stdout = o.String()
	stderr = e.String()
	return
}
