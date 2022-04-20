package restHandler

import (
	"context"
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2"
	audit "github.com/ernur-eskermes/crud-audit-log/pkg/domain"
	"github.com/ernur-eskermes/crud-audit-log/pkg/logging"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

type AuditService interface {
	GetAll(ctx context.Context) ([]audit.LogItem, error)
}

type Deps struct {
	AuditService AuditService
	Logger       *logging.Logger
}

type Handler struct {
	auditHandler *Audit
}

func New(deps Deps) *Handler {
	return &Handler{
		auditHandler: NewAudit(deps.AuditService, deps.Logger),
	}
}

func (h *Handler) InitRouter(app *fiber.App) {
	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		TimeFormat: time.RFC3339,
		TimeZone:   "Asia/Almaty",
	}))
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        20,
		Expiration: 30 * time.Second,
	}))

	app.Get("/dashboard", monitor.New())
	app.Get("/swagger/*", swagger.HandlerDefault)
	h.initAPI(app)
}

func (h *Handler) initAPI(app fiber.Router) {
	api := app.Group("/api")
	{
		h.auditHandler.initAuditRoutes(api)
	}
}

type response struct {
	Message string `json:"message"`
}
