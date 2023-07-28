package ginrestaurant

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test/common"
	"test/lib/appctx"
	restaurantbusiness "test/module/restaurant/business"
	restaurantmodel "test/module/restaurant/model"
	restaurantstorage "test/module/restaurant/storage"
)

func CreateRestaurant(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := appCtx.GetMainDBConnection()

		//go func() {
		//	defer common.AppRecover()
		//
		//	arr := []int{}
		//	log.Println(arr[0])
		//}()

		var data restaurantmodel.RestaurantCreate
		if err := c.ShouldBind(&data); err != nil {
			//c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			panic(err)
		}

		store := restaurantstorage.NewSQLStore(db)
		business := restaurantbusiness.NewCreateRestaurantBusiness(store)

		if err := business.CreateRestaurant(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
	}
}
