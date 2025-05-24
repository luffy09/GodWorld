package api

import (
	"bytes"
	"encoding/json"
	"godworld/god"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() (*gin.Engine, *god.World) {
	gin.SetMode(gin.TestMode)
	world := god.NewWorld()
	r := gin.Default()
	RegisterHandlers(r, world)
	return r, world
}

func TestCreateHandler(t *testing.T) {
	r, _ := setupRouter()

	payload := CreateRequest{
		Name: "test-entity",
		Properties: map[string]string{
			"type": "test",
		},
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "creation attempted", resp["status"])
}

func TestCreateHandlerInvalidJSON(t *testing.T) {
	r, _ := setupRouter()

	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer([]byte(`{invalid json`)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetHandler(t *testing.T) {
	r, world := setupRouter()

	// Pre-create an entity
	world.Create("entity1", map[string]string{"key": "value"})

	req, _ := http.NewRequest("GET", "/get/entity1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp GetEntityResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "entity1", resp.Entity.Name)
	assert.Equal(t, "value", resp.Entity.Properties["key"])
}

func TestGetHandlerNotFound(t *testing.T) {
	r, _ := setupRouter()

	req, _ := http.NewRequest("GET", "/get/nonexistent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDestroyHandler(t *testing.T) {
	r, world := setupRouter()

	world.Create("entityToDelete", map[string]string{})

	req, _ := http.NewRequest("DELETE", "/destroy/entityToDelete", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDumpHandler(t *testing.T) {
	r, world := setupRouter()

	world.Create("e1", map[string]string{"foo": "bar"})
	world.Create("e2", map[string]string{"baz": "qux"})

	req, _ := http.NewRequest("GET", "/dump", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]god.Entity
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "e1")
	assert.Contains(t, resp, "e2")
}

func TestDumpWorldHandler(t *testing.T) {
	r, world := setupRouter()

	world.Create("entityX", map[string]string{"a": "b"})

	req, _ := http.NewRequest("GET", "/dump/world", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp DumpResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp.Entities, "entityX")
	assert.NotNil(t, resp.Events)
}
