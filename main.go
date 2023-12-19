package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/soulteary/yi-openai-proxy/define"
	"github.com/soulteary/yi-openai-proxy/yi"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	version   = ""
	buildDate = ""
	gitCommit = ""
)

func main() {
	viper.AutomaticEnv()
	parseFlag()

	err := yi.GetInstance()
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	registerRoute(r)

	srv := &http.Server{
		Addr:    viper.GetString("listen"),
		Handler: r,
	}

	runServer(srv)
}

// registerRoute registers all routes
func registerRoute(r *gin.Engine) {
	// https://platform.openai.com/docs/api-reference
	r.HEAD("/", func(c *gin.Context) {
		c.Status(200)
	})

	r.Any("/health", func(c *gin.Context) {
		c.Status(200)
	})

	stripPrefixConverter := yi.NewStripPrefixConverter(define.REST_API_VERSION)
	apiRouter := r.Group(define.REST_API_VERSION)
	{
		apiRouter.Any("/completions", yi.ProxyWithConverter(stripPrefixConverter))
		apiRouter.Any("/chat/completions", yi.ProxyWithConverter(stripPrefixConverter))
	}
}

func runServer(srv *http.Server) {
	go func() {
		log.Printf("Server listening at %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(errors.Errorf("listen: %s\n", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server Shutdown...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func parseFlag() {
	pflag.StringP("listen", "l", ":8080", "listen address")
	pflag.BoolP("version", "v", false, "version information")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}
	if viper.GetBool("v") {
		fmt.Println("version:", version)
		fmt.Println("buildDate:", buildDate)
		fmt.Println("gitCommit:", gitCommit)
		os.Exit(0)
	}
}
