package tests

import (
	auth "github.com/Roflan4eg/auth-serivce/internal/interfaces/grpc/pb"
	"github.com/Roflan4eg/auth-serivce/tests/suite"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRegisterLogin_Login_Success(t *testing.T) {
	ctx, st := suite.New(t)

	secret := st.Cfg.JWTConfig.Secret
	deltaSec := 1

	email := gofakeit.Email()
	pass := suite.RandomPass()

	respReg, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Email:           email,
		Password:        pass,
		PasswordConfirm: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetAccessToken())
	assert.NotEmpty(t, respReg.GetRefreshToken())

	regRefClaims, err := suite.JWTParse(respReg.GetRefreshToken(), secret)
	require.NoError(t, err)

	regAccClaims, err := suite.JWTParse(respReg.GetAccessToken(), secret)
	require.NoError(t, err)

	respLog, err := st.AuthClient.Login(ctx, &auth.LoginRequest{
		Email:    email,
		Password: pass,
	})
	loginTime := time.Now()
	require.NoError(t, err)

	refreshToken := respLog.GetRefreshToken()
	require.NotEmpty(t, refreshToken)

	accessToken := respLog.GetAccessToken()
	require.NotEmpty(t, accessToken)

	logRefClaims, err := suite.JWTParse(refreshToken, secret)
	require.NoError(t, err)

	assert.Equal(t, regRefClaims["uid"].(string), logRefClaims["uid"].(string))

	logAccClaims, err := suite.JWTParse(accessToken, secret)
	require.NoError(t, err)
	assert.Equal(t, regAccClaims["uid"].(string), logAccClaims["uid"].(string))
	// check expiration
	assert.InDelta(t, loginTime.Add(st.Cfg.JWTConfig.RefreshTokenTTL).Unix(), logRefClaims["exp"].(float64), float64(deltaSec))
	assert.InDelta(t, loginTime.Add(st.Cfg.JWTConfig.AccessTokenTTL).Unix(), logAccClaims["exp"].(float64), float64(deltaSec))

	// check issued
	assert.InDelta(t, loginTime.Unix(), logRefClaims["iat"].(float64), float64(deltaSec))
	assert.InDelta(t, loginTime.Unix(), logAccClaims["iat"].(float64), float64(deltaSec))

}

func TestRegisterLogin_UserExists(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := suite.RandomPass()

	resp, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Email:           email,
		Password:        pass,
		PasswordConfirm: pass,
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.GetAccessToken())
	require.NotEmpty(t, resp.GetRefreshToken())

	resp, err = st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Email:           email,
		Password:        pass,
		PasswordConfirm: pass,
	})
	require.Error(t, err)
	require.Empty(t, resp.GetAccessToken())
	require.Empty(t, resp.GetRefreshToken())
	assert.ErrorContains(t, err, "user already exists")

}

func TestRegisterLogin_UserWrongPassword(t *testing.T) {
	ctx, st := suite.New(t)
	email := gofakeit.Email()
	pass := suite.RandomPass()
	_, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Email:           email,
		Password:        pass,
		PasswordConfirm: pass,
	})
	require.NoError(t, err)
	_, err = st.AuthClient.Login(ctx, &auth.LoginRequest{
		Email:    email,
		Password: suite.RandomPass(),
	})
	require.Error(t, err)
	assert.ErrorContains(t, err, "invalid password")

}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)
	pass := suite.RandomPass()
	cases := []struct {
		name            string
		email           string
		password        string
		confirmPassword string
		expectedErr     string
	}{
		{
			name:            "Register with empty email",
			email:           "",
			password:        pass,
			confirmPassword: pass,
			expectedErr:     "email is required",
		},
		{
			name:            "Register with empty password",
			email:           gofakeit.Email(),
			password:        "",
			confirmPassword: pass,
			expectedErr:     "password is required",
		},
		{
			name:            "Register with empty confirmPassword",
			email:           gofakeit.Email(),
			password:        pass,
			confirmPassword: "",
			expectedErr:     "password_confirm is required",
		},
		{
			name:            "Register with different passwords",
			email:           gofakeit.Email(),
			password:        suite.RandomPass(),
			confirmPassword: pass,
			expectedErr:     "must be equal",
		},
		{
			name:            "Register with invalid email",
			email:           "invalidEmail",
			password:        pass,
			confirmPassword: pass,
			expectedErr:     "email is invalid",
		},
		{
			name:            "Register with invalid password",
			email:           gofakeit.Email(),
			password:        "invalidPassword",
			confirmPassword: pass,
			expectedErr:     "weak password",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
				Email:           c.email,
				Password:        c.password,
				PasswordConfirm: c.confirmPassword,
			})
			require.Error(t, err)
			require.ErrorContains(t, err, c.expectedErr)
		})
	}

}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)
	pass := suite.RandomPass()
	cases := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Login with empty email",
			email:       "",
			password:    pass,
			expectedErr: "email is required",
		},
		{
			name:        "Login with empty password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "Login with invalid email",
			email:       "invalidEmail",
			password:    pass,
			expectedErr: "email is invalid",
		},
		{
			name:        "Login with invalid password",
			email:       gofakeit.Email(),
			password:    "invalidPassword",
			expectedErr: "weak password",
		},
		{
			name:        "Login with non-existing user",
			email:       gofakeit.Email(),
			password:    suite.RandomPass(),
			expectedErr: "user not found",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := st.AuthClient.Login(ctx, &auth.LoginRequest{
				Email:    c.email,
				Password: c.password,
			})
			require.Error(t, err)
			require.ErrorContains(t, err, c.expectedErr)
		})
	}

}
