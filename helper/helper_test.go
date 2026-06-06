package helper

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIResponse(t *testing.T) {
	data := map[string]string{"key": "value"}

	response := APIResponse("Created", http.StatusCreated, "success", data)

	assert.Equal(t, "Created", response.Meta.Message)
	assert.Equal(t, http.StatusCreated, response.Meta.Code)
	assert.Equal(t, "success", response.Meta.Status)
	assert.Equal(t, data, response.Data)
}

func TestAPIResponseAllowsNilData(t *testing.T) {
	response := APIResponse("No content", http.StatusNoContent, "success", nil)

	assert.Nil(t, response.Data)
	assert.Equal(t, http.StatusNoContent, response.Meta.Code)
}
