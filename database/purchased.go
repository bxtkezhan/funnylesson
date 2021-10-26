package database

import (
	"context"
	"database/sql"
)

func BuyLesson(ctx context.Context, user, lesson int) error {
    _, err := DB.ExecContext(ctx,
        `insert into purchased(user, lesson) values(?, ?)`,
        user, lesson)
    return err
}

func IsBought(ctx context.Context, user, lesson int) (bool, error) {
    var rows *sql.Rows
    var err error
    rows, err = DB.QueryContext(ctx,
        `select exists(select 1 from purchased where user = ? and lesson = ?)`, user, lesson)
    if err != nil {
        return false, err
    }
    var exists bool
    if rows.Next() {
        if err = rows.Scan(&exists); err != nil {
            return false, err
        }
    }
    if err = rows.Err(); err != nil {
        return false, err
    }
    return exists, nil
}

func GetBoughts(ctx context.Context, user int) ([]*Lesson, error) {
    var rows *sql.Rows
    var err error
    lessons := []*Lesson{}
    rows, err = DB.QueryContext(ctx,
        `select lesson from purchased where user = ? order by id desc`, user)
    if err != nil {
        return lessons, err
    }
    var lessonId int
    for rows.Next() {
        err = rows.Scan(&lessonId)
        if err != nil {
            return lessons, err
        }
        lesson, err := GetLessonById(ctx, lessonId)
        if err != nil {
            return lessons, err
        }
        lessons = append(lessons, lesson)
    }
    if err = rows.Err(); err != nil {
        return lessons, err
    }
    return lessons, nil
}
