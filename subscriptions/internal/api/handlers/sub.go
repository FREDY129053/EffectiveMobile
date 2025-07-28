package handlers

import (
	"log"
	"net/http"
	"strconv"
	"subscriptions/rest-service/internal/schemas"
	"subscriptions/rest-service/internal/service"
	"subscriptions/rest-service/pkg/helpers"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("mm_yyyy_date", helpers.ValidateDateMMYYYYFormatValidator)
}

func checkStartDateBeforeEndDate(startDate, endDate string) bool {
	startDateDate, _ := time.Parse("01-2006", startDate)
	endDateDate, _ := time.Parse("01-2006", endDate)

	if endDateDate.Before(startDateDate) {
		return false
	}

	return true
}

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
// @Failure 	500 	{object}  	schemas.APIError
// @Router 		/subs	[get]
func (h *SubHandler) GetAllSubscriptions(c *gin.Context) {
	res, err := h.service.GetAllSubs()
	if err != nil {
		if serviceErr, ok := err.(*schemas.AppError); ok {
			c.JSON(serviceErr.Code, gin.H{"error": serviceErr.Message})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, res)
}

// GetSubscriptionByID	godoc
// @Summary 	Get subscription info
// @Description Get subscription from database by id
// @Tags		Subs
// @Produce		json
// @Param       id    	path     	uint  	true  	"Sub ID"	Format(uint)
// @Success 	200 	{object} 	schemas.FullSubInfo
// @Failure 	400 	{object}  	schemas.APIError
// @Failure 	404 	{object}  	schemas.APIError
// @Failure 	500 	{object}  	schemas.APIError
// @Router 		/subs/{id} 	[get]
func (h *SubHandler) GetSubscriptionByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid unsigned integer parameter"})
		return
	}

	res, err := h.service.GetSub(uint(id))
	if err != nil {
		if serviceErr, ok := err.(*schemas.AppError); ok {
			c.JSON(serviceErr.Code, gin.H{"error": serviceErr.Message})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, res)
}

// CreateSubscription	godoc
// @Summary 	Create subscription
// @Description Create new subscription record
// @Tags		Subs
// @Accept		json
// @Produce 	json
// @Param       newSubscription   	body     	schemas.CreateSub 	true  	"Subscription data"
// @Success 	201 	{object} 	schemas.CreateReturn
// @Failure 	400 	{object}  	schemas.APIError
// @Failure 	500 	{object}  	schemas.APIError
// @Router 		/subs 	[post]
func (h *SubHandler) CreateSubscription(c *gin.Context) {
	var newSub schemas.CreateSub

	if err := c.ShouldBindJSON(&newSub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription data"})
		return
	}

	if err := validate.Struct(newSub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format input (must be 'mm-yyyy')"})
		return
	}

	if newSub.EndDate != nil {
		if !checkStartDateBeforeEndDate(newSub.StartDate, *newSub.EndDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "startDate cannot be after endDate"})
			return
		}
	}

	res, err := h.service.CreateSub(newSub)
	if err != nil {
		if serviceErr, ok := err.(*schemas.AppError); ok {
			c.JSON(serviceErr.Code, gin.H{"error": serviceErr.Message})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": res})
}

// FullUpdateSubscription	godoc
// @Summary 	Update subscription
// @Description Update all fields of subscription
// @Tags		Subs
// @Accept		json
// @Produce 	json
// @Param       id    				path    uint  	true  	"Subscription ID"	Format(uint)
// @Param       updateFields    	body    schemas.FullUpdateSub  	true  	"Subscription data"
// @Success 	200					{object} 	schemas.MessageReturn
// @Failure 	400 				{object}  	schemas.APIError
// @Failure 	500 				{object}  	schemas.APIError
// @Router 		/subs/{id} 	[put]
func (h *SubHandler) FullUpdateSubscription(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid unsigned integer parameter"})
		return
	}

	var subFields schemas.FullUpdateSub

	if err := c.ShouldBindJSON(&subFields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params to update subscription"})
		return
	}

	if err := validate.Struct(subFields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format input (must be 'mm-yyyy')"})
		return
	}

	if subFields.EndDate != nil {
		if !checkStartDateBeforeEndDate(subFields.StartDate, *subFields.EndDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "startDate cannot be after endDate"})
			return
		}
	}

	err = h.service.FullUpdateSub(uint(id), subFields)
	if err != nil {
		if serviceErr, ok := err.(*schemas.AppError); ok {
			c.JSON(serviceErr.Code, gin.H{"error": serviceErr.Message})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "subscription updated"})
}

// PatchUpdateSubscription	godoc
// @Summary 	Update subscription
// @Description Update the passed fields of subscription
// @Tags		Subs
// @Accept		json
// @Produce 	json
// @Param       id    				path    uint  	true  	"Subscription ID"	Format(uint)
// @Param       updateFields    	body    schemas.PatchUpdateSub  	true  	"Subscription data"
// @Success 	200 				{object} 	schemas.MessageReturn
// @Failure 	400 				{object}  	schemas.APIError
// @Failure 	500 				{object}  	schemas.APIError
// @Router 		/subs/{id} 	[patch]
func (h *SubHandler) PatchUpdateSubscription(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid unsigned integer parameter"})
		return
	}

	var subFields schemas.PatchUpdateSub

	if err := c.ShouldBindJSON(&subFields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params to update subscription"})
		return
	}

	if err := validate.Struct(subFields); err != nil {
		log.Printf("%T\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format input (must be 'mm-yyyy')"})
		return
	}

	if subFields.EndDate != nil && subFields.StartDate != nil {
		if !checkStartDateBeforeEndDate(*subFields.StartDate, *subFields.EndDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "startDate cannot be after endDate"})
			return
		}
	}

	err = h.service.PatchUpdateSub(uint(id), subFields)
	if err != nil {
		if serviceErr, ok := err.(*schemas.AppError); ok {
			c.JSON(serviceErr.Code, gin.H{"error": serviceErr.Message})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "subscription updated"})
}

// DeleteSubscription	godoc
// @Summary 	Delete subscription
// @Description Delete subscription from database
// @Tags		Subs
// @Produce 	json
// @Param       id    	path     	uint  	true  	"Subscription ID"	Format(uint)
// @Success 	200 	{object} 	schemas.MessageReturn
// @Failure 	400 	{object}  	schemas.APIError
// @Failure 	404 	{object}  	schemas.APIError
// @Failure 	500 	{object}  	schemas.APIError
// @Router 		/subs/{id} 	[delete]
func (h *SubHandler) DeleteSubscription(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid unsigned integer parameter"})
		return
	}

	err = h.service.DeleteSub(uint(id))
	if err != nil {
		if serviceErr, ok := err.(*schemas.AppError); ok {
			c.JSON(serviceErr.Code, gin.H{"error": serviceErr.Message})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "subscription deleted"})
}

// GetSubscriptionSumInfo	godoc
// @Summary 	Get subscription price
// @Description Get subscription price for period and filtered by userID or(and) serviceName
// @Tags		Subs
// @Produce 	json
// @Param       startDate    	query     	string  	true  	"Period start date('mm-yyyy')"	Format(string)
// @Param       endDate    		query     	string  	true  	"Period end date('mm-yyyy')"	Format(string)
// @Param       userID    		query     	string  	false  	"User ID"						Format(string)
// @Param       serviceName    	query     	string  	false  	"Service name"					Format(string)
// @Success 	200 	{object} 	schemas.SumReturn
// @Failure 	400 	{object}  	schemas.APIError
// @Failure 	422 	{object}  	schemas.APIError
// @Router 		/subs/sub_sum 	[get]
func (h *SubHandler) GetSubscriptionSumInfo(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	userIDInput := c.Query("userID")
	serviceNameInput := c.Query("serviceName")

	if !helpers.ValidateDateMMYYYYFormat(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}
	if !helpers.ValidateDateMMYYYYFormat(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}

	if !checkStartDateBeforeEndDate(startDate, endDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "startDate cannot be after endDate"})
		return
	}

	var userID *uuid.UUID
	if userIDInput != "" {
		userIDParse, err := uuid.Parse(userIDInput)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}
		userID = &userIDParse
	}

	resultSum, err := h.service.GetSubSum(userID, &serviceNameInput, startDate, endDate)
	if err != nil {
		if serviceErr, ok := err.(*schemas.AppError); ok {
			c.JSON(serviceErr.Code, gin.H{"error": serviceErr.Message})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"total_sum": resultSum})
}
