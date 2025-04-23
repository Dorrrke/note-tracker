package server

import (
	"context"
	"fmt"
	"time"

	"github.com/Dorrrke/note-tracker/gen/auth"
	"github.com/Dorrrke/note-tracker/internal/auth-service/domain/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var jwtKey = []byte("super_secret_key_1285@ggwg")

type Repository interface {
	LoginUser(userCreds models.UserCreds) (models.User, error)
	RegisterUser(user models.User) (string, error)
}

type GRPCServer struct {
	auth.UnimplementedAuthServiceServer
	repo Repository
}

func RegisterGrpcServer(gRPC *grpc.Server, repo Repository) {
	auth.RegisterAuthServiceServer(gRPC, &GRPCServer{repo: repo})
}

func (g *GRPCServer) Login(
	ctx context.Context,
	req *auth.UserCredentials,
) (*auth.AuthResponse, error) {
	user := models.UserCreds{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	}

	dbUser, err := g.repo.LoginUser(user)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	token, err := genJwtToken(dbUser.UID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.AuthResponse{
		Token: token,
	}, nil
}

func (g *GRPCServer) Register(
	ctx context.Context,
	req *auth.User,
) (*auth.AuthResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	user := models.User{
		Name:     req.GetName(),
		Login:    req.GetLogin(),
		Password: string(hash),
	}

	uid, err := g.repo.RegisterUser(user)
	if err != nil {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}

	token, err := genJwtToken(uid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &auth.AuthResponse{Token: token}, nil
}

func genJwtToken(uid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   uid,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
	})
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return ``, err
	}
	return tokenStr, nil
}

func validateJwtToken(tokenStr string) (string, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return ``, err
	}

	if !token.Valid {
		return ``, fmt.Errorf("invalid token")
	}

	return claims.Subject, nil
}
