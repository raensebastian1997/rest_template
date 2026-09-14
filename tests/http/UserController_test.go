package http_test

import (
	"net/http/httptest"
	"testing"

	userhttp "restapirian/internal/delivery/http"

	"github.com/gin-gonic/gin"
)

func TestUserControllerIndexDatabaseUnavailable(t *testing.T) {
	gin.SetMode(gin.TestMode)

	request := httptest.NewRequest("GET", "/api/v1/", nil)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request

	controller := userhttp.NewUserController(nil)
	controller.Index(ctx)

	if recorder.Code != 500 {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}

	expected := `{"error":"database unavailable"}`
	if recorder.Body.String() != expected {
		t.Fatalf("expected body %s, got %s", expected, recorder.Body.String())
	}
}
