package database

import (
	"context"
	"database/sql"
	"fmt"
)

/*
CREATE TABLE "lesson" (
	"id"	INTEGER NOT NULL,
	"time"	INTEGER NOT NULL,
	"title"	TEXT NOT NULL,
	"introduction"	TEXT NOT NULL,
	"keywords"	TEXT NOT NULL,
	"image"	TEXT NOT NULL,
	"source"	TEXT NOT NULL,
	"type"	INTEGER NOT NULL DEFAULT 0,
	"owner"	INTEGER NOT NULL,
	"level"	INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("id")
);
*/
type Lesson struct {
    Id int
    Time int64
    Title string
    Introduction string
    Keywords string
    Image string
    Source string
    Type int
    Owner int
    Ticket int
    Level int
}

type LessonSummary struct {
    Id int
    Time int64
    Title string
    Introduction string
    Image string
    Ticket int
    Level int
}

func AddLesson(ctx context.Context, lesson *Lesson) error {
    _, err := DB.ExecContext(ctx,
        `insert into lesson(time, title, introduction, keywords, image, source, type, owner, ticket, level)
        values(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        lesson.Time, lesson.Title, lesson.Introduction, lesson.Keywords, lesson.Image, lesson.Source,
        lesson.Type, lesson.Owner, lesson.Ticket, lesson.Level)
    return err
}

func DelLesson(ctx context.Context, id int) error {
    _, err := DB.ExecContext(ctx,
        `delete from lesson where id = ?`, id)
    return err
}

func GetLessonById(ctx context.Context, id int) (*Lesson, error) {
    rows, err := DB.QueryContext(ctx, `select * from lesson where id = ?`, id)
    if err != nil {
        return nil, err
    }
    lesson := &Lesson{}
    if rows.Next() {
        err = rows.Scan(
            &lesson.Id,
            &lesson.Time,
            &lesson.Title,
            &lesson.Introduction,
            &lesson.Keywords,
            &lesson.Image,
            &lesson.Source,
            &lesson.Type,
            &lesson.Owner,
            &lesson.Ticket,
            &lesson.Level)
        if err != nil {
            return nil, err
        }
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }
    return lesson, nil
}

func SetLesson(ctx context.Context, lesson *Lesson) error {
    if lesson.Image == "" {
        _, err := DB.ExecContext(ctx,
            `update lesson
            set time = ?, title = ?, introduction = ?, keywords = ?, source = ?, type = ?, owner = ?, ticket = ?, level = ?
            where id = ?`,
            lesson.Time, lesson.Title, lesson.Introduction, lesson.Keywords, lesson.Source, lesson.Type,
            lesson.Owner, lesson.Ticket, lesson.Level,
            lesson.Id)
        return err
    }
    _, err := DB.ExecContext(ctx,
        `update lesson
        set time = ?, title = ?, introduction = ?, keywords = ?, image = ?, source = ?, type = ?, owner = ?, ticket = ?, level = ?
        where id = ?`,
        lesson.Time, lesson.Title, lesson.Introduction, lesson.Keywords, lesson.Image, lesson.Source, lesson.Type,
        lesson.Owner, lesson.Ticket, lesson.Level,
        lesson.Id)
    return err
}

func GetLessons(
    ctx context.Context,
    limit int, offset int, where string,
    args ...interface{}) ([]*LessonSummary, error) {

    lessons := []*LessonSummary{}
    query := `select id, time, title, introduction, image, ticket, level from lesson`
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
        lesson := &LessonSummary{}
        err = rows.Scan(
            &lesson.Id,
            &lesson.Time,
            &lesson.Title,
            &lesson.Introduction,
            &lesson.Image,
            &lesson.Ticket,
            &lesson.Level)
        if err != nil {
            return nil, err
        }
        lessons = append(lessons, lesson)
    }
    err = rows.Err()
    if err != nil {
        return nil, err
    }
    return lessons, nil
}

func GetLessonsByPage(
    ctx context.Context,
    page int, size int, owner int) ([]*LessonSummary, error) {

    limit := size
    offset := page * size
    if owner != 0 {
        return GetLessons(ctx, limit, offset, "owner = ?", owner)
    }
    return GetLessons(ctx, limit, offset, "")
}

func GetLessonsTotal(ctx context.Context, owner int) (int, error) {
    var rows *sql.Rows
    var err error
    if owner != 0 {
        rows, err = DB.QueryContext(ctx, `select count(*) from lesson where owner = ?`, owner)
    } else {
        rows, err = DB.QueryContext(ctx, `select count(*) from lesson`)
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

func GetLessonsPageTotal(ctx context.Context, size int, owner int) (int, error) {
    count, err := GetLessonsTotal(ctx, owner)
    if err != nil {
        return 0, err
    }
    total := count / size
    if count % size != 0 {
        total += 1
    }
    return total, nil
}
