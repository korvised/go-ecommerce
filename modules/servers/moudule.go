package servers

import (
	"github.com/gofiber/fiber/v2"
	middlewareHandlers "github.com/korvised/go-ecommerce/modules/middlewares/handlers"
	middlewareRepositories "github.com/korvised/go-ecommerce/modules/middlewares/repositories"
	middlewareUsecases "github.com/korvised/go-ecommerce/modules/middlewares/usecase"
	monitorHandlers "github.com/korvised/go-ecommerce/modules/monitor/handlers"
	usersHandlers "github.com/korvised/go-ecommerce/modules/users/handlers"
	usersRepositories "github.com/korvised/go-ecommerce/modules/users/repositories"
	usersUsecases "github.com/korvised/go-ecommerce/modules/users/usecases"
)

type IModuleFactory interface {
	MonitorModule()
	UserModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewareHandlers.IMiddlewareHandler
}

func InitModule(r fiber.Router, s *server, mid middlewareHandlers.IMiddlewareHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func InitMiddlewares(s *server) middlewareHandlers.IMiddlewareHandler {
	repository := middlewareRepositories.MiddlewareRepository(s.db)
	usecase := middlewareUsecases.MiddlewareUsecase(repository)
	return middlewareHandlers.MiddlewareHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UserModule() {
	repository := usersRepositories.UserRepository(m.s.db)
	usecase := usersUsecases.UserUsecase(m.s.cfg, repository)
	handler := usersHandlers.UserHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")

	router.Post("/signup", handler.SignUpCustomer)
}
