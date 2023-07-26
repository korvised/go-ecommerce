package usersHandlers

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/entities"
	"github.com/korvised/go-ecommerce/modules/users"
	"github.com/korvised/go-ecommerce/modules/users/userUsecases"
	"github.com/korvised/go-ecommerce/pkg/auth"
	"strings"
)

type userHandlersErrCode string

const (
	signUpCustomerErr     userHandlersErrCode = "users-001"
	signInErr             userHandlersErrCode = "users-002"
	refreshPassportErr    userHandlersErrCode = "users-003"
	signOutErr            userHandlersErrCode = "users-004"
	signUpAdminErr        userHandlersErrCode = "users-005"
	generateAdminTokenErr userHandlersErrCode = "users-006"
	getUserProfileErr     userHandlersErrCode = "users-007"
)

type IUsersHandler interface {
	SignUpAdmin(c *fiber.Ctx) error
	SignUpCustomer(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error
	SingOut(c *fiber.Ctx) error
	GenerateAdminToken(c *fiber.Ctx) error
	GetUserProfile(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg          config.IConfig
	usersUsecase usersUsecases.IUsersUsecase
}

func UsersHandler(cfg config.IConfig, usersUsecase usersUsecases.IUsersUsecase) IUsersHandler {
	return &usersHandler{
		cfg:          cfg,
		usersUsecase: usersUsecase,
	}
}

func (h *usersHandler) GenerateAdminToken(c *fiber.Ctx) error {
	adminToken, err := auth.NewAuth(auth.Admin, h.cfg.Jwt(), nil)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(generateAdminTokenErr),
			err.Error(),
		).Res()
	}

	data := &struct {
		Token string `json:"token"`
	}{
		Token: adminToken.SignToken(),
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, data).Res()
}

func (h *usersHandler) SignUpAdmin(c *fiber.Ctx) error {
	// Request body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(signUpAdminErr), err.Error()).Res()
	}

	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(signUpAdminErr), "email pattern is invalid").Res()
	}

	// Insert
	result, err := h.usersUsecase.InsertAdmin(req)
	if err != nil {
		switch err.Error() {
		case "username is already in used", "email is already in used":
			return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(signUpAdminErr), err.Error()).Res()
		default:
			return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(signUpAdminErr), err.Error()).Res()
		}
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) SignUpCustomer(c *fiber.Ctx) error {
	// Request body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(signUpCustomerErr), err.Error()).Res()
	}

	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(signUpCustomerErr), "email pattern is invalid").Res()
	}

	// Insert
	result, err := h.usersUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username is already in used", "email is already in used":
			return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(signUpCustomerErr), err.Error()).Res()
		default:
			return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(signUpCustomerErr), err.Error()).Res()
		}
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(signInErr), err.Error()).Res()
	}

	passport, err := h.usersUsecase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(signInErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) RefreshPassport(c *fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(refreshPassportErr), err.Error()).Res()
	}

	passport, err := h.usersUsecase.RefreshPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(refreshPassportErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) SingOut(c *fiber.Ctx) error {
	req := new(users.UserRemoveCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(signOutErr), err.Error()).Res()
	}

	if err := h.usersUsecase.DeleteOauth(req.OauthID); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(signOutErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}

func (h *usersHandler) GetUserProfile(c *fiber.Ctx) error {
	// Set params
	userId := strings.Trim(c.Params("user_id"), " ")

	// Get profile
	profile, err := h.usersUsecase.GetUserProfile(userId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(getUserProfileErr), err.Error()).Res()
		default:
			return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(getUserProfileErr), err.Error()).Res()

		}
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, profile).Res()
}
