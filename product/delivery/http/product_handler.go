package http 

import (
		"net/http"
		validator "gopkg.in/go-playground/validator.v9"
	
		"github.com/labstack/echo/v4"
		"github.com/sirupsen/logrus"
		"github.com/wdwiramadhan/bookhub-api/domain"
)

type ResponseError struct{
	Success bool `json:"success"`
	Message string `json:"message"`
}

type ResponseSuccess struct{
	Success bool `json:"success"`
	Data interface{} `json:"data"`
}

type ProductHandler struct {
	PUsecase domain.ProductUseCase
}

func NewProductHandler(e *echo.Echo, us domain.ProductUseCase){
	handler := &ProductHandler{
		PUsecase: us,
	}
	e.GET("/products", handler.FetchProduct)
	e.POST("/product", handler.Store)
	e.GET("/product/:id", handler.GetById)
	e.PUT("/product/:id", handler.Update)
	e.DELETE("/product/:id", handler.Delete)
}

func (p *ProductHandler) FetchProduct(c echo.Context) error {
	ctx := c.Request().Context()
	listProduct, err := p.PUsecase.Fetch(ctx)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Success: false,Message: err.Error()})
	}
	return c.JSON(http.StatusOK, ResponseSuccess {Success: true, Data: listProduct})
}

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
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated,  ResponseSuccess {Success: true, Data: product})
}

func (p *ProductHandler) GetById(c echo.Context)  error{
	id := c.Param("id")
	ctx := c.Request().Context()
	product, err := p.PUsecase.GetById(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, product)
}

func (p *ProductHandler) Update(c echo.Context) (err error){
	id := c.Param("id")
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
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, ResponseSuccess {Success: true, Data: nil})
}

func (p *ProductHandler) Delete(c echo.Context) (err error){
	id := c.Param("id")
	ctx := c.Request().Context()
	err = p.PUsecase.Delete(ctx, id)
	if err != nil{
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, ResponseSuccess {Success: true, Data: nil})
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