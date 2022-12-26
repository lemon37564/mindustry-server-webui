package main

import (
	"context"
	"errors"
	"mindserver/log"
	"mindserver/mindustrybridge"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./static", false)))

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	// Mindustry subservice
	minServer := mindustrybridge.Route(router.Group("/mindustry"))
	defer minServer.MustExit()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Infof("listen: %s\n", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	log.Info("Recived signal " + sig.String())

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	log.Info("Shutting down the server in 5 seconds")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:" + err.Error())
	}
}
