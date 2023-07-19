package servers

import (
	"github.com/gofiber/fiber/v2"
	middlewareHandlers "github.com/korvised/go-ecommerce/modules/middlewares/handlers"
	middlewareRepositories "github.com/korvised/go-ecommerce/modules/middlewares/repositories"
	middlewareUsecases "github.com/korvised/go-ecommerce/modules/middlewares/usecase"
	monitorHandlers "github.com/korvised/go-ecommerce/modules/monitor/handlers"
)

type IModuleFactory interface {
	MonitorModule()
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
