package services

import (
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/BabyBoChen/bbljfooddiary/utils"
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

func (service *CuisineService) SaveNewCuisine(newCuisine map[string]interface{}, cuisineImage io.Reader) error {
	repo, err := service.db.GetRepository("cuisine")

	var lastInsertedId map[string]interface{}
	if err == nil {
		lastInsertedId, err = repo.Insert(newCuisine)
	}
	if cuisineImage != nil {
		var dpClient *DropboxClient
		if err == nil {
			dpClient, err = NewDropboxClient()
		}
		var img image.Image
		var ext string
		if err == nil {
			img, ext, err = image.Decode(cuisineImage)
		}
		var tmpImgPath string
		if err == nil {
			tmpImgPath, err = utils.ResizeImage(img, ext, 800)
		}
		var tmpImg *os.File
		if err == nil {
			tmpImg, err = os.Open(tmpImgPath)
		}
		var id int64
		ok := false
		var id_str string
		if err == nil {
			defer os.Remove(tmpImgPath)
			defer tmpImg.Close()
			id, ok = lastInsertedId["cuisine_id"].(int64)
			if ok {
				id_str = fmt.Sprintf("%d", id)
				_, err = dpClient.UploadFile("/"+id_str, "CuisineImage.png", tmpImg)
			} else {
				err = errors.New("cannot upload image")
			}
		}
		var resp map[string]interface{}
		if err == nil {
			resp, err = dpClient.CreateSharedLink("/" + id_str + "/CuisineImage.png")
		}
		if err == nil {
			service.putCacheCuisineImageUrl(id_str, resp["url"].(string))
		}
	}
	if err == nil {
		err = service.db.Commit()
	}
	return err
}

var cacheCuisineImageUrl map[string]string

func (service *CuisineService) putCacheCuisineImageUrl(cuisineId string, sharedLink string) {
	if cacheCuisineImageUrl == nil {
		cacheCuisineImageUrl = make(map[string]string)
	}
	key := fmt.Sprintf("k_%s", cuisineId)
	cacheCuisineImageUrl[key] = sharedLink
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

type FormData = map[string]string

func (service *CuisineService) QueryWithKeyword(formData FormData) ([]map[string]interface{}, int, error) {
	var results []map[string]interface{}
	var err error

	keys := make([]string, 0, len(formData))
	for k := range formData {
		keys = append(keys, k)
	}

	keyword := formData["Keyword"]

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
	select *, COUNT(*) OVER() AS total_count
	from VI_Cuisine
	%s
	%s
	%s`

	whereSql := "where "
	if len(keyword) > 0 {
		whereSql += " cuisine_name like $1 OR cuisine_type_name like $1 OR cast(last_order_date AS VARCHAR(10)) like $1 OR restaurant like $1 OR address like $1 OR remark like $1 "
	} else {
		whereSql += "1=1 "
	}

	orderSql := "ORDER BY "
	sortField := "x"
	sortDir := "asc"
	sortI := 0
	if slices.Contains(keys, "sort_field_0") {
		for sortField != "" {
			sortField = formData[fmt.Sprintf("sort_field_%d", sortI)]
			sortDir = formData[fmt.Sprintf("sort_dir_%d", sortI)]
			if sortField == "" {
				break
			}
			if sortI > 0 {
				orderSql += ", "
			}
			orderSql += fmt.Sprintf(" %s %s ", sortField, sortDir)
			sortI += 1
		}
	} else {
		orderSql += " last_order_date desc "
	}

	limitSql := "LIMIT %d OFFSET %d"
	size, _ := strconv.Atoi(formData["size"])
	page, _ := strconv.Atoi(formData["page"])
	limitSql = fmt.Sprintf(limitSql, size, (page-1)*size)

	sql = fmt.Sprintf(sql, whereSql, orderSql, limitSql)

	var dt *pgdbcontext.DataTable
	if len(keyword) > 0 {
		dt, err = service.db.Query(sql, "%"+keyword+"%")
	} else {
		dt, err = service.db.Query(sql)
	}

	totalCount := 0
	if err == nil {
		results = toSliceMap(dt)
		if len(results) > 0 {
			i64, _ := results[0]["total_count"].(int64)
			totalCount = int(i64)
		}
	}
	return results, totalCount, err
}

func (service *CuisineService) QueryWithFields(formData FormData) ([]map[string]interface{}, int, error) {
	var results []map[string]interface{}
	var err error

	keys := make([]string, 0, len(formData))
	for k := range formData {
		keys = append(keys, k)
	}

	sql := `WITH VI_Cuisine AS 
	(
		SELECT A.cuisine_id, A.cuisine_name, A.unit_price,B.cuisine_type_name
		,CASE WHEN A.is_one_set = true
			THEN 'YES'
			ELSE 'NO'
			END AS is_one_set
		,A.review,A.last_order_date,A.restaurant,A.address,A.remark
		,A.cuisine_type
		FROM public.cuisine AS A
		LEFT JOIN public.cuisine_type AS B ON A.cuisine_type=B.cuisine_type_id
	)
	select *, COUNT(*) OVER() AS total_count
	from VI_Cuisine
	%s
	%s
	%s`

	whereSql := "WHERE 1=1 "
	paramIndex := "$1"
	parameters := make([]interface{}, 0)
	if utils.MapContainsKey[string](formData, "CuisineName") && len(formData["CuisineName"]) > 0 {
		whereSql += fmt.Sprintf(" AND cuisine_name LIKE " + paramIndex)
		parameters = append(parameters, "%"+formData["CuisineName"]+"%")
		paramIndex = "$" + fmt.Sprintf("%d", len(parameters)+1)
	}
	cuisineTypes := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
	}
	if utils.MapContainsKey[string](formData, "CuisineType") && utils.MapContainsKey[int](cuisineTypes, formData["CuisineType"]) {
		whereSql += fmt.Sprintf(" AND cuisine_type=" + paramIndex)
		parameters = append(parameters, cuisineTypes[formData["CuisineType"]])
		paramIndex = "$" + fmt.Sprintf("%d", len(parameters)+1)
	}
	hasLastOrderDate := false
	var lastOrderDate time.Time
	hasLastOrderDateTo := false
	var lastOrderDateTo time.Time
	if utils.MapContainsKey[string](formData, "LastOrderDate") && len(formData["LastOrderDate"]) > 0 {
		lastOrderDate_str := formData["LastOrderDate"]
		lastOrderDate, err = time.Parse("2006-01-02", lastOrderDate_str)
		if err == nil {
			hasLastOrderDate = true
		}
	}
	if err == nil {
		if utils.MapContainsKey[string](formData, "LastOrderDateTo") && len(formData["LastOrderDateTo"]) > 0 {
			lastOrderDateTo_str := formData["LastOrderDateTo"]
			lastOrderDateTo, err = time.Parse("2006-01-02", lastOrderDateTo_str)
			if err == nil {
				hasLastOrderDateTo = true
			}
		}
	}
	if err == nil {
		if hasLastOrderDate && hasLastOrderDateTo {
			dateRange := make([]time.Time, 2)
			if lastOrderDate.Compare(lastOrderDateTo) <= 0 {
				dateRange[0] = lastOrderDate
				dateRange[1] = lastOrderDateTo
			} else {
				dateRange[0] = lastOrderDateTo
				dateRange[1] = lastOrderDate
			}
			whereSql += " AND last_order_date >= " + paramIndex
			parameters = append(parameters, dateRange[0])
			paramIndex = "$" + fmt.Sprintf("%d", len(parameters)+1)
			whereSql += " AND last_order_date <= " + paramIndex
			parameters = append(parameters, dateRange[1])
			paramIndex = "$" + fmt.Sprintf("%d", len(parameters)+1)
		} else if hasLastOrderDate {
			whereSql += " AND last_order_date = " + paramIndex
			parameters = append(parameters, lastOrderDate)
			paramIndex = "$" + fmt.Sprintf("%d", len(parameters)+1)
		}
	}
	if err == nil {
		if utils.MapContainsKey[string](formData, "Restaurant") && len(formData["Restaurant"]) > 0 {
			whereSql += fmt.Sprintf(" AND restaurant LIKE " + paramIndex)
			parameters = append(parameters, "%"+formData["Restaurant"]+"%")
			paramIndex = "$" + fmt.Sprintf("%d", len(parameters)+1)
		}
	}
	if err == nil {
		if utils.MapContainsKey[string](formData, "Address") && len(formData["Address"]) > 0 {
			whereSql += fmt.Sprintf(" AND address LIKE " + paramIndex)
			parameters = append(parameters, "%"+formData["Address"]+"%")
			paramIndex = "$" + fmt.Sprintf("%d", len(parameters)+1)
		}
	}
	if err == nil {
		if utils.MapContainsKey[string](formData, "Remark") && len(formData["Remark"]) > 0 {
			whereSql += fmt.Sprintf(" AND remark LIKE " + paramIndex)
			parameters = append(parameters, "%"+formData["Remark"]+"%")
			paramIndex = "$" + fmt.Sprintf("%d", len(parameters)+1)
		}
	}

	orderSql := "ORDER BY "
	sortField := "x"
	sortDir := "asc"
	sortI := 0
	if slices.Contains(keys, "sort_field_0") {
		for sortField != "" {
			sortField = formData[fmt.Sprintf("sort_field_%d", sortI)]
			sortDir = formData[fmt.Sprintf("sort_dir_%d", sortI)]
			if sortField == "" {
				break
			}
			if sortI > 0 {
				orderSql += ", "
			}
			orderSql += fmt.Sprintf(" %s %s ", sortField, sortDir)
			sortI += 1
		}
	} else {
		orderSql += " last_order_date desc "
	}

	limitSql := "LIMIT %d OFFSET %d"
	size, _ := strconv.Atoi(formData["size"])
	page, _ := strconv.Atoi(formData["page"])
	limitSql = fmt.Sprintf(limitSql, size, (page-1)*size)

	//orderBySql := ""
	// ord := map[string]string{
	// 	"0": "0",
	// 	"1": "1",
	// 	"2": "2",
	// }
	// if err == nil {
	// 	if utils.MapContainsKey[string](formData, "UnitPriceOrder") && formData["UnitPriceOrder"] != "0" && utils.MapContainsKey[string](ord, formData["UnitPriceOrder"]) {
	// 		if formData["UnitPriceOrder"] == "1" {
	// 			orderBySql += "unit_price ASC, "
	// 		} else {
	// 			orderBySql += "unit_price DESC, "
	// 		}
	// 	}
	// }
	// if err == nil {
	// 	if utils.MapContainsKey[string](formData, "ReviewOrder") && formData["ReviewOrder"] != "0" && utils.MapContainsKey[string](ord, formData["ReviewOrder"]) {
	// 		if formData["ReviewOrder"] == "1" {
	// 			orderBySql += "review ASC, "
	// 		} else {
	// 			orderBySql += "review DESC, "
	// 		}
	// 	}
	// }

	sql = fmt.Sprintf(sql, whereSql, orderSql, limitSql)
	var dt *pgdbcontext.DataTable
	if err == nil {
		dt, err = service.db.Query(sql, parameters...)
	}
	totalCount := 0
	if err == nil {
		results = toSliceMap(dt)
		if len(results) > 0 {
			i64, _ := results[0]["total_count"].(int64)
			totalCount = int(i64)
		}
	}
	return results, totalCount, err
}

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

	if err == nil {
		results["CuisineImageUrl"], _ = service.getCuisimeImageUrl(cuisineId)
	}

	return results, err
}

func (service *CuisineService) getCuisimeImageUrl(cuisineId string) (string, error) {
	var cuisineImageUrl string
	var err error

	//search from cached urls
	cuisineImageUrl, err = service.getCacheCuisineImageUrl(cuisineId)
	if err == nil && len(cuisineImageUrl) > 0 {
		cuisineImageUrl = cuisineImageUrl + "&raw=1"
	} else {
		//get url from dropbox service
		var dpClient *DropboxClient
		err = nil

		dpClient, err = NewDropboxClient()
		var urls []string
		if err == nil {
			urls, err = dpClient.GetSharedLink("/" + cuisineId + "/CuisineImage.png")
		}

		if err == nil {
			if len(urls) >= 1 {
				cuisineImageUrl = urls[0]
				service.putCacheCuisineImageUrl(cuisineId, cuisineImageUrl) //put into cache
				cuisineImageUrl = cuisineImageUrl + "&raw=1"
			} else {
				err = errors.New("shared link not found")
			}
		}
	}

	return cuisineImageUrl, err
}

func (service *CuisineService) getCacheCuisineImageUrl(cuisineId string) (string, error) {
	var sharedLink string = "/assets/food_dummy_250x250.png"
	var err error
	if cacheCuisineImageUrl == nil {
		cacheCuisineImageUrl = make(map[string]string)
	}
	key := fmt.Sprintf("k_%s", cuisineId)
	if utils.MapContainsKey[string](cacheCuisineImageUrl, key) {
		sharedLink = cacheCuisineImageUrl[key]
	} else {
		err = errors.New("url not found")
	}
	return sharedLink, err
}

func (service *CuisineService) SaveCuisine(cuisine map[string]interface{}, cuisineImage io.Reader) error {
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
		if cuisineImage != nil {
			var dpClient *DropboxClient
			if err == nil {
				dpClient, err = NewDropboxClient()
			}
			var img image.Image
			var ext string
			if err == nil {
				img, ext, err = image.Decode(cuisineImage)
			}
			var tmpImgPath string
			if err == nil {
				tmpImgPath, err = utils.ResizeImage(img, ext, 800)
			}
			var tmpImg *os.File
			if err == nil {
				tmpImg, err = os.Open(tmpImgPath)
			}
			var folderPath string
			if err == nil {
				defer os.Remove(tmpImgPath)
				folderPath = fmt.Sprintf("/%d", cuisine["cuisine_id"])
				_, err = dpClient.UploadFile(folderPath, "CuisineImage.png", tmpImg)
			}
			if err == nil {
				dpClient.CreateSharedLink(folderPath + "/CuisineImage.png")
			}
		}
	}

	if err == nil {
		err = service.db.Commit()
	}

	return err
}

func (service *CuisineService) DeleteCuisine(cuisineId string) error {
	repo, err := service.db.GetRepository("cuisine")

	var id int64
	if err == nil {
		id, err = strconv.ParseInt(cuisineId, 10, 64)
	}

	var dt *pgdbcontext.DataTable
	if err == nil {
		dt, err = repo.Select("cuisine_id=$1", id)
	}

	if err == nil {
		if len(dt.Rows) != 1 {
			err = errors.New("not found")
		}
	}

	if err == nil {
		key := make(map[string]interface{})
		key["cuisine_id"] = id
		err = repo.Delete(key)
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

func ClearCache() {
	cacheCuisineImageUrl = make(map[string]string)
}
