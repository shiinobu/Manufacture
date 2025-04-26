package purchase_model

import (
	"math"
	"net/http"
	"strconv"

	EXE "id.benderaku.manufacture/app/helpers"
)

type Product struct {
	SKU           string  `json:"code"`
	PriceBuy      float64 `json:"pricebuy"`
	PriceSell     float64 `json:"pricesell"`
	Quantity      int     `json:"quantity"`
	TotalBuy      float64 `json:"totalbuy"`
	TotalSell     float64 `json:"totalsell"`
	TotalDiscount float64 `json:"totaldiscount"`
	Note          string  `json:"note"`
}

type Purchase struct {
	ID            int       `json:"id"`
	Div           string    `json:"-"`
	OrderNumber   string    `json:"-"`
	Reff          string    `json:"noreff"`
	Date          string    `json:"date"`
	AuthId        int       `json:"auth_id"`
	CustomerId    int       `json:"customer_id"`
	SupplierId    int       `json:"supplier_id"`
	PaymentType   int       `json:"payment_type"`
	TotalBuy      float64   `json:"totalbuy"`
	TotalSell     float64   `json:"totalsell"`
	TotalGross    float64   `json:"totalgross"`
	TotalDiscount float64   `json:"totaldiscount"`
	TotalPpn      float64   `json:"totalppn"`
	TotalPayment  float64   `json:"totalpayment"`
	PpnType       int       `json:"ppn_type"`
	Note          string    `json:"note"`
	Products      []Product `json:"products"`
	Posting       int       `json:"posting"`
}

func GetListPurchase(writ http.ResponseWriter, req *http.Request) ([]Purchase, int, error) {
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
	rows, err := EXE.QueryParams("SELECT id, fOrderNumber, fReffNumber, fTanggal FROM tpurchaseoders WHERE fDelete = 0 LIMIT ? OFFSET ?", []any{limit, offset})
	if err != nil {
		return nil, 0, err
	}

	purchases := []Purchase{}
	for rows.Next() {
		var purchase Purchase
		if err := rows.Scan(&purchase.ID, &purchase.OrderNumber, &purchase.Reff, &purchase.Date); err != nil {
			return nil, 0, err
		}
		purchases = append(purchases, purchase)
	}
	var totalCount int
	res, err := EXE.Query("SELECT COUNT(id) FROM tpurchaseoders WHERE fDelete = 0 LIMIT 1")
	if err != nil {
		return nil, 0, err
	}

	if res.Next() {
		if err := res.Scan(&totalCount); err != nil {
			return nil, 0, err
		}
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
	return purchases, totalPages, nil
}

func GetPurchaseByID(id int) (*Purchase, error) {
	rows, err := EXE.QueryParams("SELECT id, fOrderNumber, fReffNumber, fTanggal FROM tpurchaseoders WHERE id = ? AND fDelete = 0 LIMIT 1", []any{id})
	purchase := &Purchase{}
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		if err = rows.Scan(&purchase.ID, &purchase.OrderNumber, &purchase.Reff, &purchase.Date); err != nil {
			return nil, err
		}
	}
	return purchase, nil
}

func CreatePurchase(purchase *Purchase) (*int, error) {
	status := 0
	_, err := EXE.QueryExec("INSERT INTO tpurchaseoders (fOrderNumber, fReffNumber, fTanggal) VALUES (?, ?, ?)", []any{purchase.OrderNumber, purchase.Reff, purchase.Date})
	return &status, err
}
/*
func UpdatePurchase(id int, purchase *Purchase) error {
	_, err := EXE.QueryExec("UPDATE tpurchaseoders SET fNama = ?, fAlamat = ? WHERE id = ? AND fDelete = 0", []any{purchase.Name, purchase.Address, id})
	return err
}

func DeletePurchase(id int) error {
	_, err := EXE.QueryExec("UPDATE tpurchaseoders SET fDelete = 1 WHERE id = ?", []any{id})
	return err
}
*/
