package handlers

import (
    "context"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "time"

    db "../database"

    "github.com/google/uuid"
    "github.com/gorilla/sessions"
)

func init() {
    Routes["/lessons"] = Lessons
    SecretRoutes["/lesson"] = Lesson
    SecretRoutes["/addlesson"] = AddLesson
}

func Lesson(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
    ctx := r.Context()
    id := s.Values["userid"].(int)
    ctx_, cancel := context.WithCancel(ctx)
    user, err := db.GetUserById(ctx_, id)
    cancel()
    if err != nil {
        log.Panic(err)
    }
    lessonId, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil {
        lessonId = 1
    }
    lesson, err := db.GetLessonById(ctx, lessonId)
    if user.Level < lesson.Level {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    if err != nil {
        log.Panic(err)
    }
    Jsonify(w, lesson)
}

func Lessons(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil {
        page = 0
    }
    size, err := strconv.Atoi(r.URL.Query().Get("size"))
    if err != nil {
        size = 1
    }
    owner, err := strconv.Atoi(r.URL.Query().Get("owner"))
    if err != nil {
        owner = 0
    }
    lessons, err := db.GetLessonsByPage(ctx, page, size, owner)
    if err != nil {
        log.Panic(err)
    }
    Jsonify(w, lessons)
}

func AddLesson(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
    ctx := r.Context()
    id := s.Values["userid"].(int)
    ctx_, cancel := context.WithCancel(ctx)
    user, err := db.GetUserById(ctx_, id)
    cancel()
    if err != nil {
        log.Panic(err)
    }
    if user.Level < LevelWorker {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    title := r.PostFormValue("title")
    introduction := r.PostFormValue("introduction")
    keywords := r.PostFormValue("keywords")
    source := r.PostFormValue("source")
    if len(title) < 1 {
        log.Println("title is empty")
        http.Redirect(w, r, "/addlesson.html", http.StatusSeeOther)
        return
    }
    if len(introduction) < 2 {
        log.Println("introduction is too short")
        http.Redirect(w, r, "/addlesson.html", http.StatusSeeOther)
        return
    }
    if len(source) < 1 {
        log.Println("source is empty")
        http.Redirect(w, r, "/addlesson.html", http.StatusSeeOther)
        return
    }
    type_, err := strconv.Atoi(r.PostFormValue("type"))
    if err != nil {
        type_ = 0
    }
    level := LevelDefault
    switch r.PostFormValue("level") {
    case "", "default":
        level = LevelDefault
    case "vip":
        level = LevelVip
    case "worker":
        level = LevelWorker
    case "admin":
        level = LevelAdmin
    default:
        log.Println("no support level:", r.PostFormValue("level"))
        http.Redirect(w, r, "/addlesson.html", http.StatusSeeOther)
    }
    filename := ""
    srcFile, header, err := r.FormFile("image")
    if err == nil {
        defer srcFile.Close()
        filename = uuid.New().String() + filepath.Ext(header.Filename)
        dstFile, err := os.Create(filepath.Join("./img/", filename))
        if err != nil {
            log.Panic(err)
        }
        defer dstFile.Close()
        io.Copy(dstFile, srcFile)
        log.Println(header.Filename)
        log.Println(header.Size)
    } else if err != http.ErrMissingFile {
        log.Panic(err)
    }
    lesson := &db.Lesson{
        Time: time.Now().Unix(),
        Title: title,
        Introduction: introduction,
        Keywords: keywords,
        Image: filename,
        Source: source,
        Type: type_,
        Owner: id,
        Level: level}
    lesson.Id, err = strconv.Atoi(r.PostFormValue("id"))
    if err != nil {
        db.AddLesson(ctx, lesson)
    } else {
        db.SetLesson(ctx, lesson)
    }
    http.Redirect(w, r, "/lessons.html", http.StatusSeeOther)
}
