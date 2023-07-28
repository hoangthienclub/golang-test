package ginrestaurant

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test/common"
	"test/lib/appctx"
	restaurantbusiness "test/module/restaurant/business"
	restaurantstorage "test/module/restaurant/storage"
)

func DeleteRestaurant(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := appCtx.GetMainDBConnection()

		//id, err := strconv.Atoi(c.Param("id"))

		uid, err := common.FromBase58(c.Param("id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := restaurantstorage.NewSQLStore(db)
		business := restaurantbusiness.NewDeleteRestaurantBusiness(store)

		if err := business.DeleteRestaurant(c.Request.Context(), int(uid.GetLocalID())); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
