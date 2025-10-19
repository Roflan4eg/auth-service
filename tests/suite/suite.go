package suite

import (
	"context"
	"github.com/Roflan4eg/auth-serivce/config"
	auth "github.com/Roflan4eg/auth-serivce/internal/interfaces/grpc/pb"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/caarlos0/env/v11"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient auth.AuthServiceClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()
	cfg, err := LoadFromFile("../../config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.WriteTimeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.NewClient(cfg.GRPC.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	return ctx, &Suite{Cfg: cfg, AuthClient: auth.NewAuthServiceClient(cc)}

}

func RandomPass() string {
	pass := gofakeit.Password(true, true, true, true, false, rand.Intn(10)+8)
	pass += "Test0!"
	return pass
}

func JWTParse(token, secret string) (jwt.MapClaims, error) {
	pToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims := pToken.Claims.(jwt.MapClaims)
	return claims, nil
}

func LoadFromFile(path string) (*config.Config, error) {
	var cfg config.Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, err
	}
	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
