package supplier_model

import (
	"golang.org/x/crypto/bcrypt"

	DB "id.benderaku.manufacture/app/helpers"
)

type Supplier struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
    Email    string `json:"email"`
    Status int `json:"-"`
}

func GetSupplierByID(id int) (*Supplier, error) {
    query := "SELECT id, fUsername, fEmail FROM tusers WHERE id = ? AND fDelete = 0"
    rows, err := DB.QueryParams(query, []any{id})
    user := &Supplier{}
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

func CheckStatusSupplier(username string, email string) (*Supplier, error) {
    user := &Supplier{}
    query := "SELECT id, fDelete FROM tusers WHERE fUsername = ? AND fEmail = ?"
    rows, err := DB.QueryParams(query, []any{username, email})
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

func GetAllSupplier() ([]Supplier, error) {
    query := "SELECT id, fUsername, fPassword, fEmail FROM tusers WHERE fDelete = 0"
    rows, err := DB.Query(query)
    if err != nil {
        return nil, err
    }

    users := []Supplier{}
    for rows.Next() {
        var user Supplier
        if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}

func CreateSupplier(user *Supplier) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    query := "INSERT INTO tusers (fUsername, fPassword, fEmail) VALUES (?, ?, ?)"
    _, err = DB.QueryExec(query, []any{user.Username, string(hashedPassword), user.Email})
    return err
}

func UpdateSupplier(id int, user *Supplier) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    query := "UPDATE tusers SET fUsername = ?, fPassword = ?, fEmail = ? WHERE id = ? AND fDelete = 0"
    _, err = DB.QueryExec(query, []any{user.Username, string(hashedPassword), user.Email, id})
    return err
}

func UpdateStatusSupplier(id int, user *Supplier) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    query := "UPDATE tusers SET fUsername = ?, fPassword = ?, fEmail = ?, fDelete = 0 WHERE id = ?"
    _, err = DB.QueryExec(query, []any{user.Username, string(hashedPassword), user.Email, id})
    return err
}

func DeleteSupplier(id int) error {
    query := "UPDATE tusers SET fDelete = 1 WHERE id = ?"
    _, err := DB.QueryExec(query, []any{id})
    return err
}