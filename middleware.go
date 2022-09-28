package middleware

import (
	"github.com/alifcapital/auth-middleware/config"
	//"github.com/dgrijalva/jwt-go"
	"net/http"

	//"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gopkg.in/square/go-jose.v2/jwt"

	"strings"
)

var Module = fx.Options(
	fx.Provide(NewMiddleware),
	config.Module,
)

type Params struct {
	fx.In
	Config config.Config
}

type Middleware interface {
	Middleware(nextHandler gin.HandlerFunc) gin.HandlerFunc
}

type middleware struct {
	config config.Config
}

func NewMiddleware(p Params) (Middleware, error) {

	middlewareItem := &middleware{
		config: p.Config,
	}

	return middlewareItem, nil
}

type TokenPayload struct {
	Exp               int      `json:"exp"`
	Iat               int      `json:"iat"`
	AuthTime          int      `json:"auth_time"`
	Jti               string   `json:"jti"`
	Iss               string   `json:"iss"`
	Sub               string   `json:"sub"`
	Typ               string   `json:"typ"`
	Azp               string   `json:"azp"`
	Nonce             string   `json:"nonce"`
	SessionState      string   `json:"session_state"`
	Acr               string   `json:"acr"`
	AllowedOrigins    []string `json:"allowed-origins"`
	Scope             string   `json:"scope"`
	Sid               string   `json:"sid"`
	EmailVerified     bool     `json:"email_verified"`
	Systems           []string `json:"systems"`
	Name              string   `json:"name"`
	PreferredUsername string   `json:"preferred_username"`
	GivenName         string   `json:"given_name"`
	FamilyName        string   `json:"family_name"`
	Email             string   `json:"email"`
}

func (m *middleware) Middleware(nextHandler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

		tok, err := jwt.ParseSigned(raw)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization failed: " + err.Error()})
			return
		}

		var out TokenPayload

		if err := tok.UnsafeClaimsWithoutVerification(&out); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token " + err.Error()})
			return
		}

		authenticated := false

		for _, system := range out.Systems {
			if system == m.config.Get().AuthMiddlewareServiceName {
				authenticated = true
			}
		}

		if !authenticated {
			c.JSON(http.StatusForbidden, gin.H{"error": "no access to the system"})
			return
		}

		c.Set("id", out.Sub)
		c.Set("username", out.GivenName)

		nextHandler(c)
	}
}
