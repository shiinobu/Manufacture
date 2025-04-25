package customer_model

import (
	"math"
	"net/http"
	"strconv"

	EXE "id.benderaku.manufacture/app/helpers"
)

type Customer struct {
	ID      int    `json:"id"`
	Code    string `json:"code"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Status  int    `json:"-"`
}

func GetListCustomer(writ http.ResponseWriter, req *http.Request) ([]Customer, int, error) {
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
	rows, err := EXE.QueryParams("SELECT id, fKode, fNama, fAlamat FROM tcustomer WHERE fDelete = 0 LIMIT ? OFFSET ?", []any{limit, offset})
	if err != nil {
		return nil, 0, err
	}

	customers := []Customer{}
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.Code, &customer.Name, &customer.Address); err != nil {
			return nil, 0, err
		}
		customers = append(customers, customer)
	}
	var totalCount int
	res, err := EXE.Query("SELECT COUNT(id) FROM tcustomer WHERE fDelete = 0 LIMIT 1")
	if err != nil {
		return nil, 0, err
	}

	if res.Next() {
		if err := res.Scan(&totalCount); err != nil {
			return nil, 0, err
		}
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
	return customers, totalPages, nil
}

func GetCustomerByID(id int) (*Customer, error) {
	rows, err := EXE.QueryParams("SELECT id, fKode, fNama, fAlamat FROM tcustomer WHERE id = ? AND fDelete = 0 LIMIT 1", []any{id})
	customer := &Customer{}
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		if err = rows.Scan(&customer.ID, &customer.Code, &customer.Name, &customer.Address); err != nil {
			return nil, err
		}
	}
	return customer, nil
}

func CheckStatusCustomer(code string) (string, error) {
	customer := &Customer{}
	rows, err := EXE.QueryParams("SELECT fDelete FROM tcustomer WHERE fKode = ? LIMIT 1", []any{code})
	if err != nil {
		return "", err
	}
	if rows.Next() {
		if err = rows.Scan(&customer.Status); err != nil {
			return "", err
		}
		if customer.Status == 0 {
			return "EXIST", nil
		}
		if customer.Status == 1 {
			return "UPDATE", nil
		}
	}
	return "INSERT", nil
}

func CreateCustomer(customer *Customer) (*int, error) {
	var status int
	result, err := EXE.Nullable(customer.Address)
	if err != nil {
		result = customer.Address
	}
	action, err := CheckStatusCustomer(customer.Code)
	if err != nil {
		return nil, err
	}
	if action == "EXIST" {
		status = 0
	} else if action == "INSERT" {
		status = 1
		_, err = EXE.QueryExec("INSERT INTO tcustomer (fKode, fNama, fAlamat) VALUES (?, ?, ?)", []any{customer.Code, customer.Name, result})
	} else {
		status = 2
		_, err = EXE.QueryExec("UPDATE tcustomer SET fNama = ?, fAlamat = ?, fDelete = 0 WHERE fKode = ?", []any{customer.Name, result, customer.Code})
	}
	return &status, err
}

func UpdateCustomer(id int, customer *Customer) error {
	_, err := EXE.QueryExec("UPDATE tcustomer SET fNama = ?, fAlamat = ? WHERE id = ? AND fDelete = 0", []any{customer.Name, customer.Address, id})
	return err
}

func DeleteCustomer(id int) error {
	_, err := EXE.QueryExec("UPDATE tcustomer SET fDelete = 1 WHERE id = ?", []any{id})
	return err
}
