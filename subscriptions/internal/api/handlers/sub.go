package handlers

import (
	"net/http"
	"strconv"
	"subscriptions/rest-service/internal/schemas"
	"subscriptions/rest-service/internal/service"

	"github.com/gin-gonic/gin"
)

type SubHandler struct {
	service service.SubscriptionService
}

func NewHandler(serviceInput service.SubscriptionService) SubHandler {
	return SubHandler{
		service: serviceInput,
	}
}

// GetAllSubscriptions	godoc
// @Summary 	Get subscriptions
// @Description Get all subscriptions from database
// @Tags		Subs
// @Produce		json
// @Success 	200 	{object} 	[]schemas.FullSubInfo
// @Failure 	500 	{object}  	map[string]string
// @Router 		/subs	[get]
func (h *SubHandler) GetAllSubscriptions(c *gin.Context) {
	res, err := h.service.GetAllSubs()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, res)
}

// GetSubscriptionByID	godoc
// @Summary 	Get subscription info
// @Description Get subscription from database by id
// @Tags		Subs
// @Produce		json
// @Param       id    	path     	uint  	true  	"Sub ID"	Format(uint)
// @Success 	200 	{object} 	schemas.FullSubInfo
// @Failure 	400 	{object}  	map[string]string
// @Failure 	404 	{object}  	map[string]string
// @Failure 	500 	{object}  	map[string]string
// @Router 		/subs/{id} 	[get]
func (h *SubHandler) GetSubscriptionByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid unsigned integer parameter"})
		return
	}

	res, err := h.service.GetSub(uint(id))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, res)
}

// CreateSubscription	godoc
// @Summary 	Create subscription
// @Description Create new subscription record
// @Tags		Subs
// @Accept		json
// @Produce 	json
// @Param       newSubscription   	body     	schemas.CreateSub 	true  	"Subscription data"
// @Success 	201 	{object} 	map[string]uint
// @Failure 	400 	{object}  	map[string]string
// @Failure 	500 	{object}  	map[string]string
// @Router 		/subs 	[post]
func (h *SubHandler) CreateSubscription(c *gin.Context) {
	var newSub schemas.CreateSub

	if err := c.ShouldBindJSON(&newSub); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid subscription data"})
		return
	}

	res, err := h.service.CreateSub(newSub)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"id": res})
}

// FullUpdateSubscription	godoc
// @Summary 	Update subscription
// @Description Update all fields of subscription
// @Tags		Subs
// @Accept		json
// @Produce 	json
// @Param       id    				path    uint  	true  	"Subscription ID"	Format(uint)
// @Param       updateFields    	body    schemas.FullUpdateSub  	true  	"Subscription data"
// @Success 	200					{object} 	map[string]string
// @Failure 	400 				{object}  	map[string]string
// @Failure 	500 				{object}  	map[string]string
// @Router 		/subs/{id} 	[put]
func (h *SubHandler) FullUpdateSubscription(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid unsigned integer parameter"})
		return
	}

	var subFields schemas.FullUpdateSub

	if err := c.ShouldBindJSON(&subFields); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid params to update subscription"})
		return
	}

	if err := h.service.FullUpdateSub(uint(id), subFields); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "cannot update subscription"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "subscription updated"})
}

// PatchUpdateSubscription	godoc
// @Summary 	Update subscription
// @Description Update the passed fields of subscription
// @Tags		Subs
// @Accept		json
// @Produce 	json
// @Param       id    				path    uint  	true  	"Subscription ID"	Format(uint)
// @Param       updateFields    	body    schemas.PatchUpdateSub  	true  	"Subscription data"
// @Success 	200 				{object} 	map[string]string
// @Failure 	400 				{object}  	map[string]string
// @Failure 	500 				{object}  	map[string]string
// @Router 		/subs/{id} 	[patch]
func (h *SubHandler) PatchUpdateSubscription(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid unsigned integer parameter"})
		return
	}

	var subFields schemas.PatchUpdateSub

	if err := c.ShouldBindJSON(&subFields); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid params to update subscription"})
		return
	}

	if err := h.service.PatchUpdateSub(uint(id), subFields); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "cannot update subscription"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "subscription updated"})
}

// DeleteSubscription	godoc
// @Summary 	Delete subscription
// @Description Delete subscription from database
// @Tags		Subs
// @Produce 	json
// @Param       id    	path     	uint  	true  	"Subscription ID"	Format(uint)
// @Success 	200 	{object} 	map[string]string
// @Failure 	400 	{object}  	map[string]string
// @Failure 	404 	{object}  	map[string]string
// @Failure 	500 	{object}  	map[string]string
// @Router 		/subs/{id} 	[delete]
func (h *SubHandler) DeleteSubscription(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid unsigned integer parameter"})
		return
	}

	if err := h.service.DeleteSub(uint(id)); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "subscription deleted"})
}