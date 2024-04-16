package main

import (
	"errors"
	"github.com/spf13/viper"
	"gomc/src/data"
	"gomc/src/server"
	"gomc/src/status"
)

func main() {
	if err := loadConfig(); err != nil {
		panic(err)
	}

	if err := status.Init(); err != nil {
		panic(err)
	}

	data.PrismarineStairs{
		Facing:      data.PrismarineStairsFacingSouth,
		Half:        data.PrismarineStairsHalfBottom,
		Shape:       data.PrismarineStairsShapeInnerLeft,
		Waterlogged: true,
	}.Id()

	srv := server.NewServer(&server.Config{
		ViewDistance: viper.GetInt("view_distance"),
		OnlineMode:   viper.GetBool("online_mode"),
	})
	srv.Start()
}

func loadConfig() error {
	viper.SetDefault("bind_addr", "")
	viper.SetDefault("port", 25565)
	viper.SetDefault("motd", "Hello world!")
	viper.SetDefault("max_players", 100)
	viper.SetDefault("online_mode", true)
	viper.SetDefault("view_distance", 10)

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	var configFileNotFoundError viper.ConfigFileNotFoundError
	if errors.As(err, &configFileNotFoundError) {
		return viper.SafeWriteConfigAs("config.toml")
	}
	return err
}
