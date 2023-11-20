package services

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/BabyBoChen/pgdbcontext"
)

type CuisineService struct {
	db *pgdbcontext.DbContext
}

func NewCuisineService() (*CuisineService, error) {
	service := CuisineService{}
	var err error
	envVars := ReadEnvironmentVariables()
	service.db, err = pgdbcontext.NewDbContext(envVars.ConnStr)
	return &service, err
}

func NewCuisineServiceWithApplicationName(applicationName string) (*CuisineService, error) {
	service := CuisineService{}
	var err error
	envVars := ReadEnvironmentVariables()
	connStr := envVars.ConnStr + " application_name=" + applicationName
	service.db, err = pgdbcontext.NewDbContext(connStr)
	return &service, err
}

type DataSet = map[string][]map[string]interface{}

func (service *CuisineService) GetTop10Cuisines() (DataSet, error) {
	var err error
	ds := make(DataSet)
	var top10Main []map[string]interface{}
	var top10Dessert []map[string]interface{}
	var top10Buffet []map[string]interface{}
	var top10NewCuisine []map[string]interface{}

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
		ds["Top10Main"] = top10Main
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
		ds["Top10Dessert"] = top10Dessert
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
		ds["Top10Buffet"] = top10Buffet
		sql = `
		SELECT A.cuisine_id, A.cuisine_name, A.unit_price,B.cuisine_type_name
		,CASE WHEN A.is_one_set = true
			THEN 'YES'
			ELSE 'NO'
			END AS is_one_set
		,A.review,A.last_order_date,A.restaurant,A.address,A.remark
		FROM public.cuisine AS A
		LEFT JOIN public.cuisine_type AS B ON A.cuisine_type=B.cuisine_type_id
		ORDER BY A.cuisine_id DESC
		LIMIT 10`
		dt, err = service.db.Query(sql)
	}

	if err == nil {
		top10NewCuisine = toSliceMap(dt)
		ds["Top10NewCuisine"] = top10NewCuisine
	}

	return ds, err
}

func toSliceMap(dt *pgdbcontext.DataTable) []map[string]interface{} {
	sliceMap := make([]map[string]interface{}, len(dt.Rows))
	for i, row := range dt.Rows {
		rowDict := row.ToMap()
		sliceMap[i] = rowDict
	}
	return sliceMap
}

func (service *CuisineService) SaveNewCuisine(newCuisine map[string]interface{}) error {
	repo, err := service.db.GetRepository("cuisine")
	if err == nil {
		_, err = repo.Insert(newCuisine)
	}
	if err == nil {
		err = service.db.Commit()
	}

	return err
}

func (service *CuisineService) ListAllCuisine() ([]map[string]interface{}, error) {
	var err error
	var allCuisine []map[string]interface{}
	var dt *pgdbcontext.DataTable
	sql := `
	SELECT A.cuisine_id, A.cuisine_name, A.unit_price,B.cuisine_type_name
	,CASE WHEN A.is_one_set = true
		THEN 'YES'
		ELSE 'NO'
		END AS is_one_set
	,A.review,A.last_order_date,A.restaurant,A.address,A.remark
	FROM public.cuisine AS A
	LEFT JOIN public.cuisine_type AS B ON A.cuisine_type=B.cuisine_type_id
	ORDER BY last_order_date DESC`
	dt, err = service.db.Query(sql)
	if err == nil {
		allCuisine = toSliceMap(dt)
	}
	return allCuisine, err
}

func (service *CuisineService) GetApplicationName() error {
	dt, err := service.db.Query("SELECT current_setting('application_name')")
	if err == nil {
		for _, row := range dt.Rows {
			fmt.Println(row.ToMap())
		}
	}
	return err
}

func (service *CuisineService) QueryWithKeyword(keyword string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	var err error

	sql := `with VI_Cuisine as 
	(
		SELECT A.cuisine_id, A.cuisine_name, A.unit_price,B.cuisine_type_name
		,CASE WHEN A.is_one_set = true
			THEN 'YES'
			ELSE 'NO'
			END AS is_one_set
		,A.review,A.last_order_date,A.restaurant,A.address,A.remark
		FROM public.cuisine AS A
		LEFT JOIN public.cuisine_type AS B ON A.cuisine_type=B.cuisine_type_id
	)
	select *
	from VI_Cuisine
	%s
	ORDER BY last_order_date desc`
	whereSql := ""
	if len(keyword) > 0 {
		whereSql = "where cuisine_name like $1 OR cuisine_type_name like $1 OR cast(last_order_date AS VARCHAR(10)) like $1 OR restaurant like $1 OR address like $1 OR remark like $1 "
	}
	sql = fmt.Sprintf(sql, whereSql)

	var dt *pgdbcontext.DataTable
	if len(keyword) > 0 {
		dt, err = service.db.Query(sql, "%"+keyword+"%")
	} else {
		dt, err = service.db.Query(sql)
	}

	if err == nil {
		results = toSliceMap(dt)
	}
	return results, err
}

//query by fields

func (service *CuisineService) GetCuisineByCuisineId(cuisineId string) (map[string]interface{}, error) {
	var results map[string]interface{}
	var err error

	var cid int64
	cid, err = strconv.ParseInt(cuisineId, 10, 32)

	var dt *pgdbcontext.DataTable
	if err == nil {
		sql := `SELECT *
		FROM public.cuisine
		WHERE cuisine_id=$1`
		dt, err = service.db.Query(sql, cid)
	}

	var arr []map[string]interface{}
	if err == nil {
		arr = toSliceMap(dt)
		if len(arr) == 1 {
			results = arr[0]
		} else {
			err = errors.New("not found")
		}
	}

	return results, err
}

func (service *CuisineService) SaveCuisine(cuisine map[string]interface{}) error {
	repo, err := service.db.GetRepository("cuisine")

	var dt *pgdbcontext.DataTable
	if err == nil {
		dt, err = repo.Select("cuisine_id=$1", cuisine["cuisine_id"])
	}

	if err == nil {
		if len(dt.Rows) != 1 {
			err = errors.New("not found")
		}
	}

	if err == nil {
		err = repo.Update(cuisine)
	}

	if err == nil {
		err = service.db.Commit()
	}

	return err
}

func (service *CuisineService) Dispose() {
	if service.db != nil {
		service.db.Dispose()
	}
}
