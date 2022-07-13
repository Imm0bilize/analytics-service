package auth

import (
	"analytic-service/internal/config"
	v1 "analytic-service/pkg/auth/proto/v1"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcAuth struct {
	conn   *grpc.ClientConn
	client v1.AuthClient
}

func New(cfg *config.Config) (*GrpcAuth, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", cfg.Auth.Host, cfg.Auth.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	client := v1.NewAuthClient(conn)
	return &GrpcAuth{conn: conn, client: client}, nil
}

func (g *GrpcAuth) RefreshTokens(ctx context.Context, accessToken, refreshToken string) (string, string, error) {
	request := &v1.RefreshTokensRequest{AccessToken: accessToken, RefreshToken: refreshToken}
	r, err := g.client.RefreshTokens(ctx, request)
	if err != nil {
		return "", "", err
	}
	return r.AccessToken, r.RefreshToken, nil
}

func (g *GrpcAuth) Shutdown() error {
	return g.conn.Close()
}
