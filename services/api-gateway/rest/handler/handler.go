package handler

import (
	"errors"
	"fmt"
	"net/http"

	model3 "github.com/GEtBUsyliVn/url-shortener/services/api-gateway/model"
	"github.com/GEtBUsyliVn/url-shortener/services/api-gateway/rest"
	model2 "github.com/GEtBUsyliVn/url-shortener/services/api-gateway/rest/model"
	"github.com/GEtBUsyliVn/url-shortener/services/api-gateway/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Handler struct {
	log       *zap.Logger
	srvc      *service.GatewayService
	validator *validator.Validate
}

func NewHandler(log *zap.Logger, srvc *service.GatewayService, validator *validator.Validate) *Handler {
	return &Handler{
		log:       log.Named("rest_handler"),
		srvc:      srvc,
		validator: validator,
	}
}

func (h *Handler) CreateShortCode(c *gin.Context) {
	req := &model2.CreateShortCodeRequest{}

	if err := c.ShouldBindJSON(req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, rest.CommonError(err, "invalid request body"))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, rest.ErrValidate(err))
		return
	}

	shortCode := &model3.ShortCode{}
	shortCode.BindRestCreate(req)
	code, err := h.srvc.CreateShortUrl(c, shortCode)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, rest.CommonError(err, "failed to create short code"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": fmt.Sprintf("%s/%s", c.Request.Host, code)})
}

func (h *Handler) GetUrl(c *gin.Context) {
	param := c.Param("code")
	if param == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, rest.CommonError(nil, "code is required"))
		return
	}
	event := &model3.ClickEvent{}
	event.WriteClickEvent(c)
	h.log.Info("got click event", zap.String("short_code", param), zap.Any("event", event), zap.Any("client_ip", c.ClientIP()))

	url, err := h.srvc.GetOriginalUrl(c, param, event.BindToRequest())
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, rest.CommonError(err, "short code not found"))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, rest.CommonError(err, "failed to get url"))
		return
	}
	c.Header("Location", url)
	c.Redirect(http.StatusMovedPermanently, url)
}

func (h *Handler) GetStats(c *gin.Context) {
	param := c.Param("code")
	if param == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, rest.CommonError(nil, "code is required"))
		return
	}
	stats, err := h.srvc.GetAnalytics(c, param)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, rest.CommonError(err, "short code not found"))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, rest.CommonError(err, "failed to get stats"))
		return
	}
	c.JSON(http.StatusOK, stats.BindRest())
}
