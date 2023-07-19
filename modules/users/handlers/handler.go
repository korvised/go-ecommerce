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
	signUpParserErr userHandlersErrCode = "users-001"
)

type IUserHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
}

type userHandler struct {
	cfg         config.IConfig
	userUsecase usersUsecases.IUserUsecase
}

func UserHandler(cfg config.IConfig, userUsecase usersUsecases.IUserUsecase) IUserHandler {
	return &userHandler{
		cfg:         cfg,
		userUsecase: userUsecase,
	}
}

func (h userHandler) SignUpCustomer(c *fiber.Ctx) error {
	// Request body parser
	fmt.Println("call signup")
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		fmt.Println("error")
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(signUpParserErr), err.Error()).Res()
	}

	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(signUpParserErr), "invalid email pattern").Res()
	}

	// Insert
	result, err := h.userUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username have been used", "email have been used":
			return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(signUpParserErr), err.Error()).Res()
		default:
			return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(signUpParserErr), err.Error()).Res()
		}
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}
