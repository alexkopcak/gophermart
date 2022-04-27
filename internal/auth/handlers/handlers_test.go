package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexkopcak/gophermart/internal/auth/usecase"
	"github.com/alexkopcak/gophermart/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	router := gin.Default()

	auc := new(usecase.AuthUseCaseMock)

	RegisterHTTPEndpoints(router, auc)

	user := &models.User{
		UserName: "testuser",
		Password: "testpassword",
	}

	body, err := json.Marshal(user)
	assert.NoError(t, err)

	auc.On("SignUp", user.UserName, user.Password).Return(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	//	bodyBytes, _ := io.ReadAll(w.Body)
	//	fmt.Printf("\n\n%v\n\n", string(bodyBytes))

	//assert.Equal(t, http.StatusOK, w.Code)
}
