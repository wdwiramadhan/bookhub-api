package http

import (
	"net/http"
	"strconv"

	validator "gopkg.in/go-playground/validator.v9"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/wdwiramadhan/bookhub-api/domain"
	"github.com/wdwiramadhan/bookhub-api/helper/response"
)

var successResponse response.ResponseSuccess = response.ResponseSuccess{Success: true, Data: nil}
var failedResponse response.ResponseFailed = response.ResponseFailed{Success: false, Message: ""}

// ProductHandler  represent the httphandler for product
type ProductHandler struct {
	PUsecase domain.ProductUseCase
}

// NewProductHandler will initialize the product/ resources endpoint
func NewProductHandler(e *echo.Echo, us domain.ProductUseCase) {
	handler := &ProductHandler{
		PUsecase: us,
	}
	e.GET("/product", handler.FetchProduct)
	e.POST("/product", handler.Store)
	e.GET("/product/:productId", handler.GetByID)
	e.PUT("/product/:productId", handler.Update)
	e.DELETE("/product/:productId", handler.Delete)
}

// FetchProduct will fetch the article based on given params
func (p *ProductHandler) FetchProduct(c echo.Context) error {
	ctx := c.Request().Context()
	listProduct, err := p.PUsecase.Fetch(ctx)
	if err != nil {
		failedResponse.Message = err.Error()
		return c.JSON(getStatusCode(err), failedResponse)
	}
	successResponse.Data = listProduct
	return c.JSON(http.StatusOK, successResponse)
}

// Store will store the article by given request body
func (p *ProductHandler) Store(c echo.Context) (err error) {
	var product domain.Product
	err = c.Bind(&product)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&product); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = p.PUsecase.Store(ctx, &product)
	if err != nil {
		failedResponse.Message = err.Error()
		return c.JSON(getStatusCode(err), failedResponse)
	}
	return c.JSON(http.StatusCreated, successResponse)
}

// GetByID will get product by given id
func (p *ProductHandler) GetByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("productId"))
	ctx := c.Request().Context()
	product, err := p.PUsecase.GetByID(ctx, id)
	if err != nil {
		failedResponse.Message = err.Error()
		return c.JSON(getStatusCode(err), failedResponse)
	}
	successResponse.Data = product
	return c.JSON(http.StatusOK, successResponse)
}

// Update will update the product by given request body and params id
func (p *ProductHandler) Update(c echo.Context) (err error) {
	id, _ := strconv.Atoi(c.Param("productId"))
	var product domain.Product
	err = c.Bind(&product)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	var ok bool
	if ok, err = isRequestValid(&product); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	err = p.PUsecase.Update(ctx, &product, id)
	if err != nil {
		failedResponse.Message = err.Error()
		return c.JSON(getStatusCode(err), failedResponse)
	}

	return c.JSON(http.StatusOK, successResponse)
}

// Delete will delete product by given param
func (p *ProductHandler) Delete(c echo.Context) (err error) {
	id, _ := strconv.Atoi(c.Param("productId"))
	ctx := c.Request().Context()
	err = p.PUsecase.Delete(ctx, id)
	if err != nil {
		failedResponse.Message = err.Error()
		return c.JSON(getStatusCode(err), failedResponse)
	}
	return c.JSON(http.StatusOK, successResponse)
}

func isRequestValid(m *domain.Product) (bool, error) {
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
