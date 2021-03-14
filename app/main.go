package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	_productHttpDelivery "github.com/wdwiramadhan/bookhub-api/product/delivery/http"
	_productHttpDeliveryMiddleware "github.com/wdwiramadhan/bookhub-api/product/delivery/http/middleware"
	_productRepo "github.com/wdwiramadhan/bookhub-api/product/repository/mysql"
	_productUcase "github.com/wdwiramadhan/bookhub-api/product/usecase"

	_authorHttDelivery "github.com/wdwiramadhan/bookhub-api/author/delivery/http"
	_authorRepo "github.com/wdwiramadhan/bookhub-api/author/repository/mysql"
	_authorUcase "github.com/wdwiramadhan/bookhub-api/author/usecase"
)

func main() {
	fmt.Println(os.Getenv("APP_ENV"))
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "5000"
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := _productHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	pr := _productRepo.NewMysqlProductRepository(dbConn)
	ar := _authorRepo.NewMysqlAuthorRepository(dbConn)

	timeoutContext := time.Duration(2) * time.Second
	pu := _productUcase.NewProductUsecase(pr, timeoutContext)
	_productHttpDelivery.NewProductHandler(e, pu)
	au := _authorUcase.NewAuthorUsecase(ar, timeoutContext)
	_authorHttDelivery.NewAuthorHandler(e, au)
	e.Logger.Fatal(e.Start(":" + Port))
}
