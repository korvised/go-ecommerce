package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/korvised/go-ecommerce/modules/appinfo/appinfoHandlers"
	"github.com/korvised/go-ecommerce/modules/appinfo/appinfoRepositories"
	"github.com/korvised/go-ecommerce/modules/appinfo/appinfoUsecases"
	"github.com/korvised/go-ecommerce/modules/files/filesHandlers"
	"github.com/korvised/go-ecommerce/modules/files/filesUsecases"
	"github.com/korvised/go-ecommerce/modules/middlewares"
	"github.com/korvised/go-ecommerce/modules/middlewares/middlewaresHandlers"
	"github.com/korvised/go-ecommerce/modules/middlewares/middlewaresRepositories"
	"github.com/korvised/go-ecommerce/modules/middlewares/middlewaresUsecases"
	"github.com/korvised/go-ecommerce/modules/monitor/MonitorHandlers"
	"github.com/korvised/go-ecommerce/modules/orders/ordersHandlers"
	"github.com/korvised/go-ecommerce/modules/orders/ordersRepositories"
	"github.com/korvised/go-ecommerce/modules/orders/ordersUsecases"
	"github.com/korvised/go-ecommerce/modules/products/productsHandlers"
	"github.com/korvised/go-ecommerce/modules/products/productsRepositories"
	"github.com/korvised/go-ecommerce/modules/products/productsUsecases"
	"github.com/korvised/go-ecommerce/modules/users/userHandlers"
	"github.com/korvised/go-ecommerce/modules/users/userRepositories"
	"github.com/korvised/go-ecommerce/modules/users/userUsecases"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
	FilesModule()
	ProductsModule()
	OrdersModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewareUsecase(repository)
	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.s.db)
	usecase := usersUsecases.UsersUsecase(m.s.cfg, repository)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")

	router.Post("/signup", m.mid.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.mid.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.mid.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.mid.ApiKeyAuth(), handler.SingOut)
	router.Post("/signup-admin", m.mid.ApiKeyAuth(), handler.SignUpAdmin)

	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(middlewares.RoleAdmin), handler.GenerateAdminToken)
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.s.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.s.cfg, usecase)

	router := m.r.Group("/appinfo")

	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorize(middlewares.RoleAdmin), handler.GenerateApiKey)
	router.Get("/categories", m.mid.ApiKeyAuth(), handler.FindCategories)
	router.Post("/categories", m.mid.JwtAuth(), m.mid.Authorize(middlewares.RoleAdmin), handler.AddCategories)
	router.Delete("/categories/:category_id", m.mid.JwtAuth(), m.mid.Authorize(middlewares.RoleAdmin), handler.DeleteCategory)
}

func (m *moduleFactory) FilesModule() {
	usecase := filesUsecases.FilesUsecase(m.s.cfg)
	handler := filesHandlers.FilesHandler(m.s.cfg, usecase)

	router := m.r.Group("/files")

	router.Post("/upload", m.mid.JwtAuth(), m.mid.Authorize(middlewares.RoleAdmin), handler.UploadFile)
	router.Patch("/delete", m.mid.JwtAuth(), m.mid.Authorize(middlewares.RoleAdmin), handler.DeleteFile)
}

func (m *moduleFactory) ProductsModule() {
	fileUsecase := filesUsecases.FilesUsecase(m.s.cfg)

	repository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, fileUsecase)
	usecase := productsUsecases.ProductsUsecase(repository)
	handler := productsHandlers.ProductsHandler(m.s.cfg, usecase, fileUsecase)

	router := m.r.Group("/products")

	router.Post("/", m.mid.JwtAuth(), m.mid.Authorize(middlewares.RoleAdmin), handler.AddProduct)

	router.Patch("/:product_id", m.mid.JwtAuth(), m.mid.Authorize(middlewares.RoleAdmin), handler.UpdateProduct)

	router.Get("/", m.mid.ApiKeyAuth(), handler.FindManyProducts)
	router.Get("/:product_id", m.mid.ApiKeyAuth(), handler.FindOneProduct)

	router.Delete("/:product_id", m.mid.JwtAuth(), m.mid.Authorize(middlewares.RoleAdmin), handler.DeleteProduct)
}

func (m *moduleFactory) OrdersModule() {
	fileUsecase := filesUsecases.FilesUsecase(m.s.cfg)
	productsRepository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, fileUsecase)

	repository := ordersRepositories.OrdersRepository(m.s.db)
	usecase := ordersUsecases.OrdersUsecase(repository, productsRepository)
	handler := ordersHandlers.OrdersHandler(m.s.cfg, usecase)

	router := m.r.Group("/orders")

	router.Get("/", m.mid.JwtAuth(), m.mid.Authorize(middlewares.RoleAdmin), handler.FindManyOrders)
	router.Get("/:order_id", m.mid.JwtAuth(), handler.FindOneOrder)
}
