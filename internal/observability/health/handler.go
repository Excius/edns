package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CheckResult struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type Response struct {
	Status string        `json:"status"`
	Checks []CheckResult `json:"checks"`
}

type Handler struct {
	checkers []Checker
}

func NewHandler(checkers ...Checker) *Handler {
	return &Handler{
		checkers: checkers,
	}
}

func (h *Handler) Ready(c *gin.Context) {
	response := Response{
		Status: "ready",
	}

	for _, checker := range h.checkers {

		result := checker.Check(c.Request.Context())

		if result.Status != "up" {
			response.Status = "not_ready"
		}

		response.Checks = append(response.Checks, result)
	}

	statusCode := http.StatusOK
	if response.Status != "ready" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}
