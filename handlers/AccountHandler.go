package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hi-supergirl/go-microservice-template/handlers/services"
	"github.com/hi-supergirl/go-microservice-template/handlers/services/dto"
	"github.com/hi-supergirl/go-microservice-template/helper"
	"github.com/hi-supergirl/go-microservice-template/logging"
	"github.com/hi-supergirl/go-microservice-template/middlewares"
	"go.uber.org/zap"
)

type AccountHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Me(c *gin.Context)
}

type accountHandler struct {
	accountService services.AccountService
}

func NewAccountHandler(logger *zap.Logger, accountService services.AccountService) AccountHandler {
	return &accountHandler{accountService: accountService}
}

func (ah *accountHandler) Register(c *gin.Context) {
	execute(c, func(c *gin.Context) *Response {
		logger := logging.FromContext(c)
		logger.Infow("[accountHandler]", "Register", "")
		var auth dto.AccountDTO
		if err := c.ShouldBindJSON(&auth); err != nil {
			return NewResponse(http.StatusBadRequest, failed, err)
		}

		encodedPassword, err := helper.EncodePassword(auth.Password)
		if err != nil {
			return NewResponse(http.StatusBadRequest, failed, err)
		}
		accDto := dto.AccountDTO{UserName: auth.UserName, Password: encodedPassword}
		savedAccDto, err := ah.accountService.Save(c.Request.Context(), accDto)
		if err != nil {
			return NewResponse(http.StatusBadRequest, failed, err)
		}
		savedAccDto.Password = ""
		return NewResponse(http.StatusOK, success, savedAccDto)
	})
}

func (ah *accountHandler) Login(c *gin.Context) {
	execute(c, func(c *gin.Context) *Response {
		logger := logging.FromContext(c)
		logger.Infow("[accountHandler]", "Login", "")
		var auth dto.AccountDTO
		if err := c.ShouldBindJSON(&auth); err != nil {
			return NewResponse(http.StatusBadRequest, failed, err)
		}
		accDto, err := ah.accountService.GetByName(c.Request.Context(), auth.UserName)
		if err != nil {
			return NewResponse(http.StatusBadRequest, failed, err)
		}
		if err := accDto.ValidatePassword(auth.Password); err != nil {
			return NewResponse(http.StatusBadRequest, failed, err)
		}
		jwt, err := helper.GenerateJWT(accDto)
		if err != nil {
			return NewResponse(http.StatusBadRequest, failed, err)
		}
		return NewResponse(http.StatusOK, success, jwt)
	})
}

func (ah *accountHandler) Me(c *gin.Context) {
	execute(c, func(c *gin.Context) *Response {
		logger := logging.FromContext(c)
		logger.Infow("[accountHandler]", "Me", "'")
		curAccount, err := ah.getCurrentAccount(c)
		if err != nil {
			return NewResponse(http.StatusBadRequest, failed, err)
		}
		curAccount.Password = ""
		return NewResponse(http.StatusOK, success, curAccount)
	})

}

func (ah *accountHandler) getCurrentAccount(c *gin.Context) (*dto.AccountDTO, error) {
	currentId, err := helper.GetCurrentAccountId(c)
	if err != nil {
		return nil, err
	}
	accDto, err := ah.accountService.GetById(c.Request.Context(), currentId)
	if err != nil {
		return nil, err
	}
	return accDto, nil
}

func AccountRoute(ah AccountHandler, logger *zap.Logger, c *gin.Engine) {
	api := c.Group("/api")
	api.Use(middlewares.RequestTraceMiddleWare())
	{
		api.POST("/account/register", ah.Register)
		api.POST("/account/login", ah.Login)
	}

	api.Use(middlewares.JwtTokenMiddleWare())
	{
		api.GET("/account/me", ah.Me)
	}
}
