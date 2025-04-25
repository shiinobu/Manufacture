package product_model

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	EXE "id.benderaku.manufacture/app/helpers"
)

type Material struct {
	SKU      string `json:"code"`
	Quantity int    `json:"quantity"`
}

type Product struct {
	ID      int        `json:"id"`
	SKU     string     `json:"code"`
	Name    string     `json:"name"`
	Type    int        `json:"type"`
	Buy     float64    `json:"buy"`
	Sell    float64    `json:"sell"`
	Quality string     `json:"quality"`
	Unit    string     `json:"unit"`
	Stock   int        `json:"stock"`
	Status  int        `json:"-"`
	Items   []Material `json:"items"`
}

func GetListProducts(writ http.ResponseWriter, req *http.Request) ([]Product, int, error) {
	param := req.URL.Query()
	pageStr := param.Get("page")
	limitStr := param.Get("limit")
	tipeStr := param.Get("type")
	page := 1
	limit := 10
	tipe := 0

	if pageStr != "" {
		pageInt, err := strconv.Atoi(pageStr)
		if err != nil || pageInt < 1 {
			EXE.ERROR.Println("Invalid page parameter:", pageStr)
			EXE.SendResponse(writ, "", http.StatusBadRequest, "Invalid page parameter", "")
			return nil, 0, err
		}
		page = pageInt
	}

	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err != nil || limitInt < 1 {
			EXE.ERROR.Println("Invalid limit parameter:", limitStr)
			EXE.SendResponse(writ, "", http.StatusBadRequest, "Invalid limit parameter", "")
			return nil, 0, err
		}
		limit = limitInt
	}

	if tipeStr != "" {
		tipeInt, err := strconv.Atoi(tipeStr)
		if err != nil || tipeInt < 0 {
			EXE.ERROR.Println("Invalid tipe parameter:", tipeStr)
			EXE.SendResponse(writ, "", http.StatusBadRequest, "Invalid tipe parameter", "")
			return nil, 0, err
		}
		tipe = tipeInt
	}

	offset := (page - 1) * limit
	rows, err := EXE.QueryParams("SELECT id, fKodeBrg, fNamaBrg, fSatuan, fType, fBeli, fJual, fQuality, fStock, fMaterial FROM tproducts WHERE fType = ? AND fDelete = 0 LIMIT ? OFFSET ?", []any{tipe, limit, offset})
	if err != nil {
		return nil, 0, err
	}

	products := []Product{}
	for rows.Next() {
		var product Product
		var materialsJSON []byte
		if err := rows.Scan(&product.ID, &product.SKU, &product.Name, &product.Unit, &product.Type, &product.Buy, &product.Sell, &product.Quality, &product.Stock, &materialsJSON); err != nil {
			return nil, 0, err
		}
		if len(materialsJSON) > 0 {
			var materials []Material
			if err := json.Unmarshal(materialsJSON, &materials); err != nil {
				return nil, 0, err
			}
			product.Items = materials
		} else {
			product.Items = nil
		}
		products = append(products, product)
	}

	var totalCount int
	res, err := EXE.QueryParams("SELECT COUNT(id) FROM tproducts WHERE fType = ? AND fDelete = 0 LIMIT 1", []any{tipe})
	if err != nil {
		return nil, 0, err
	}

	if res.Next() {
		if err := res.Scan(&totalCount); err != nil {
			return nil, 0, err
		}
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
	return products, totalPages, nil
}

func GetProductMaterials() ([]Product, error) {
	rows, err := EXE.Query("SELECT id, fKodeBrg, fNamaBrg, fSatuan, fType, fBeli, fJual, fQuality, fStock FROM tproducts WHERE (fType = 0 OR fType = 1) AND fDelete = 0")
	if err != nil {
		return nil, err
	}

	products := []Product{}
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.SKU, &product.Name, &product.Unit, &product.Type, &product.Buy, &product.Sell, &product.Quality, &product.Stock); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func GetProductByID(id int) (*Product, error) {
	product := &Product{}
	rows, err := EXE.QueryParams("SELECT id, fKodeBrg, fNamaBrg, fSatuan, fType, fBeli, fJual, fQuality, fStock, fMaterial FROM tproducts WHERE id = ? AND fDelete = 0 LIMIT 1", []any{id})
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		if err = rows.Scan(&product.ID, &product.SKU, &product.Name, &product.Unit, &product.Type, &product.Buy, &product.Sell, &product.Quality, &product.Stock, &product.Items); err != nil {
			return nil, err
		}
	}
	return product, nil
}

func CheckStatusProduct(code string) (int, error) {
	product := &Product{}
	rows, err := EXE.QueryParams("SELECT fDelete FROM tproducts WHERE fKodeBrg = ? LIMIT 1", []any{code})
	if err != nil {
		return 0, err
	}
	if rows.Next() {
		if err = rows.Scan(&product.Status); err != nil {
			return 0, err
		}
	}
	return product.Status, nil
}

func CreateProduct(product *Product) error {
	var material []byte
	if len(product.Items) > 0 {
		json, err := json.Marshal(product.Items)
		if err != nil {
			return err
		}
		material = json
	} else {
		material = nil
	}
	convert, err := EXE.Nullable(material)
	if err != nil {
		return err
	}
	status, err := CheckStatusProduct(product.SKU)
	if err != nil {
		if (status == 1) {
			_, err = EXE.QueryExec("UPDATE tproducts SET fNamaBrg = ?, fSatuan = ?, fType = ?, fBeli = ?, fJual = ?, fQuality = ?, fStock = ?, fDelete = 0 WHERE fKodeBrg = ?", []any{product.Name, product.Unit, product.Type, product.Buy, product.Sell, product.Quality, product.Stock, product.SKU})
		} else {
			_, err = EXE.QueryExec("INSERT INTO tproducts (fKodeBrg, fNamaBrg, fSatuan, fType, fBeli, fJual, fQuality, fStock, fMaterial) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", []any{product.SKU, product.Name, product.Unit, product.Type, product.Buy, product.Sell, product.Quality, product.Stock, convert})
		}
	}
	return err
}

func UpdateProduct(id int, product *Product) error {
	_, err := EXE.QueryExec("UPDATE tproducts SET fNamaBrg = ?, fSatuan = ?, fType = ?, fBeli = ?, fJual = ?, fQuality = ?, fStock = ? WHERE id = ?", []any{product.Name, product.Unit, product.Type, product.Buy, product.Sell, product.Quality, product.Stock, id})
	return err
}

func DeleteProduct(id int) error {
	_, err := EXE.QueryExec("UPDATE tproducts SET fDelete = 1 WHERE id = ? AND fDelete = 0", []any{id})
	return err
}
