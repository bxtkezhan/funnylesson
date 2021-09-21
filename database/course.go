package database

import (
	"context"
	"database/sql"
	"fmt"
)

/*
CREATE TABLE "course" (
	"id"	INTEGER NOT NULL,
	"time"	INTEGER NOT NULL,
	"title"	TEXT NOT NULL,
	"introduction"	TEXT NOT NULL,
	"keywords"	TEXT NOT NULL,
	"image"	TEXT NOT NULL,
	"category"	TEXT NOT NULL,
	"total"	INTEGER NOT NULL DEFAULT 0,
	"owner"	INTEGER NOT NULL,
	"level"	INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("id")
);
*/
type Course struct {
    Id int
    Time int64
    Title string
    Introduction string
    Keywords string
    Image string
    Category string
    Owner int
    Total int
    Level int
}

type CourseSummary struct {
    Id int
    Time int64
    Title string
    Introduction string
    Image string
    Level int
}

func AddCourse(ctx context.Context, course *Course) error {
    _, err := DB.ExecContext(ctx,
        `insert into course(time, title, introduction, keywords, image, category, owner, total, level)
        values(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        course.Time, course.Title, course.Introduction, course.Keywords, course.Image, course.Category, course.Owner,
        course.Total, course.Level)
    return err
}

func DelCourse(ctx context.Context, id int) error {
    _, err := DB.ExecContext(ctx,
        `delete from course where id = ?`, id)
    return err
}

func GetCourseById(ctx context.Context, id int) (*Course, error) {
    rows, err := DB.QueryContext(ctx, `select * from course where id = ?`, id)
    if err != nil {
        return nil, err
    }
    course := &Course{}
    if rows.Next() {
        err = rows.Scan(
            &course.Id,
            &course.Time,
            &course.Title,
            &course.Introduction,
            &course.Keywords,
            &course.Image,
            &course.Category,
            &course.Owner,
            &course.Total,
            &course.Level)
        if err != nil {
            return nil, err
        }
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }
    return course, nil
}

func SetCourse(ctx context.Context, course *Course) error {
    if course.Image == "" {
        _, err := DB.ExecContext(ctx,
            `update course
            set time = ?, title = ?, introduction = ?, keywords = ?, category = ?, owner = ?, total = ?, level = ?
            where id = ?`,
            course.Time, course.Title, course.Introduction, course.Keywords, course.Category, course.Owner,
            course.Total, course.Level,
            course.Id)
        return err
    }
    _, err := DB.ExecContext(ctx,
        `update course
        set time = ?, title = ?, introduction = ?, keywords = ?, image = ?, category = ?, owner = ?, total = ?, level = ?
        where id = ?`,
        course.Time, course.Title, course.Introduction, course.Keywords, course.Image, course.Category, course.Owner,
        course.Total, course.Level,
        course.Id)
    return err
}

func GetCourses(
    ctx context.Context,
    limit int, offset int,where string,
    args ...interface{}) ([]*CourseSummary, error) {

    courses := []*CourseSummary{}
    query := `select id, time, title, introduction, image, level from course`
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
        course := &CourseSummary{}
        err = rows.Scan(
            &course.Id,
            &course.Time,
            &course.Title,
            &course.Introduction,
            &course.Image,
            &course.Level)
        if err != nil {
            return nil, err
        }
        courses = append(courses, course)
    }
    err = rows.Err()
    if err != nil {
        return nil, err
    }
    return courses, nil
}

func GetCoursesByPage(
    ctx context.Context,
    page int, size int, category string) ([]*CourseSummary, error) {

    limit := size
    offset := page * size
    if category != "" {
        return GetCourses(ctx, limit, offset, "category = ?", category)
    }
    return GetCourses(ctx, limit, offset, "")
}

func GetCoursesTotal(ctx context.Context, category string) (int, error) {
    var rows *sql.Rows
    var err error
    if category != "" {
        rows, err = DB.QueryContext(ctx, `select count(*) from course where category = ?`, category)
    } else {
        rows, err = DB.QueryContext(ctx, `select count(*) from course`)
    }
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

func GetCoursesPageTotal(ctx context.Context, size int, category string) (int, error) {
    count, err := GetCoursesTotal(ctx, category)
    if err != nil {
        return 0, err
    }
    total := count / size
    if count % size != 0 {
        total += 1
    }
    return total, nil
}
