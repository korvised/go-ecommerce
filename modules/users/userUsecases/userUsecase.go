package usersUsecases

import (
	"fmt"
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/users"
	"github.com/korvised/go-ecommerce/modules/users/userRepositories"
	"github.com/korvised/go-ecommerce/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecase interface {
	InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error)
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
	RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
	DeleteOauth(oauthID string) error
	GetUserProfile(userID string) (*users.User, error)
}

type usersUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUsersRepository
}

func UsersUsecase(cfg config.IConfig, userRepository usersRepositories.IUsersRepository) IUsersUsecase {
	return &usersUsecase{
		cfg:             cfg,
		usersRepository: userRepository,
	}
}

func (u usersUsecase) InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error) {
	// Hashing password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// Insert user
	result, err := u.usersRepository.InsertUser(req, true)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	// Hashing password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// Insert user
	result, err := u.usersRepository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u usersUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	// Find user
	user, err := u.usersRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// Compare password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Sign tokens
	accessToken, err := auth.NewAuth(auth.Access, u.cfg.Jwt(), &users.UserClaims{
		ID:     user.ID,
		RoleID: user.RoleID,
	})

	refreshToken, err := auth.NewAuth(auth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		ID:     user.ID,
		RoleID: user.RoleID,
	})

	// Set password
	passport := &users.UserPassport{
		User: &users.User{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			RoleID:   user.RoleID,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	// Insert oauth session
	if err = u.usersRepository.InsertOauth(passport); err != nil {
		return nil, err
	}

	return passport, nil
}

func (u usersUsecase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
	// Parse token
	claims, err := auth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Check oauth
	oauth, err := u.usersRepository.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Find profile
	profile, err := u.usersRepository.GetProfile(oauth.UserID)
	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		ID:     profile.ID,
		RoleID: profile.RoleID,
	}

	accessToken, err := auth.NewAuth(auth.Access, u.cfg.Jwt(), newClaims)
	if err != nil {
		return nil, err
	}

	refreshToken := auth.RepeatToken(u.cfg.Jwt(), newClaims, claims.ExpiresAt.Unix())

	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserToken{
			ID:           oauth.ID,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}

	if err = u.usersRepository.UpdateOauth(passport.Token); err != nil {
		if err != nil {
			return nil, err
		}
	}

	return passport, nil
}

func (u *usersUsecase) DeleteOauth(oauthID string) error {
	if err := u.usersRepository.DeleteOauth(oauthID); err != nil {
		return err
	}

	return nil
}

func (u *usersUsecase) GetUserProfile(userID string) (*users.User, error) {
	profile, err := u.usersRepository.GetProfile(userID)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
