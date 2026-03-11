package handlers

import (
	"net/http"
	"notification-system/api/internal/dto"
	"notification-system/api/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

func (n *NotificationHandler) CreateNotification(c *gin.Context) {

	var req dto.CreateNotificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification, err := n.notificationService.CreateNotification(c.Request.Context(), req.UserID, req.Message, req.Channels)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, notification)
}
