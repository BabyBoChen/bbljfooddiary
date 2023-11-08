package services

import (
	"github.com/BabyBoChen/pgdbcontext"
)

type CuisineService struct {
	db *pgdbcontext.DbContext
}

func NewCuisineService() (CuisineService, error) {
	service := CuisineService{}
	var err error
	envVars := ReadEnvironmentVariables()
	service.db, err = pgdbcontext.NewDbContext(envVars.ConnStr)
	return service, err
}

func (service *CuisineService) GetTop10Cuisines() ([]map[string]interface{}, error) {
	var top10Cuisines []map[string]interface{}
	var err error

	sql := `
	SELECT cuisine_id, cuisine_name, unit_price
	,CASE WHEN is_one_set = true
		THEN 'YES'
		ELSE 'NO'
		END AS is_one_set
	,review,last_order_date,restaurant,address,remark
	FROM public.cuisine
	ORDER BY review DESC, cuisine_name ASC
	LIMIT 10`

	var dt *pgdbcontext.DataTable
	dt, err = service.db.Query(sql)

	if err == nil {
		top10Cuisines = toSliceMap(dt)
	}

	return top10Cuisines, err
}

func toSliceMap(dt *pgdbcontext.DataTable) []map[string]interface{} {
	sliceMap := make([]map[string]interface{}, len(dt.Rows))
	for i, row := range dt.Rows {
		rowDict := row.ToMap()
		sliceMap[i] = rowDict
	}
	return sliceMap
}

func (service *CuisineService) Dispose() {
	if service.db != nil {
		service.db.Dispose()
	}
}
