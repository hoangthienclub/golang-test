package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"test/lib/appctx"
	"test/middleware"
	"test/module/restaurant/transport/ginrestaurant"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Restaurant struct {
	Id   int    `json:"id" gorm:"column:id;"` //tag
	Name string `json:"name" gorm:"column:name;"`
	Addr string `json:"addr" gorm:"column:addr;"`
}

func (Restaurant) TableName() string { return "restaurants" }

type RestaurantUpdate struct {
	Name *string `json:"name" gorm:"column:name;"`
	Addr *string `json:"addr" gorm:"column:addr;"`
}

func (RestaurantUpdate) TableName() string { return Restaurant{}.TableName() }

func main() {
	//test := Restaurant{
	//	Id:   1,
	//	Name: "thien",
	//	Addr: "la",
	//}
	//
	//jsonByte, err := json.Marshal(test)
	//log.Println(string(jsonByte), err)
	//
	//json.Unmarshal([]byte("{\"id\":1,\"name\":\"thien\",\"addr\":\"la\"}"), &test)
	//log.Println(test)

	dsn := os.Getenv("MYSQL_CONN_STRING")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db = db.Debug()

	appContext := appctx.NewAppContext(db)

	r := gin.Default()
	r.Use(middleware.Recover(appContext))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v1 := r.Group("/v1")

	restaurants := v1.Group("/restaurants")

	restaurants.POST("", ginrestaurant.CreateRestaurant(appContext))

	restaurants.GET("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		var data Restaurant

		db.Where("id = ?", id).First(&data)

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})
	restaurants.GET("", ginrestaurant.ListRestaurant(appContext))
	restaurants.DELETE("/:id", ginrestaurant.DeleteRestaurant(appContext))

	r.Run()
}
