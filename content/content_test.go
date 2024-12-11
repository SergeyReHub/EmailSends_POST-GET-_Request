package content

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostContents_success(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(responseRecorder)
	engine.POST("/v1/api/emails", PostContents)
	requestBody := `{"to":"serzh.rybakov.06@gmail.com","carbon_copy_recipients":["serzh.rybakov.06@mail.ru","jopa342@mail.ru"],"subject":"Test Subject","body":"This is a test email."}`
	ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/api/emails", bytes.NewBuffer([]byte(requestBody)))
	engine.ServeHTTP(responseRecorder, ctx.Request)
	
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, requestBody, responseRecorder.Body.String())
}

func TestPostContents_failure(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(responseRecorder)
	engine.POST("/v1/api/emails", PostContents)

	ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/api/emails", nil)
	engine.ServeHTTP(responseRecorder, ctx.Request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}
