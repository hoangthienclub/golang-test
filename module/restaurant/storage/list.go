package restaurantstorage

import (
	"context"
	"test/common"
	restaurantmodel "test/module/restaurant/model"
)

// cách tăng tải: không dùng cho dữ liệu thường xuyên thay đổi
// Step 1: fetch 1000 ids => đưa vào catching
// Step 2: fetch 50 first ids => dage (page 1)
// Step 3: Page 2, omit 50 ids and fetch the next 50 ids

func (s *sqlStore) ListDataWithCondition(
	context context.Context,
	filter *restaurantmodel.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]restaurantmodel.Restaurant, error) {
	var result []restaurantmodel.Restaurant

	db := s.db.Table(restaurantmodel.Restaurant{}.TableName())

	if f := filter; f != nil {
		if f.OwnerId > 0 {
			db = db.Where("owner_id = ?", f.OwnerId)
		}

		if len(f.Status) > 0 {
			db = db.Where("status in (?)", f.Status)
		}
	}

	if err := db.Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	if v := paging.FakeCursor; v != "" {
		uid, err := common.FromBase58(v)

		if err != nil {
			return nil, common.ErrDB(err)
		}

		db = db.Where("id < ?", uid.GetLocalID())
	} else {
		offset := (paging.Page - 1) * paging.Limit
		db = db.Offset(offset)
	}

	if err := db.
		Limit(paging.Limit).
		Order("id desc").
		Find(&result).Error; err != nil {
		return nil, err
	}

	if len(result) > 0 {
		last := result[len(result)-1]
		last.Mask(false)
		paging.NextCursor = last.FakeId.String()
	}

	return result, nil
}
