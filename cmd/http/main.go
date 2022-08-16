package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"social/internal/app/http"
	"social/pkg/log"
)

func main()  {
	loadConfig()
	log.Init("http", "log", "debug")
	http.Run()
}

func loadConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("load config error:", err)
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("file change: %s, %s\n",e.Name, viper.GetString("websocket.port"))
	})
	log.Info("load config successfully")
}
