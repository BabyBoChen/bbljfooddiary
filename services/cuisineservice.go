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

func (service *CuisineService) GetTop10Cuisines() ([]map[string]interface{}, []map[string]interface{}, []map[string]interface{}, error) {
	var err error
	var top10Main []map[string]interface{}
	var top10Dessert []map[string]interface{}
	var top10Buffet []map[string]interface{}

	var dt *pgdbcontext.DataTable
	var sql string
	if err == nil {
		sql = `
		SELECT A.cuisine_id, A.cuisine_name, A.unit_price,B.cuisine_type_name
		,CASE WHEN A.is_one_set = true
			THEN 'YES'
			ELSE 'NO'
			END AS is_one_set
		,A.review,A.last_order_date,A.restaurant,A.address,A.remark
		FROM public.cuisine AS A
		LEFT JOIN public.cuisine_type AS B ON A.cuisine_type=B.cuisine_type_id
		WHERE A.cuisine_type=1
		ORDER BY review DESC, cuisine_name ASC
		LIMIT 10`
		dt, err = service.db.Query(sql)
	}

	if err == nil {
		top10Main = toSliceMap(dt)
		sql = `
		SELECT A.cuisine_id, A.cuisine_name, A.unit_price,B.cuisine_type_name
		,CASE WHEN A.is_one_set = true
			THEN 'YES'
			ELSE 'NO'
			END AS is_one_set
		,A.review,A.last_order_date,A.restaurant,A.address,A.remark
		FROM public.cuisine AS A
		LEFT JOIN public.cuisine_type AS B ON A.cuisine_type=B.cuisine_type_id
		WHERE A.cuisine_type=2
		ORDER BY review DESC, cuisine_name ASC
		LIMIT 10`
		dt, err = service.db.Query(sql)
	}

	if err == nil {
		top10Dessert = toSliceMap(dt)
		sql = `
		SELECT A.cuisine_id, A.cuisine_name, A.unit_price,B.cuisine_type_name
		,CASE WHEN A.is_one_set = true
			THEN 'YES'
			ELSE 'NO'
			END AS is_one_set
		,A.review,A.last_order_date,A.restaurant,A.address,A.remark
		FROM public.cuisine AS A
		LEFT JOIN public.cuisine_type AS B ON A.cuisine_type=B.cuisine_type_id
		WHERE A.cuisine_type=3
		ORDER BY review DESC, cuisine_name ASC
		LIMIT 10`
		dt, err = service.db.Query(sql)
	}

	if err == nil {
		top10Buffet = toSliceMap(dt)
	}

	return top10Main, top10Dessert, top10Buffet, err
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
