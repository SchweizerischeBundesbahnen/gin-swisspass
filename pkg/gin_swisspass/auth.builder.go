package gin_swisspass

import "github.com/gin-gonic/gin"

type AuthBuilder struct {
	swisspass           *GinSwisspass
	allowedRoles        []string
	allowedSwisspassIds []string
}

func NewBuilder(swisspass *GinSwisspass) *AuthBuilder {
	return &AuthBuilder{swisspass: swisspass}
}

func (builder *AuthBuilder) AllowRole(role ...string) *AuthBuilder {
	builder.allowedRoles = role
	return builder
}

func (builder *AuthBuilder) AllowSwisspassId(id ...string) *AuthBuilder {
	builder.allowedSwisspassIds = id
	return builder
}

func (builder *AuthBuilder) Build() gin.HandlerFunc {
	return builder.swisspass.Auth(builder.allowedRoles, builder.allowedSwisspassIds)
}
