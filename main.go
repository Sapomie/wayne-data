package main

import (
	"fmt"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/routers"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	err := BeforeStarting()
	if err != nil {
		panic(err)
	}

	s := setServer()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = s.ListenAndServe()
		if err != nil {
			fmt.Println("server error")
		}
	}()

	go func() {
		quit := make(chan os.Signal, 2)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
		for {
			select {
			case <-quit:
				wg.Done()
			}
		}
	}()

	wg.Wait()

	err = ShutDown()
	if err != nil {
		panic(err)
	}

}

func setServer() *http.Server {
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	return &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
}

func ShutDown() error {

	err := model.CloseDb()
	if err != nil {
		return err
	}

	fmt.Println("Gracefully shut down")

	return nil
}
