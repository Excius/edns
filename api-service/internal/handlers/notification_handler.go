package handlers

import (
	"net/http"
	"strings"

	"github.com/excius/edns/api-service/internal/dto"
	"github.com/excius/edns/api-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

func (n *NotificationHandler) GetNotificationByID(c *gin.Context) {
	id := c.Param("id")

	notificationID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification_id"})
		return
	}

	notification, err := n.notificationService.GetNotificationByID(c.Request.Context(), notificationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notification)
}

func (n *NotificationHandler) GetDeliveriesByNotificationID(c *gin.Context) {
	id := c.Param("id")

	notificationID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	delivery, err := n.notificationService.GetDeliveriesByNotificatoinID(c.Request.Context(), notificationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, delivery)
}

func (n *NotificationHandler) CreateNotification(c *gin.Context) {

	var req dto.CreateNotificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title cannot be empty"})
	}

	message := strings.TrimSpace(req.Message)
	if message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message cannot be empty"})
		return
	}

	if len(req.Channels) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one channel must be specified"})
		return
	}

	notification, err := n.notificationService.CreateNotification(c.Request.Context(), userID, title, message, req.Channels)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, notification)
}
