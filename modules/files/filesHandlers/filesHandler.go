package filesHandlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/entities"
	"github.com/korvised/go-ecommerce/modules/files"
	"github.com/korvised/go-ecommerce/modules/files/filesUsecases"
	"github.com/korvised/go-ecommerce/pkg/utils"
	"math"
	"path/filepath"
	"strings"
)

type filesHandlersErrCode string

const (
	uploadFileErr filesHandlersErrCode = "files-001"
	deleteFileErr filesHandlersErrCode = "files-002"
)

type IFilesHandler interface {
	UploadFile(c *fiber.Ctx) error
	DeleteFile(c *fiber.Ctx) error
}

type filesHandler struct {
	cfg          config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func FilesHandler(cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IFilesHandler {
	return &filesHandler{
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}

func (h filesHandler) UploadFile(c *fiber.Ctx) error {
	req := make([]*files.FileReq, 0)

	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(uploadFileErr), err.Error()).Res()
	}

	filesReq := form.File["files"]
	destination := c.FormValue("destination")

	// Files ext validation
	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
	}

	for _, file := range filesReq {
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Filename), "."))
		if extMap[ext] != ext || extMap[ext] == "" {
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(uploadFileErr),
				"files are not acceptable",
			).Res()
		}

		if file.Size > int64(h.cfg.App().FileLimit()) {
			maxMiB := int(math.Ceil(float64(h.cfg.App().FileLimit()) / math.Pow(1024, 2)))

			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(uploadFileErr),
				fmt.Sprintf("file size must less than than %d MiB", maxMiB),
			).Res()
		}

		filename := utils.RandFileName(ext)

		req = append(req, &files.FileReq{
			File:        file,
			FileName:    filename,
			Destination: destination + "/" + filename,
			Extension:   ext,
		})
	}

	res, err := h.filesUsecase.UploadToStorage(req)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(uploadFileErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, res).Res()
}

func (h filesHandler) DeleteFile(c *fiber.Ctx) error {
	req := make([]*files.DeleteFileReq, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(deleteFileErr), err.Error()).Res()
	}

	if err := h.filesUsecase.DeleteFileOnStorage(req); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(deleteFileErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}
