package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/merrors"
	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/models"
	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/services"
	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	SubService *services.SubscriptionService
}

func NewSubscriptionHandler(svc *services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{SubService: svc}
}

// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Create a subscription with service name, price, user ID, start date, and optional end date
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.SubscriptionCreateReq true "Subscription data"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	slog.Info("Creating new subscription", "path", c.Request.URL.Path)

	var req models.SubscriptionCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slog.Info("Parsed subscription creation request", "user_id", req.UserID, "service_name", req.ServiceName, "price", req.Price, "start_date", req.StartDate, "end_date", req.EndDate)

	if err := h.SubService.Create(&req); err != nil {
		slog.Error("Failed to create subscription", "status", merrors.ErrorsToHTTP(err), "error", err)
		merrors.GinReturnError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Subscription created"})
}

// GetSubscription godoc
// @Summary Get subscription by ID
// @Description Retrieve a subscription by its ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} models.SubscriptionModel
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	idStr := c.Param("id")
	slog.Info("Getting subscription by ID", "id", idStr)

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	sub, err := h.SubService.GetByID(id)
	if err != nil {
		slog.Error("Failed to get subscription", "id", id, "status", merrors.ErrorsToHTTP(err), "error", err)
		merrors.GinReturnError(c, err)
		return
	}

	c.JSON(http.StatusOK, sub)
}

// UpdateSubscription godoc
// @Summary Update subscription
// @Description Update an existing subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param subscription body models.SubscriptionUpdateReq true "Updated subscription data"
// @Success 200 {object} map[string]string "Updated"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/subscriptions/{id} [patch]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var sub models.SubscriptionUpdateReq
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Parsed subscription update request", "id", id, "user_id", sub.UserID, "service_name", sub.ServiceName, "price", sub.Price, "start_date", sub.StartDate, "end_date", sub.EndDate)

	if err := h.SubService.Update(id, &sub); err != nil {
		slog.Error("Failed to update subscription", "id", id, "status", merrors.ErrorsToHTTP(err), "error", err)
		merrors.GinReturnError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription updated"})
}

// DeleteSubscription godoc
// @Summary Delete subscription
// @Description Delete a subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} map[string]string "Deleted"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	idStr := c.Param("id")
	slog.Info("Deleting subscription", "id", idStr)

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.SubService.Delete(id); err != nil {
		slog.Error("Failed to delete subscription", "id", id, "status", merrors.ErrorsToHTTP(err), "error", err)
		merrors.GinReturnError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription deleted"})
}

// ListSubscriptions godoc
// @Summary List subscriptions by user ID
// @Description Retrieve all subscriptions for a user
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID (UUID)"
// @Param name query string false "Service name"
// @Param from query string false "Start date (MM-YYYY)"
// @Param till query string false "End date (MM-YYYY)"
// @Success 200 {array} models.SubscriptionModel
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
	slog.Info("Listing subscriptions", "query", c.Request.URL.RawQuery)

	filter, err := models.NewSubscriptionFilterFromURL(c.Request.URL.Query())
	if err != nil {
		merrors.GinReturnError(c, err)
		return
	}

	subs, err := h.SubService.GetByFilters(filter)
	if err != nil {
		slog.Error("Failed to list subscriptions", "status", merrors.ErrorsToHTTP(err), "query", c.Request.URL.RawQuery, "error", err)
		merrors.GinReturnError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"subscriptions": subs,
	})
}

// GetSum godoc
// @Summary Get sum of subscription costs
// @Description Calculate total cost of subscriptions with optional filters
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID (UUID)"
// @Param name query string false "Service name"
// @Param from query string false "Start date (MM-YYYY)"
// @Param till query string false "End date (MM-YYYY)"
// @Success 200 {object} map[string]float64 "sum"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/subscriptions/sum [get]
func (h *SubscriptionHandler) GetSum(c *gin.Context) {
	slog.Info("Getting subscription sum", "query", c.Request.URL.RawQuery)

	filter, err := models.NewSubscriptionFilterFromURL(c.Request.URL.Query())
	if err != nil {
		merrors.GinReturnError(c, err)
		return
	}

	sum, err := h.SubService.GetSum(filter)
	if err != nil {
		slog.Error("Failed to calculate subscription sum", "status", merrors.ErrorsToHTTP(err), "query", c.Request.URL.RawQuery, "error", err)

		merrors.GinReturnError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"sum": sum})
}
