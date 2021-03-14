package main

import (
	"fmt"
	"log"
	"time"
	"os"
	"net/url"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/go-sql-driver/mysql"

	_articleHttpDeliveryMiddleware "github.com/wdwiramadhan/bookhub-api/product/delivery/http/middleware"
	_productHttpDelivery "github.com/wdwiramadhan/bookhub-api/product/delivery/http"
	_productRepo "github.com/wdwiramadhan/bookhub-api/product/repository/mysql"
	_productUcase "github.com/wdwiramadhan/bookhub-api/product/usecase"
)

func main(){
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
	middL := _articleHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	pr := _productRepo.NewMysqlProductRepository(dbConn)

	timeoutContext := time.Duration(2) * time.Second
	pu := _productUcase.NewProductUsecase(pr, timeoutContext)
	_productHttpDelivery.NewProductHandler(e, pu)
	e.Logger.Fatal(e.Start(":"+Port))
}