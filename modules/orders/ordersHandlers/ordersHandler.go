package ordersHandlers

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/entities"
	"github.com/korvised/go-ecommerce/modules/middlewares"
	"github.com/korvised/go-ecommerce/modules/middlewares/middlewaresHandlers"
	"github.com/korvised/go-ecommerce/modules/orders"
	"github.com/korvised/go-ecommerce/modules/orders/ordersUsecases"
	"strings"
	"time"
)

type ordersHandlersErrCode string

const (
	fineOneOrderErr   ordersHandlersErrCode = "orders-001"
	fineManyOrdersErr ordersHandlersErrCode = "orders-002"
	insertOrderErr    ordersHandlersErrCode = "orders-003"
	updateOrderErr    ordersHandlersErrCode = "orders-004"
)

type IOrdersHandler interface {
	FindOneOrder(c *fiber.Ctx) error
	FindManyOrders(c *fiber.Ctx) error
	InsertOrder(c *fiber.Ctx) error
	UpdateOrder(c *fiber.Ctx) error
}

type ordersHandler struct {
	cfg           config.IConfig
	ordersUsecase ordersUsecases.IOrdersUsecase
}

func OrdersHandler(cfg config.IConfig, ordersUsecase ordersUsecases.IOrdersUsecase) IOrdersHandler {
	return &ordersHandler{
		cfg:           cfg,
		ordersUsecase: ordersUsecase,
	}
}

func (h *ordersHandler) FindOneOrder(c *fiber.Ctx) error {
	orderID := strings.Trim(c.Params("order_id"), " ")

	order, err := h.ordersUsecase.FindOneOrder(orderID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(fineOneOrderErr),
				"order not found",
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.StatusInternalServerError,
				string(fineOneOrderErr),
				err.Error(),
			).Res()
		}
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}

func (h *ordersHandler) FindManyOrders(c *fiber.Ctx) error {
	req := &orders.OrderFilter{
		SortReq:       &entities.SortReq{},
		PaginationReq: &entities.PaginationReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(fineManyOrdersErr), err.Error()).Res()
	}

	// Pagination
	if req.Page < 1 {
		req.Page = 1
	}

	if req.Size < 5 {
		req.Size = 5
	}

	// Sort
	req.OrderBy = strings.ToLower(req.OrderBy)
	orderByMap := map[string]string{
		"id":         `o.id`,
		"created_at": `o.created_at`,
	}
	if orderByMap[req.OrderBy] == "" {
		req.OrderBy = orderByMap["id"]
	} else {
		req.OrderBy = orderByMap[req.OrderBy]
	}

	req.Sort = strings.ToUpper(req.Sort)
	sortMap := map[string]string{
		"ASC":  "ASC",
		"DESC": "DESC",
	}
	if sortMap[req.Sort] == "" {
		req.Sort = sortMap["DESC"]
	} else {
		req.Sort = sortMap[req.Sort]
	}

	// * Start Date  format: YYYY-MM-DD
	if req.StartDate != "" {
		start, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(fineManyOrdersErr),
				"start date is invalid",
			).Res()
		}

		req.StartDate = start.Format("2006-01-02")
	}

	// * End Date  format: YYYY-MM-DD
	if req.EndDate != "" {
		end, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(fineManyOrdersErr),
				"end date is invalid",
			).Res()
		}

		req.EndDate = end.Format("2006-01-02")
	}

	// Find orders
	data := h.ordersUsecase.FindManyOrders(req)

	return entities.NewResponse(c).Success(fiber.StatusOK, data).Res()
}

func (h *ordersHandler) InsertOrder(c *fiber.Ctx) error {
	userID := c.Locals(middlewaresHandlers.UserID).(string)

	req := &orders.Order{
		Products: make([]*orders.ProductsOrder, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(insertOrderErr), err.Error()).Res()
	}

	if len(req.Products) == 0 {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(insertOrderErr),
			"products are empty",
		).Res()
	}

	req.UserID = userID
	req.Status = "waiting"
	req.TotalPaid = 0

	order, err := h.ordersUsecase.InsertOrder(req)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(insertOrderErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}

func (h *ordersHandler) UpdateOrder(c *fiber.Ctx) error {
	orderID := strings.Trim(c.Params("order_id"), " ")
	roleID := c.Locals(middlewaresHandlers.UserRoleID).(int)

	req := new(orders.UpdateOrderReq)

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.StatusBadRequest, string(updateOrderErr), err.Error()).Res()
	}

	req.ID = orderID
	req.Status = strings.ToLower(req.Status)

	statusMap := map[string]string{
		"waiting":   "waiting",
		"shipping":  "shipping",
		"completed": "completed",
		"canceled":  "canceled",
	}

	// Check role user can update only "canceled" status
	if roleID == middlewares.RoleUser && req.Status != statusMap["canceled"] {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(updateOrderErr),
			"incorrect order status",
		).Res()
	}

	if req.TransferSlip != nil {
		if req.TransferSlip.ID == "" {
			req.TransferSlip.ID = uuid.NewString()
		}

		// YYYY-MM-DD HH:MM:SS
		// 2006-01-02 15:04:05
		loc, _ := time.LoadLocation("Asia/Vientiane")
		req.TransferSlip.CreatedAt = time.Now().In(loc).Format("2006-01-02 15:04:05")
	}

	order, err := h.ordersUsecase.UpdateOrder(req)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.StatusInternalServerError, string(updateOrderErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}
