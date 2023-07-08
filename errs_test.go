package errs_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/thirathawat/errs"
)

func TestNewError(t *testing.T) {
	err := errs.New(errs.CodeBadRequest, "Bad request")
	assert.NotNil(t, err)
	assert.Equal(t, errs.CodeBadRequest, err.Code)
	assert.Equal(t, "Bad request", err.Message)
	assert.Empty(t, err.Info)
	assert.False(t, err.Timestamp.IsZero())
}

func TestNewErrorWithInfo(t *testing.T) {
	info := map[string]interface{}{
		"field": "value",
	}
	err := errs.New(errs.CodeInternalServerError, "Internal server error",
		errs.WithInfo(info),
	)
	assert.NotNil(t, err)
	assert.Equal(t, errs.CodeInternalServerError, err.Code)
	assert.Equal(t, "Internal server error", err.Message)
	assert.Equal(t, info, err.Info)
	assert.False(t, err.Timestamp.IsZero())
}

func TestInvalidStructError(t *testing.T) {
	type test struct {
		Field string `json:"field" binding:"required"`
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/struct", func(c *gin.Context) {
		var t test
		if err := c.ShouldBindJSON(&t); err != nil {
			errs.ResponseError(c, errs.InvalidStructError(err))
			return
		}

		c.JSON(http.StatusOK, t)
	})

	w := performRequest(router, http.MethodPost, "/struct", bytes.NewBufferString(`{"field":""}`))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResponseErrorWithErrsError(t *testing.T) {
	err := errs.New(errs.CodeNotFound, "Not found")
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/errs", func(c *gin.Context) {
		errs.ResponseError(c, err)
	})

	w := performRequest(router, http.MethodGet, "/errs", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestResponseErrorWithNonErrsError(t *testing.T) {
	err := errors.New("Some error")
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/non-err", func(c *gin.Context) {
		errs.ResponseError(c, err)
	})

	w := performRequest(router, http.MethodGet, "/non-err", nil)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func performRequest(router *gin.Engine, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
