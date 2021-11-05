package main

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/routers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	err := BeforeStarting()
	if err != nil {
		panic(err)
	}

	gin.SetMode(global.ServerSetting.RunMode)

	r := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        r,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	err = s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
