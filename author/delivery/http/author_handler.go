package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/wdwiramadhan/bookhub-api/domain"
	"gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ResponseSuccess represent the reseponse success struct
type ResponseSuccess struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

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
		return c.JSON(getStatusCode(err), ResponseError{Success: false, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, ResponseSuccess{Success: true, Data: authors})
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
		return c.JSON(getStatusCode(err), ResponseError{Success: false, Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, ResponseSuccess{Success: true, Data: nil})
}

func (a *AuthorHandler) GetAuthorById(c echo.Context) (err error) {
	ctx := c.Request().Context()
	authorId, _ := strconv.Atoi(c.Param("authorId"))
	author, err := a.AUsecase.GetAuthorById(ctx, authorId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Success: false, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, ResponseSuccess{Success: true, Data: author})
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
		return c.JSON(getStatusCode(err), ResponseError{Success: false, Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, ResponseSuccess{Success: true, Data: nil})
}

func (a *AuthorHandler) DeleteAuthorById(c echo.Context) (err error) {
	authorId, _ := strconv.Atoi(c.Param("authorId"))
	ctx := c.Request().Context()
	err = a.AUsecase.DeleteAuthorById(ctx, authorId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Success: false, Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, ResponseSuccess{Success: true, Data: nil})
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
