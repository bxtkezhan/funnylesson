package database

import (
	"context"
	"database/sql"
)

func Follow(ctx context.Context, user, course int) error {
    _, err := DB.ExecContext(ctx,
        `insert into subscribe(user, course) values(?, ?)`,
        user, course)
    return err
}

func UnFollow(ctx context.Context, user, course int) error {
    _, err := DB.ExecContext(ctx,
        `delete from subscribe where user = ? and course = ?`,
        user, course)
    return err
}

func CountFollowers(ctx context.Context, course int) (int, error) {
    var rows *sql.Rows
    var err error
    rows, err = DB.QueryContext(ctx,
        `select count(*) from subscribe where course = ?`, course)
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

func GetLikes(ctx context.Context, user int) ([]*Course, error) {
    var rows *sql.Rows
    var err error
    courses := []*Course{}
    rows, err = DB.QueryContext(ctx,
        `select course from subscribe where user = ?`, user)
    if err != nil {
        return courses, err
    }
    var courseId int
    for rows.Next() {
        err = rows.Scan(&courseId)
        if err != nil {
            return courses, err
        }
        course, err := GetCourseById(ctx, courseId)
        if err != nil {
            return courses, err
        }
        courses = append(courses, course)
    }
    if err = rows.Err(); err != nil {
        return courses, err
    }
    return courses, nil
}
