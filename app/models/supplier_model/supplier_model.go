package supplier_model

import (
	"math"
	"net/http"
	"strconv"

	EXE "id.benderaku.manufacture/app/helpers"
)

type Supplier struct {
    ID      int    `json:"id"`
    Code    string `json:"code"`
    Name    string `json:"name"`
    Address string `json:"address"`
    Status  int    `json:"-"`
}

func GetListSupplier(writ http.ResponseWriter, req *http.Request) ([]Supplier, int, error) {
    param := req.URL.Query()
	pageStr := param.Get("page")
	limitStr := param.Get("limit")
	page := 1
	limit := 10

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

	offset := (page - 1) * limit
    rows, err := EXE.QueryParams("SELECT id, fKode, fNama, fAlamat FROM tsupplier WHERE fDelete = 0 LIMIT ? OFFSET ?", []any{limit, offset})
    if err != nil {
        return nil, 0, err
    }

    suppliers := []Supplier{}
    for rows.Next() {
		var supplier Supplier
		if err := rows.Scan(&supplier.ID, &supplier.Code, &supplier.Name, &supplier.Address); err != nil {
			return nil, 0, err
		}
		suppliers = append(suppliers, supplier)
	}
    var totalCount int
	res, err := EXE.Query("SELECT COUNT(id) FROM tsupplier WHERE fDelete = 0 LIMIT 1")
	if err != nil {
		return nil, 0, err
	}

	if res.Next() {
		if err := res.Scan(&totalCount); err != nil {
			return nil, 0, err
		}
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
	return suppliers, totalPages, nil
}

func GetSupplierByID(id int) (*Supplier, error) {
    rows, err := EXE.QueryParams("SELECT id, fKode, fNama, fAlamat FROM tsupplier WHERE id = ? AND fDelete = 0 LIMIT 1", []any{id})
    supplier := &Supplier{}
    if err != nil {
        return nil, err
    }
    if rows.Next() {
        if err = rows.Scan(&supplier.ID, &supplier.Code, &supplier.Name, &supplier.Address); err != nil {
            return nil, err
        }
    }
    return supplier, nil
}

func CheckStatusSupplier(code string) (int, error) {
    supplier := &Supplier{}
    rows, err := EXE.QueryParams("SELECT fDelete FROM tsupplier WHERE fKode = ? LIMIT 1", []any{code})
    if err != nil {
        return 0, err
    }
    if rows.Next() {
        if err = rows.Scan(&supplier.Status); err != nil {
            return 0, err
        }
    }
    return supplier.Status, nil
}

func CreateSupplier(supplier *Supplier) error {
	status, err := CheckStatusSupplier(supplier.Code)
	if err != nil {
		if (status == 1) {
			_, err = EXE.QueryExec("UPDATE tsupplier SET fNama = ?, fAlamat = ?, fDelete = 0 WHERE fKode = ?", []any{supplier.Name, supplier.Address, supplier.Code})
		} else {
			_, err = EXE.QueryExec("INSERT INTO tsupplier (fKode, fNama, fAlamat) VALUES (?, ?, ?)", []any{supplier.Code, supplier.Name, supplier.Address})
		}
	}
    return err
}

func UpdateSupplier(id int, supplier *Supplier) error {
    _, err := EXE.QueryExec("UPDATE tsupplier SET fNama = ?, fAlamat = ? WHERE id = ? AND fDelete = 0", []any{supplier.Name, supplier.Address, id})
    return err
}

func UpdateStatusSupplier(id int, supplier *Supplier) error {
    _, err := EXE.QueryExec("UPDATE tsupplier SET fNama = ?, fAlamat = ?, fDelete = 0 WHERE id = ?", []any{supplier.Name, supplier.Address, id})
    return err
}

func DeleteSupplier(id int) error {
    _, err := EXE.QueryExec("UPDATE tsupplier SET fDelete = 1 WHERE id = ?", []any{id})
    return err
}