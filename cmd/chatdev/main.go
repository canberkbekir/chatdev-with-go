package main

import (
	"chatdev-with-go/internal/config"
	"chatdev-with-go/internal/utils"
)

func main() {

	utils.InitLog()
	utils.Logger.Info().Msg("Starting ChatDev with Go")
	config.NewConfig()
	utils.Logger.Info().Msg("ChatDev with Go started successfully")

}
