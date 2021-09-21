package database

import (
	"context"
	"database/sql"
)

func AddContent(ctx context.Context, course, lesson, index int) error {
    _, err := DB.ExecContext(ctx,
        `insert into contents(course, lesson, "index") values(?, ?, ?)`,
        course, lesson, index)
    return err
}

func DelContent(ctx context.Context, course, lesson int) error {
    _, err := DB.ExecContext(ctx,
        `delete from contents where course = ? and lesson = ?`,
        course, lesson)
    return err
}

func GetMaxIndex(ctx context.Context, course int) (int, error) {
    var rows *sql.Rows
    var err error
    rows, err = DB.QueryContext(ctx,
        `select max("index") from contents where course = ?`, course)
    if err != nil {
        return 0, err
    }
    var index int
    if rows.Next() {
        if err = rows.Scan(&index); err != nil {
            return 0, err
        }
    }
    if err = rows.Err(); err != nil {
        return 0, err
    }
    return index, nil
}

func IsExistContent(ctx context.Context, course, lesson int) (bool, error) {
    var rows *sql.Rows
    var err error
    rows, err = DB.QueryContext(ctx,
        `select count(*) from contents where course = ? and lesson = ?`,
        course, lesson)
    if err != nil {
        return false, err
    }
    var count int
    if rows.Next() {
        if err = rows.Scan(&count); err != nil {
            return false, err
        }
    }
    if err = rows.Err(); err != nil {
        return false, err
    }
    return count > 0, nil
}

func GetContents(ctx context.Context, course int) ([]*Lesson, error) {
    var rows *sql.Rows
    var err error
    lessons := []*Lesson{}
    rows, err = DB.QueryContext(ctx,
        `select lesson from contents where course = ? order by "index" asc`,
        course)
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
