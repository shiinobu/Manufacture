package user_model

import (
	"golang.org/x/crypto/bcrypt"

	DB "id.benderaku.manufacture/app/helpers"
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
    rows, err := DB.QueryParams(query, []any{id})
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
    rows, err := DB.QueryParams(query, []any{username})
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

func GetAllUsers() ([]User, error) {
    query := "SELECT id, fUsername, fPassword, fEmail FROM tusers WHERE fDelete = 0"
    rows, err := DB.Query(query)
    if err != nil {
        return nil, err
    }

    users := []User{}
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}

func CreateUser(user *User) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    query := "INSERT INTO tusers (fUsername, fPassword, fEmail) VALUES (?, ?, ?)"
    _, err = DB.QueryExec(query, []any{user.Username, string(hashedPassword), user.Email})
    return err
}

func UpdateUser(id int, user *User) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    query := "UPDATE tusers SET fUsername = ?, fPassword = ?, fEmail = ? WHERE id = ? AND fDelete = 0"
    _, err = DB.QueryExec(query, []any{user.Username, string(hashedPassword), user.Email, id})
    return err
}

func UpdateStatusUser(id int, user *User) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    query := "UPDATE tusers SET fUsername = ?, fPassword = ?, fEmail = ?, fDelete = 0 WHERE id = ?"
    _, err = DB.QueryExec(query, []any{user.Username, string(hashedPassword), user.Email, id})
    return err
}

func DeleteUser(id int) error {
    query := "UPDATE tusers SET fDelete = 1 WHERE id = ?"
    _, err := DB.QueryExec(query, []any{id})
    return err
}