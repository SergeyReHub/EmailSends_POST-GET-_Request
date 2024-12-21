package content

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_post_contents_success(t *testing.T) {
	response_recorder := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(response_recorder)
	engine.POST("/v1/api/emails", Post_contents)
	request_body := `{"to":"serzh.rybakov.06@gmail.com","carbon_copy_recipients":["serzh.rybakov.06@mail.ru","jopa342@mail.ru"],"subject":"Test Subject","body":"This is a test email."}`
	ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/api/emails", bytes.NewBuffer([]byte(request_body)))
	engine.ServeHTTP(response_recorder, ctx.Request)

	assert.Equal(t, http.StatusOK, response_recorder.Code)
	assert.Equal(t, request_body, response_recorder.Body.String())
}

func Test_post_contents_failure(t *testing.T) {
	response_recorder := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(response_recorder)
	engine.POST("/v1/api/emails", Post_contents)

	ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/api/emails", nil)
	engine.ServeHTTP(response_recorder, ctx.Request)

	assert.Equal(t, http.StatusBadRequest, response_recorder.Code)
}

// func Test_get_contents_success(t *testing.T) {
// 	response_recorder := httptest.NewRecorder()
// 	ctx, engine := gin.CreateTestContext(response_recorder)
// 	engine.GET("/v1/api/emails", Post_contents)
// 	request_body := `{"to":"serzh.rybakov.06@gmail.com","carbon_copy_recipients":["serzh.rybakov.06@mail.ru","jopa342@mail.ru"],"subject":"Test Subject","body":"This is a test email."}`
// 	ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/api/emails", bytes.NewBuffer([]byte(request_body)))
// 	engine.ServeHTTP(response_recorder, ctx.Request)

// 	assert.Equal(t, http.StatusOK, response_recorder.Code)
// 	assert.Equal(t, request_body, response_recorder.Body.String())
// }
