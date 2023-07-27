package servers

import (
	"github.com/korvised/go-ecommerce/modules/files/filesHandlers"
	"github.com/korvised/go-ecommerce/modules/files/filesUsecases"
	"github.com/korvised/go-ecommerce/modules/middlewares"
)

type IFileModule interface {
	Init()
	Usecase() filesUsecases.IFilesUsecase
	Handler() filesHandlers.IFilesHandler
}

type fileModule struct {
	*moduleFactory
	usecase filesUsecases.IFilesUsecase
	handler filesHandlers.IFilesHandler
}

func (m *moduleFactory) FilesModule() IFileModule {
	usecase := filesUsecases.FilesUsecase(m.s.cfg)
	handler := filesHandlers.FilesHandler(m.s.cfg, usecase)

	return &fileModule{
		moduleFactory: m,
		usecase:       usecase,
		handler:       handler,
	}
}

func (f *fileModule) Init() {
	router := f.r.Group("/files")

	router.Post("/upload", f.mid.JwtAuth(), f.mid.Authorize(middlewares.RoleAdmin), f.handler.UploadFile)
	router.Patch("/delete", f.mid.JwtAuth(), f.mid.Authorize(middlewares.RoleAdmin), f.handler.DeleteFile)
}

func (f *fileModule) Usecase() filesUsecases.IFilesUsecase { return f.usecase }

func (f *fileModule) Handler() filesHandlers.IFilesHandler { return f.handler }
