package aliyundrive

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthService_RefreshToken(t *testing.T) {
	refreshToken, empty := os.LookupEnv("REFRESH_TOKEN")
	if !empty {
		return
	}

	cli := New()

	// 获取 refresh_token 可以通过网页端： JSON.parse(localStorage.token).refresh_token
	token, err := cli.Auth.RefreshToken(context.TODO(), &RefreshTokenReq{
		RefreshToken: refreshToken,
	})

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.RefreshToken)
}
