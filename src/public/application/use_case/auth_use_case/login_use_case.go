package auth_use_case

import (
	"errors"
	"github.com/kouhei-github/golang-ddd-boboilerplate/application/use_case/external"
	"github.com/kouhei-github/golang-ddd-boboilerplate/domain/models/user_models"
	"github.com/kouhei-github/golang-ddd-boboilerplate/domain/repositories"
)

type LoginUseCase struct {
	ur  repositories.UserRepository
	jte external.JWTTokenExternal
}

func NewLoginUseCase(
	userRepo repositories.UserRepository,
	jte external.JWTTokenExternal,
) LoginUseCase {
	return LoginUseCase{ur: userRepo, jte: jte}
}

type LoginResponse struct {
	UserId             int    `json:"user_id"`
	Token              string `json:"access_token"`
	AccessTokenExpires int    `json:"access_token_expires"`
	RefreshToken       string `json:"refresh_token"`
	AvatarURL          string `json:"avatar_url"`
}

func (lu LoginUseCase) Execute(email, password string) (*LoginResponse, error) {
	emailVo, err := user_models.NewEmail(email)
	if err != nil {
		return nil, err
	}
	passwordVo, err := user_models.NewPassword(&password)
	if err != nil {
		return nil, err
	}

	user, err := lu.ur.GetByEmail(string(emailVo))
	if err != nil {
		return nil, err
	}

	authUser, err := lu.ur.GetUserAuthByID(user.ID.(int))
	if err != nil {
		return nil, err
	}

	if !authUser.CheckPassword(string(passwordVo)) {
		return nil, errors.New("パスワードが一致しません。")
	}

	// Create token
	accessToken, err := lu.jte.GenerateToken(AccessTokenExpires, user.ID.(int), user.UserName, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := lu.jte.GenerateToken(RefreshTokenExpires, user.ID.(int), user.UserName, user.Email)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		UserId:             user.ID.(int),
		Token:              accessToken,
		AccessTokenExpires: int(AccessTokenExpires.Seconds()),
		RefreshToken:       refreshToken,
		AvatarURL:          user.Image,
	}, nil

}