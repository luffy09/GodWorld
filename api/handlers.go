package api

import (
	"godworld/god"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateRequest represents the JSON payload for /create
type CreateRequest struct {
	Name       string            `json:"name" binding:"required"`
	Properties map[string]string `json:"properties"`
}
type DumpResponse struct {
	Entities map[string]map[string]string // map of entity name → properties
	Events   []string
}
type GetEntityResponse struct {
	Entity   god.Entity `json:"entity"`
	ChaosMsg string     `json:"chaos_msg,omitempty"`
}

// registerHandlers registers all routes and handlers on the gin router.
func RegisterHandlers(r *gin.Engine, world *god.World) {
	r.POST("/create", createHandler(world))
	r.GET("/get/:name", getHandler(world))
	r.DELETE("/destroy/:name", destroyHandler(world))
	r.GET("/dump", dumpHandler(world))
	r.GET("/dump/world", dumpWorldHandler(world))
}

// @Summary Create an entity
// @Description Creates a new entity in God World. Chaos may interfere.
// @Tags entities
// @Accept json
// @Produce json
// @Param entity body createRequest true "Entity info"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /create [post]
func createHandler(world *god.World) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		chaosMsg := world.Create(req.Name, req.Properties)

		response := gin.H{"status": "creation attempted"}
		if chaosMsg != "" {
			response["chaos_msg"] = chaosMsg
		}

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Get an entity
// @Description Retrieves an entity by name, may include chaos message if chaos interfered
// @Tags entities
// @Produce json
// @Param name path string true "Entity Name"
// @Success 200 {object} god.GetEntityResponse
// @Failure 404 {object} map[string]string
// @Router /get/{name} [get]
func getHandler(world *god.World) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		entity, ok, chaosMsg := world.Get(name)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "entity not found"})
			return
		}
		resp := GetEntityResponse{
			Entity:   entity,
			ChaosMsg: chaosMsg,
		}
		c.JSON(http.StatusOK, resp)
	}
}

// @Summary Destroy an entity
// @Description Removes an entity by name. Chaos may interfere.
// @Tags entities
// @Produce json
// @Param name path string true "Entity Name"
// @Success 200 {object} map[string]string
// @Router /destroy/{name} [delete]
func destroyHandler(world *god.World) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		chaosMsg := world.Destroy(name)

		response := gin.H{"status": "destroy attempted"}
		if chaosMsg != "" {
			response["chaos_msg"] = chaosMsg
		}

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Display world Entities
// @Description Shows all entities currently in the world
// @Tags world
// @Produce json
// @Success 200 {object} map[string]god.Entity
// @Router /dump [get]
func dumpHandler(world *god.World) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, world.AllEntities())
	}
}

// @Summary Display world state
// @Description Shows everything about the current world
// @Tags world
// @Produce json
// @Success 200 {object} god.DumpResponse
// @Router /dump/world [get]
func dumpWorldHandler(world *god.World) gin.HandlerFunc {
	return func(c *gin.Context) {
		entities, events := world.Dump()
		data := DumpResponse{
			Entities: entities,
			Events:   events,
		}
		c.JSON(http.StatusOK, data)
	}
}
