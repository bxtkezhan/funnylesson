package database

import (
	"fmt"
	"log"
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

/*
CREATE TABLE "user" (
	"id"	INTEGER NOT NULL,
	"username"	TEXT NOT NULL,
	"password"	TEXT NOT NULL,
	"email"	TEXT NOT NULL,
	"introduction"	TEXT,
	"image"	TEXT,
    "time" INTEGER NOT NULL,
	"level"	INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("id")
);
*/
type User struct {
    Id int
    Username string
    Password string
    Email string
    Introduction string
    Image string
    Ticket int
    Time int64
    Level int
}

type UserSummary struct {
    Id int
    Username string
    Email string
    Introduction string
    Image string
    Ticket int
    Time int64
    Level int
}

func AddUser(ctx context.Context, user *User) error {
    _, err := DB.ExecContext(ctx,
        `insert into user(username, password, email, time, level) values(?, ?, ?, ?, ?)`,
        user.Username, user.Password, user.Email, user.Time, user.Level)
    return err
}

func DelUser(ctx context.Context, id int) error {
    _, err := DB.ExecContext(ctx,
        `delete from user where id = ?`, id)
    return err
}

func GetUserById(ctx context.Context, id int) (*User, error) {
    rows, err := DB.QueryContext(ctx,
        `select id, username, email, introduction, image, ticket, time, level from user where id = ?`, id)
    if err != nil {
        return nil, err
    }
    user := &User{}
    if rows.Next() {
        err := rows.Scan(
            &user.Id,
            &user.Username,
            &user.Email,
            &user.Introduction,
            &user.Image,
            &user.Ticket,
            &user.Time,
            &user.Level)
        if err != nil {
            return nil, err
        }
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return user, nil
}

func HashString(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}

func ValidEmail(ctx context.Context, email string) bool {
    rows, err := DB.QueryContext(ctx,
        `select count(*) from user where email = ?`, email)
    if err != nil {
        log.Println(err)
        return false
    }
    var count int
    if rows.Next() {
        err := rows.Scan(&count)
        if err != nil {
            log.Println(err)
            return false
        }
    }
    if err := rows.Err(); err != nil {
        log.Println(err)
        return false
    }
    if count != 0 {
        return false
    }
    return true
}

func ValidPassword(ctx context.Context, email, password string) int {
    var userId int
    var passwordHash string
    rows, err := DB.QueryContext(ctx,
        `select id, password from user where email = ?`, email)
    if err != nil {
        log.Println(err)
        return 0
    }
    if rows.Next() {
        err := rows.Scan(&userId, &passwordHash)
        if err != nil {
            log.Println(err)
            return 0
        }
    }
    if err := rows.Err(); err != nil {
        log.Println(err)
        return 0
    }
    if HashString(password) != passwordHash {
        return 0
    }
    return userId
}

func SetUser(ctx context.Context, user *User) error {
    _, err := DB.ExecContext(ctx,
        `update user
        set username = ?, password = ?, email = ?, introduction = ?, image = ?, ticket = ?, time = ?, level = ?
        where id = ?`,
        user.Username, user.Password, user.Email, user.Introduction, user.Image, user.Ticket, user.Time, user.Level, user.Id)
    return err
}

func SetUserTicket(ctx context.Context, user *User) error {
    _, err := DB.ExecContext(ctx, `update user set ticket = ?  where id = ?`, user.Ticket, user.Id)
    return err
}

func GetUsers(
    ctx context.Context,
    limit int, offset int,where string,
    args ...interface{}) ([]*UserSummary, error) {

    users := []*UserSummary{}
    query := `select id, username, email, introduction, image, ticket, time, level from user`
    if where != "" {
        query = fmt.Sprintf("%s where %s", query, where)
    }
    query = fmt.Sprintf("%s order by id desc", query)
    if limit != 0 {
        query = fmt.Sprintf("%s limit %d offset %d", query, limit, offset)
    }
    rows, err := DB.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    for rows.Next() {
        user := &UserSummary{}
        err = rows.Scan(
            &user.Id,
            &user.Username,
            &user.Email,
            &user.Introduction,
            &user.Image,
            &user.Ticket,
            &user.Time,
            &user.Level)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    err = rows.Err()
    if err != nil {
        return nil, err
    }
    return users, nil
}

func GetUsersByPage(ctx context.Context, page int, size int) ([]*UserSummary, error) {

    limit := size
    offset := page * size
    return GetUsers(ctx, limit, offset, "")
}

func GetUsersTotal(ctx context.Context) (int, error) {
    var rows *sql.Rows
    var err error
    rows, err = DB.QueryContext(ctx, `select count(*) from user`)
    if err != nil {
        return 0, err
    }
    var count int
    if rows.Next() {
        if err = rows.Scan(&count); err != nil {
            return 0, err
        }
    }
    if err = rows.Err(); err != nil {
        return 0, err
    }
    return count, nil
}

func GetUsersPageTotal(ctx context.Context, size int) (int, error) {
    count, err := GetUsersTotal(ctx)
    if err != nil {
        return 0, err
    }
    total := count / size
    if count % size != 0 {
        total += 1
    }
    return total, nil
}
