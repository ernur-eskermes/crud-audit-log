package restHandler

import (
	"context"

	"github.com/ernur-eskermes/crud-audit-log/pkg/logging"
	"github.com/gofiber/fiber/v2"
)

type Audit struct {
	service AuditService
	logger  *logging.Logger
}

func NewAudit(service AuditService, logger *logging.Logger) *Audit {
	return &Audit{
		service: service,
		logger:  logger,
	}
}

func (h *Audit) initAuditRoutes(api fiber.Router) {
	auditLog := api.Group("/logs")
	{
		auditLog.Get("", h.getAllLogs)
	}
}

// @Summary Get Logs
// @Tags logs
// @Description Get all logs
// @ModuleID getAllLogs
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.LogItem
// @Router /logs [get]
func (h *Audit) getAllLogs(c *fiber.Ctx) error {
	logs, err := h.service.GetAll(context.TODO())
	if err != nil {
		h.logger.Error(err)

		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.JSON(logs)
}
