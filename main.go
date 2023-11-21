package main

import (
	"encoding/json"
	"io"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	// register web server
	e := echo.New()

	// allow from localhost
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
	}))

	// register db handler
	dsn := "root:@tcp(127.0.0.1:3306)/pos?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	} else {
		log.Println("Connected successfully")
	}

	// define model
	type Invoice struct {
		gorm.Model
		WorkspaceName    string
		SubscriptionPlan string
		BillingAmount    float64
		BillingPeriod    string
		PONumber         string `gorm:"size:15"`
	}

	e.GET("/", func(c echo.Context) error {

		db.AutoMigrate(Invoice{})

		return c.JSON(200, "hellow world")

	})

	e.POST("/save-po", func(c echo.Context) error {

		defer c.Request().Body.Close()
		jsonbody, err := ParseRequestBody(c.Request().Body)

		if err != nil {
			return c.JSON(500, err.Error())
		}

		amount, err := strconv.ParseFloat(jsonbody["billingAmount"].(string), 64)

		if err != nil {
			return c.JSON(500, "amount cannot be parsed to a valid number.")
		}

		invoice := Invoice{
			WorkspaceName:    jsonbody["workspaceName"].(string),
			SubscriptionPlan: jsonbody["subscriptionPlan"].(string),
			BillingAmount:    amount,
			BillingPeriod:    `["2023-01-11", "2023-01-11"]`,
			PONumber:         jsonbody["poNumber"].(string),
		}

		err = db.Create(&invoice).Error

		if err != nil {
			return c.JSON(500, err.Error())
		}

		return c.JSON(200, "success")

	})

	e.Logger.Fatal(e.Start(":1323"))
}

func ParseRequestBody(body io.Reader) (map[string]interface{}, error) {
	var jsonBody map[string]interface{}
	err := json.NewDecoder(body).Decode(&jsonBody)

	if err != nil {
		return nil, err
	} else {
		return jsonBody, nil
	}
}
