package usersHandlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/entities"
	"github.com/korvised/go-ecommerce/modules/users"
	usersUsecases "github.com/korvised/go-ecommerce/modules/users/usecases"
)

type userHandlersErrCode string

const (
	signUpCustomerErr userHandlersErrCode = "users-001"
	signInErr         userHandlersErrCode = "users-002"
)

type IUsersHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
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
	fmt.Println("start inserting")
	result, err := h.usersUsecase.InsertCustomer(req)
	if err != nil {
		fmt.Println(err.Error())
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
		return entities.NewResponse(c).Error(fiber.StatusUnauthorized, string(signInErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}
