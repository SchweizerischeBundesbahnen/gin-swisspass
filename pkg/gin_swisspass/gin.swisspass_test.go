package gin_swisspass

import (
	"code.sbb.ch/ki_sjm/gin-swisspass/pkg/swisspass.authenticator"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var ginSwisspass *GinSwisspass
var dummySwissId = "1234"
var mockRouter *gin.Engine

const DUMMYPATHPUBLIC = "/public"
const DUMMYPATHPRIVATE = "/private"

type MockVerifier struct {
	SimulateFail bool
}

func (mock *MockVerifier) VerifyToken(clientId string, token string) (swisspass_authenticator.SwisspassUserInfo, error) {
	userinfo := swisspass_authenticator.SwisspassUserInfo{}
	if mock.SimulateFail {
		return userinfo, errors.New("simulated fail")
	}
	userinfo.FirstName = "first"
	userinfo.FirstName = "last"
	userinfo.Roles = []string{"ADMIN", "USER"}
	userinfo.SwissPassId = dummySwissId
	return userinfo, nil
}

func initRouter(authbuilder *AuthBuilder, simulateFail bool) {
	ginSwisspass = New(
		"",
		ClientIdFromString("clientid"),
		1)

	mockRouter = gin.New()
	ginSwisspass.verifier = &MockVerifier{SimulateFail: simulateFail}
	if authbuilder != nil {
		mockRouter.Use(authbuilder.Build())
	}
	mockRouter.GET(DUMMYPATHPUBLIC, func(ctx *gin.Context) {})
	mockRouter.GET(DUMMYPATHPRIVATE, func(ctx *gin.Context) {
		if _, ok := ctx.Get(CONTEXT_USERINFO_KEY); !ok {
			panic("userinfo is not in context")
		}
	})
}

func TestGinSwisspass_Public_Path(t *testing.T) {
	initRouter(nil, false)
	req, _ := http.NewRequest(http.MethodGet, DUMMYPATHPUBLIC, nil)

	answer := PerformRequest(mockRouter, req)

	assert.Equal(t, http.StatusOK, answer.Code)

}

func TestGinSwisspass_Authenticated(t *testing.T) {
	builder := NewBuilder(ginSwisspass)
	initRouter(builder, false)
	req, _ := http.NewRequest(http.MethodGet, DUMMYPATHPRIVATE, nil)

	answer := PerformRequest(mockRouter, req)

	assert.Equal(t, http.StatusOK, answer.Code)

}

func TestGinSwisspass_Authorized_And_Role_OK(t *testing.T) {
	builder := NewBuilder(ginSwisspass).AllowRole("ADMIN")
	initRouter(builder, false)
	req, _ := http.NewRequest(http.MethodGet, DUMMYPATHPRIVATE, nil)

	answer := PerformRequest(mockRouter, req)

	assert.Equal(t, http.StatusOK, answer.Code)

}

func TestGinSwisspass_Authorized_And_Role_NOK(t *testing.T) {
	builder := NewBuilder(ginSwisspass).AllowRole("DEVELOPER")
	initRouter(builder, false)
	req, _ := http.NewRequest(http.MethodGet, DUMMYPATHPRIVATE, nil)

	answer := PerformRequest(mockRouter, req)

	assert.Equal(t, http.StatusUnauthorized, answer.Code)

}

func TestGinSwisspass_Authorized_By_ID(t *testing.T) {
	builder := NewBuilder(ginSwisspass).
		AllowRole("DEVELOPER").
		AllowSwisspassId(dummySwissId)

	initRouter(builder, false)
	req, _ := http.NewRequest(http.MethodGet, DUMMYPATHPRIVATE, nil)

	answer := PerformRequest(mockRouter, req)

	assert.Equal(t, http.StatusOK, answer.Code)

}

func PerformRequest(router http.Handler, req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}
