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

func ListRestaurant(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := appCtx.GetMainDBConnection()

		var pagingData common.Paging

		if err := c.ShouldBind(&pagingData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		pagingData.Fulfill()

		var filter restaurantmodel.Filter

		if err := c.ShouldBind(&filter); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		filter.Status = []int{1}

		store := restaurantstorage.NewSQLStore(db)
		business := restaurantbusiness.NewListRestaurantBusiness(store)

		result, err := business.ListRestaurant(c.Request.Context(), &filter, &pagingData)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, pagingData, filter))
	}
}
