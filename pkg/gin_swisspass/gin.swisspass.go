package gin_swisspass

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/schweizerischebundesbahnen/gin-swisspass/pkg/swisspass.authenticator"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ClientIdFn func(*gin.Context) string

const HTTPHEADER_AUTHORIZATION = "Authorization"
const CLIENTID_HEADERKEY = "client_id"
const CONTEXT_USERINFO_KEY = "gin-sp-userinfo"

type GinSwisspass struct {
	clientIdFn ClientIdFn
	verifier   swisspass_authenticator.TokenVerifier
}

func New(endpoint string, clientIdFn ClientIdFn, timeoutInSeconds int) *GinSwisspass {
	return &GinSwisspass{clientIdFn, swisspass_authenticator.New(endpoint, timeoutInSeconds)}
}

func (sp *GinSwisspass) Auth(allowedRoles []string, allowedSwisspassIds []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientId := sp.clientIdFn(ctx)
		bearerToken := ctx.Request.Header.Get(HTTPHEADER_AUTHORIZATION)
		userInfo, err := sp.verifier.VerifyToken(clientId, bearerToken)

		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.Wrap(err, "could not get userInfo from GinSwisspass-token"))
			log.WithField("path", ctx.Request.URL.Path).Warn("[Swisspass-Auth] access not allowed")
			return
		}

		if authorized := sp.validateRolesAndUsers(userInfo, allowedRoles, allowedSwisspassIds); !authorized {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("authorization failed"))
			log.WithField("path", ctx.Request.URL.Path).Warn("[Swisspass-Auth] authorization failed")
			return
		}

		ctx.Set(CONTEXT_USERINFO_KEY, userInfo)
		log.WithField("path", ctx.Request.URL.Path).
			WithField("user", userInfo.SwissPassId).
			WithField("name", userInfo.FirstName+" "+userInfo.LastName).
			Trace("[Swisspass-Auth] access granted")
	}
}

func (sp *GinSwisspass) validateRolesAndUsers(userinfo swisspass_authenticator.SwisspassUserInfo, allowedRoles []string, allowedSwisspassIds []string) bool {

	for _, allowedId := range allowedSwisspassIds {
		if allowedId == userinfo.SwissPassId {
			log.Tracef("user %s authorized by allowd principal id %s", userinfo.SwissPassId, allowedId)
			return true
		}
	}

	for _, allowedRole := range allowedRoles {
		if sp.contains(userinfo.Roles, allowedRole) {
			log.Tracef("user %s is authorized by role %s", userinfo.Roles, allowedRole)
			return true
		}
	}
	return len(allowedRoles) == 0 && len(allowedSwisspassIds) == 0 //just authenticated users
}

func (sp *GinSwisspass) isEmptyAndThereforAllowed(array []string) bool {
	return len(array) == 0
}

func (sp *GinSwisspass) contains(array []string, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func ClientIdFromHeader() func(*gin.Context) string {
	return ClientIdFromHeaderWithKey(CLIENTID_HEADERKEY)
}

func ClientIdFromString(clientId string) func(*gin.Context) string {
	return func(*gin.Context) string {
		return clientId
	}
}

func ClientIdFromHeaderWithKey(key string) func(*gin.Context) string {
	return func(ctx *gin.Context) string {
		return ctx.Request.Header.Get(key)
	}
}
