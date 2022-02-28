package aliyundrive

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthService_LoginByRefreshToken(t *testing.T) {
	refreshToken, empty := os.LookupEnv("REFRESH_TOKEN")
	if !empty {
		return
	}

	cli := New()

	// 获取 refresh_token 可以通过网页端： JSON.parse(localStorage.token).refresh_token
	userInfo, err := cli.Auth.LoginByRefreshToken(context.TODO(), refreshToken)

	assert.NoError(t, err)
	assert.NotNil(t, userInfo)
	assert.NotEmpty(t, userInfo.UserID)
}
