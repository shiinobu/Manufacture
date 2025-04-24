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
	Password string `json:"password"`
	Email    string `json:"email"`
	Status   int    `json:"-"`
}

func GetUserByID(id int) (*User, error) {
	user := &User{}
	query := "SELECT id, fUsername, fEmail FROM tusers WHERE id = ? AND fDelete = 0"
	rows, err := EXE.QueryParams(query, []any{id})
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		if err = rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	query := "SELECT id, fUsername, fPassword, fEmail FROM tusers WHERE fUsername = ? AND fDelete = 0"
	rows, err := EXE.QueryParams(query, []any{username})
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		if err = rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func CheckStatusUser(username string, email string) (*User, error) {
	user := &User{}
	query := "SELECT id, fDelete FROM tusers WHERE fUsername = ? AND fEmail = ?"
	rows, err := EXE.QueryParams(query, []any{username, email})
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		if err = rows.Scan(&user.ID, &user.Status); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func GetAllUsers(writ http.ResponseWriter, req *http.Request) ([]User, int, error) {
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
	query := "SELECT id, fUsername, fPassword, fEmail FROM tusers WHERE fDelete = 0 LIMIT ? OFFSET ?"
	rows, err := EXE.QueryParams(query, []any{limit, offset})
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
	res, err := EXE.Query("SELECT COUNT(id) FROM tusers WHERE fDelete = 0")
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

func CreateUser(user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := "INSERT INTO tusers (fUsername, fPassword, fEmail) VALUES (?, ?, ?)"
	_, err = EXE.QueryExec(query, []any{user.Username, string(hashedPassword), user.Email})
	return err
}

func UpdateUser(id int, user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := "UPDATE tusers SET fUsername = ?, fPassword = ?, fEmail = ? WHERE id = ? AND fDelete = 0"
	_, err = EXE.QueryExec(query, []any{user.Username, string(hashedPassword), user.Email, id})
	return err
}

func UpdateStatusUser(id int, user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := "UPDATE tusers SET fUsername = ?, fPassword = ?, fEmail = ?, fDelete = 0 WHERE id = ?"
	_, err = EXE.QueryExec(query, []any{user.Username, string(hashedPassword), user.Email, id})
	return err
}

func DeleteUser(id int) error {
	query := "UPDATE tusers SET fDelete = 1 WHERE id = ?"
	_, err := EXE.QueryExec(query, []any{id})
	return err
}
