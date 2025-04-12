package auth

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/net/ghttp"
	"golang.org/x/crypto/bcrypt"
)

type login struct {
	repo *repository.Repository
}

func Login(repo *repository.Repository) *login {
	return &login{
		repo: repo,
	}
}

func (uc *login) Usecase(ctx context.Context, req *dto.LoginRequest) (*ghttp.ResponseBody, error) {
	mUser := &model.User{}

	mUser, err := uc.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if err1 := bcrypt.CompareHashAndPassword([]byte(mUser.Password), []byte(req.Password)); err1 != nil {
		return nil, err1
	}

	token, err := generateToken(
		withUsername(mUser.Username),
		withRole(mUser.Role),
	)
	if err != nil {
		return nil, err
	}
	resp := dto.LoginResponse{
		Token: token,
	}

	return ghttp.ResponseBodyOK(resp), nil
}
