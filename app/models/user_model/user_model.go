package user_model

import (
	"math"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	EXE "id.benderaku.manufacture/app/helpers"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     int    `json:"role"`
	Islogin  int    `json:"islogin"`
	Status   int    `json:"-"`
}

func GetListUsers(writ http.ResponseWriter, req *http.Request) ([]User, int, error) {
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
	rows, err := EXE.QueryParams("SELECT id, fUsername, fPassword, fEmail, fFullname, fRole, fIsLogin FROM tusers WHERE fDelete = 0 LIMIT ? OFFSET ?", []any{limit, offset})
	if err != nil {
		return nil, 0, err
	}

	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	var totalCount int
	res, err := EXE.Query("SELECT COUNT(id) FROM tusers WHERE fDelete = 0 LIMIT 1")
	if err != nil {
		return nil, 0, err
	}

	if res.Next() {
		if err := res.Scan(&totalCount); err != nil {
			return nil, 0, err
		}
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
	return users, totalPages, nil
}

func GetUserByID(id int) (*User, error) {
	user := &User{}
	rows, err := EXE.QueryParams("SELECT id, fUsername, fPassword, fEmail, fFullname, fRole, fIsLogin FROM tusers WHERE id = ? AND fDelete = 0 LIMIT 1", []any{id})
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		if err = rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Fullname, &user.Role, &user.Islogin); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	rows, err := EXE.QueryParams("SELECT id, fPassword FROM tusers WHERE fEmail = ? AND fDelete = 0 AND fIsLogin = 1 LIMIT 1", []any{email})
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		if err = rows.Scan(&user.ID, &user.Password); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func CheckStatusUser(email string) (int, error) {
	user := &User{}
	rows, err := EXE.QueryParams("SELECT fDelete FROM tusers WHERE fEmail = ? LIMIT 1", []any{email})
	if err != nil {
		return 0, err
	}
	if rows.Next() {
		if err = rows.Scan(&user.Status); err != nil {
			return 0, err
		}
	}
	return user.Status, nil
}

func CreateUser(user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	status, err := CheckStatusUser(user.Email)
	if err != nil {
		if (status == 1) {
			_, err = EXE.QueryExec("UPDATE tusers SET fUsername = ?, fPassword = ?, fIsLogin = 0, fDelete = 0, fFullname = ? WHERE fEmail = ?", []any{user.Username, string(hashedPassword), user.Fullname, user.Email})
		} else {
			_, err = EXE.QueryExec("INSERT INTO tusers (fEmail, fUsername, fFullname, fPassword, fRole) VALUES (?, ?, ?, ?)", []any{user.Email, user.Username, user.Email, string(hashedPassword), user.Role})
		}
	}
	return err
}

func UpdateUser(id int, user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = EXE.QueryExec("UPDATE tusers SET fUsername = ?, fPassword = ?, fRole = ?, fIsLogin = ? WHERE id = ? AND fDelete = 0", []any{user.Username, string(hashedPassword), user.Role, user.Islogin, id})
	return err
}

func DeleteUser(id int) error {
	_, err := EXE.QueryExec("UPDATE tusers SET fDelete = 1 WHERE id = ?", []any{id})
	return err
}
