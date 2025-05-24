package main

import (
	"fmt"
	_ "godworld/docs"
	"godworld/god"

	"github.com/gin-gonic/gin"

	"github.com/pkg/browser"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title God World API
// @version 1.0
// @description A divine world simulator with unpredictable outcomes.
// @host localhost:8080
// @BasePath /

var world = god.NewWorld()

// createRequest represents the JSON payload for /create
type createRequest struct {
	Name       string            `json:"name" binding:"required"`
	Properties map[string]string `json:"properties"`
}

// @title God World API
// @version 1.0
// @description A divine world simulator with unpredictable outcomes.
// @host localhost:8080
// @BasePath /

// @Summary Create an entity
// @Description Creates a new entity in God World. Chaos may interfere.
// @Tags entities
// @Accept json
// @Produce json
// @Param entity body createRequest true "Entity info"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /create [post]
func createHandler(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	chaosMsg := world.Create(req.Name, req.Properties)

	response := gin.H{
		"status": "creation attempted",
	}
	if chaosMsg != "" {
		response["chaos_msg"] = chaosMsg
	}

	c.JSON(200, response)
}

// @Summary Get an entity
// @Description Retrieves an entity by name, may include chaos message if chaos interfered
// @Tags entities
// @Produce json
// @Param name path string true "Entity Name"
// @Success 200 {object} god.GetEntityResponse
// @Failure 404 {object} map[string]string
// @Router /get/{name} [get]
func getHandler(c *gin.Context) {
	name := c.Param("name")
	entity, ok, chaosMsg := world.Get(name)
	if !ok {
		c.JSON(404, gin.H{"error": "entity not found"})
		return
	}
	resp := god.GetEntityResponse{
		Entity:   entity,
		ChaosMsg: chaosMsg,
	}
	c.JSON(200, resp)
}

// @Summary Destroy an entity
// @Description Removes an entity by name. Chaos may interfere.
// @Tags entities
// @Produce json
// @Param name path string true "Entity Name"
// @Success 200 {object} map[string]string
// @Router /destroy/{name} [delete]
func destroyHandler(c *gin.Context) {
	name := c.Param("name")
	chaosMsg := world.Destroy(name)

	response := gin.H{
		"status": "destroy attempted",
	}
	if chaosMsg != "" {
		response["chaos_msg"] = chaosMsg
	}

	c.JSON(200, response)
}

// @Summary Display world Entities
// @Description Shows all entities currently in the world
// @Tags world
// @Produce json
// @Success 200 {object} map[string]god.Entity
// @Router /dump [get]
func dumpHandler(c *gin.Context) {
	c.JSON(200, world.AllEntities())
}

// @Summary Display world state
// @Description Shows everything about the current world
// @Tags world
// @Produce json
// @Success 200 {object} god.DumpResponse
// @Router /dump/world [get]
func dumpWorldHandler(c *gin.Context) {
	data := world.Dump()
	c.JSON(200, data)
}

func initEndpoints() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/create", createHandler)
	r.GET("/get/:name", getHandler)
	r.DELETE("/destroy/:name", destroyHandler)
	r.GET("/dump", dumpHandler)
	r.GET("/dump/world", dumpWorldHandler)

	r.Run(":8080")
}

func main() {

	initEndpoints()

	url := "http://localhost:8080/swagger/index.html"
	err := browser.OpenURL(url)
	if err != nil {
		fmt.Println("Failed to open browser:", err)
	} else {
		fmt.Println("Opened browser to", url)
	}
}
