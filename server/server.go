package server

import (
	"faceit/appdata"
	"faceit/router"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var routerLoaders = []func(e *gin.Engine){
	router.LoadUsers,
	router.LoadHealthCheck,
}

func Run() {
	e := gin.Default()

	for _, loader := range routerLoaders {
		loader(e)
	}

	err := e.Run(appdata.Config.Api.ListenHost)

	if err != nil {
		log.Fatal().Err(err)
	}
}
