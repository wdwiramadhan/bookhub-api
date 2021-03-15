package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/wdwiramadhan/bookhub-api/domain"
	"github.com/wdwiramadhan/bookhub-api/helper/response"
	"gopkg.in/go-playground/validator.v9"
)

var successResponse response.ResponseSuccess = response.ResponseSuccess{Success: true, Data: nil}
var failedResponse response.ResponseFailed = response.ResponseFailed{Success: false, Message: ""}

// AuthorHandler represent the httphandler for product
type AuthorHandler struct {
	AUsecase domain.AuthorUsecase
}

// NewAuthorHandler will initialize the author endpoint
func NewAuthorHandler(e *echo.Echo, us domain.AuthorUsecase) {
	handler := &AuthorHandler{
		AUsecase: us,
	}
	e.GET("/author", handler.Fetch)
	e.POST("/author", handler.Store)
	e.GET("/author/:authorId", handler.GetAuthorById)
	e.PUT("/author/:authorId", handler.UpdateAuthorById)
	e.DELETE("/author/:authorId", handler.DeleteAuthorById)
}

// Fetch will fetch the author based on given params
func (a *AuthorHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()
	authors, err := a.AUsecase.Fetch(ctx)
	if err != nil {
		failedResponse.Message = err.Error()
		return c.JSON(getStatusCode(err), failedResponse)
	}
	successResponse.Data = authors
	return c.JSON(http.StatusOK, successResponse)
}

func (a *AuthorHandler) Store(c echo.Context) (err error) {
	var author domain.Author
	err = c.Bind(&author)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	var ok bool
	if ok, err = isRequestValid(&author); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	err = a.AUsecase.Store(ctx, &author)
	if err != nil {
		failedResponse.Message = err.Error()
		return c.JSON(getStatusCode(err), failedResponse)
	}

	return c.JSON(http.StatusCreated, successResponse)
}

func (a *AuthorHandler) GetAuthorById(c echo.Context) (err error) {
	ctx := c.Request().Context()
	authorId, _ := strconv.Atoi(c.Param("authorId"))
	author, err := a.AUsecase.GetAuthorById(ctx, authorId)
	if err != nil {
		failedResponse.Message = err.Error()
		return c.JSON(getStatusCode(err), failedResponse)
	}
	successResponse.Data = author
	return c.JSON(http.StatusOK, successResponse)
}

func (a *AuthorHandler) UpdateAuthorById(c echo.Context) (err error) {
	authorId, _ := strconv.Atoi(c.Param("authorId"))
	var author domain.Author
	err = c.Bind(&author)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	ctx := c.Request().Context()
	err = a.AUsecase.UpdateAuthorById(ctx, authorId, &author)
	if err != nil {
		failedResponse.Message = err.Error()
		return c.JSON(getStatusCode(err), failedResponse)
	}
	return c.JSON(http.StatusCreated, successResponse)
}

func (a *AuthorHandler) DeleteAuthorById(c echo.Context) (err error) {
	authorId, _ := strconv.Atoi(c.Param("authorId"))
	ctx := c.Request().Context()
	err = a.AUsecase.DeleteAuthorById(ctx, authorId)
	if err != nil {
		failedResponse.Message = err.Error()
		return c.JSON(getStatusCode(err), failedResponse)
	}
	return c.JSON(http.StatusCreated, successResponse)
}

func isRequestValid(m *domain.Author) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
