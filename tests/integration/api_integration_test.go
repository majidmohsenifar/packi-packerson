package tests

import (
	"encoding/json"
	"job-test/api"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	api.InitApi(r)
	return r
}

func TestIntegration_CorrectEndpoint(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/correct?x=750", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, "status code mismatch")

	var resp struct {
		Status bool `json:"status"`
		Data   struct {
			Ordered int         `json:"ordered"`
			Packs   map[int]int `json:"packs"`
		} `json:"data"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.True(t, resp.Status)
	require.Equal(t, 750, resp.Data.Ordered)
	require.Equal(t, map[int]int{500: 1, 250: 1}, resp.Data.Packs)
}

func TestIntegration_InCorrectEndpoint(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/incorrect?x=251", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Status bool `json:"status"`
		Data   struct {
			Ordered int           `json:"ordered"`
			Packs   []map[int]int `json:"packs"`
		} `json:"data"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.True(t, resp.Status)
	require.Equal(t, 251, resp.Data.Ordered)
	require.Contains(t, resp.Data.Packs, map[int]int{250: 2})
}
