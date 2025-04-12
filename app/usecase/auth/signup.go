package auth

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/model"
	"github.com/ngoctd314/c72-api-server/pkg/repository"
	"github.com/ngoctd314/common/apperror"
	"github.com/ngoctd314/common/net/ghttp"
	"golang.org/x/crypto/bcrypt"
)

type signup struct {
	repo *repository.Repository
}

func Signup(repo *repository.Repository) *signup {
	return &signup{
		repo: repo,
	}
}

func (uc *signup) Usecase(ctx context.Context, req *dto.SignupRequest) (*ghttp.ResponseBody, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	mUser := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     model.ERoleFromString(req.Role),
	}
	if err := uc.repo.CreateUser(ctx, mUser); err != nil {
		if apperror.IsMySQLDuplicate(err) {
			return nil, apperror.ErrConflict("user is already exists")
		}
		return nil, err
	}

	return ghttp.ResponseBodyOK(nil, ghttp.ResponseBodyWithMessage("account is created successfully")), nil
}
