package servers

import (
	"github.com/korvised/go-ecommerce/modules/middlewares"
	"github.com/korvised/go-ecommerce/modules/products/productsHandlers"
	"github.com/korvised/go-ecommerce/modules/products/productsRepositories"
	"github.com/korvised/go-ecommerce/modules/products/productsUsecases"
)

type IProductModule interface {
	Init()
	Repository() productsRepositories.IProductsRepository
	Usecase() productsUsecases.IProductsUsecase
	Handler() productsHandlers.IProductsHandler
}

type productModule struct {
	*moduleFactory
	repository productsRepositories.IProductsRepository
	usecase    productsUsecases.IProductsUsecase
	handler    productsHandlers.IProductsHandler
}

func (m *moduleFactory) ProductsModule() IProductModule {
	repository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, m.FilesModule().Usecase())
	usecase := productsUsecases.ProductsUsecase(repository)
	handler := productsHandlers.ProductsHandler(m.s.cfg, usecase, m.FilesModule().Usecase())

	return &productModule{
		moduleFactory: m,
		usecase:       usecase,
		handler:       handler,
	}
}

func (p *productModule) Init() {

	router := p.r.Group("/products")

	router.Post("/", p.mid.JwtAuth(), p.mid.Authorize(middlewares.RoleAdmin), p.handler.AddProduct)

	router.Patch("/:product_id", p.mid.JwtAuth(), p.mid.Authorize(middlewares.RoleAdmin), p.handler.UpdateProduct)

	router.Get("/", p.mid.ApiKeyAuth(), p.handler.FindManyProducts)
	router.Get("/:product_id", p.mid.ApiKeyAuth(), p.handler.FindOneProduct)

	router.Delete("/:product_id", p.mid.JwtAuth(), p.mid.Authorize(middlewares.RoleAdmin), p.handler.DeleteProduct)
}

func (p *productModule) Repository() productsRepositories.IProductsRepository { return p.repository }

func (p *productModule) Usecase() productsUsecases.IProductsUsecase { return p.usecase }

func (p *productModule) Handler() productsHandlers.IProductsHandler { return p.handler }
