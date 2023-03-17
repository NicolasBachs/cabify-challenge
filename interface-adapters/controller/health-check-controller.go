package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheckController is the controller responsible for handling health check requests.
type HealthCheckController struct {
}

func NewHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}

// GetStatus
//
//	@Summary		Status of the service
//	@Description	Indicate the service has started up correctly and is ready to accept requests.
//	@Tags			Healtcheck
//	@Produce		json
//	@Success		200
//	@Router			/status [get]
func (c *HealthCheckController) GetStatus(ctx *gin.Context) {
	ctx.String(http.StatusOK, `{"status":"ok"}`)
}
