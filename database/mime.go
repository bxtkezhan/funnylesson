package database

import (
	"context"
	"database/sql"

    "github.com/google/uuid"
)

func Mining(ctx context.Context, user int, code string) (int, error) {
    var rows *sql.Rows
    var err error
    ctx_, cancel := context.WithCancel(ctx)
    defer cancel()
    rows, err = DB.QueryContext(ctx_,
        `select id, "count" from mime where code = ?`, code)
    if err != nil {
        return 0, err
    }
    var id, count int
    if rows.Next() {
        if err = rows.Scan(&id, &count); err != nil {
            return 0, err
        }
    }
    if err = rows.Err(); err != nil {
        return 0, err
    }
    cancel()
    _, err = DB.ExecContext(ctx, `delete from mime where id = ?`, id)
    if err != nil {
        return 0, err
    }
    return count, nil
}

func GenMime(ctx context.Context, total, count int) error {
    for i := 0; i < total; i++ {
        code := uuid.New().String()
        _, err := DB.ExecContext(ctx,
            `insert into mime(code, count) values(?, ?)`, code, count)
        if err != nil {
            return err
        }
    }
    return nil
}

func GetMimes(ctx context.Context, total, count int) ([]string, error) {
    codes := []string{}
    if total < 1 {
        total = 100
    }
    rows, err := DB.QueryContext(ctx,
        `select code from mime where count = ? limit ?`, count, total)
    if err != nil {
        return nil, err
    }
    for rows.Next() {
        var code string
        err = rows.Scan(&code)
        if err != nil {
            return nil, err
        }
        codes = append(codes, code)
    }
    err = rows.Err()
    if err != nil {
        return nil, err
    }
    return codes, nil
}
