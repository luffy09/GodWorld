// main.go
package main

import (
	"fmt"
	"godworld/api"
	_ "godworld/docs"
	"godworld/god"

	"github.com/gin-gonic/gin"
	"github.com/pkg/browser"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var world = god.NewWorld()

func main() {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register routes (defined in handlers.go)
	api.RegisterHandlers(r, world)

	go func() {
		url := "http://localhost:8080/swagger/index.html"
		if err := browser.OpenURL(url); err != nil {
			fmt.Println("Failed to open browser:", err)
		}
	}()

	r.Run(":8080")
}
