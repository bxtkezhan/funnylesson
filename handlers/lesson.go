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
    Routes["/lesson"] = Lesson
    Routes["/lessons"] = Lessons
    SecretRoutes["/addlesson"] = AddLesson
}

func Lesson(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil {
        id = 1
    }
    ctx_, cancel := context.WithCancel(ctx)
    lesson, err := db.GetLessonById(ctx_, id)
    cancel()
    if err != nil {
        log.Panic(err)
    }
    user, err := Authority(r)
    if err != nil {
        log.Panic(err)
    }
    ctx_, cancel = context.WithCancel(ctx)
    isBought, err := db.IsBought(ctx_, user.Id, lesson.Id)
    cancel()
    if err != nil {
        log.Panic(err)
    }
    frame := struct {*db.Lesson; Status int}{lesson, 0}
    cost := lesson.Ticket
    if isBought || user.Level > LevelVip {
        cost = 0
    }
    if user.Ticket < cost {
        frame.Status = 1
    }
    if user.Level < lesson.Level {
        frame.Status = 2
    }
    if frame.Status != 0 {
        frame.Lesson = &db.Lesson{}
    } else if cost > 0 {
        user.Ticket -= cost
        db.SetUserTicket(ctx, user)
        db.BuyLesson(ctx, user.Id, lesson.Id)
    }
    Jsonify(w, frame)
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
    total, err := db.GetLessonsPageTotal(ctx, size, owner)
    if err != nil {
        log.Panic(err)
    }
    Jsonify(w, struct{
        Lessons []*db.LessonSummary
        Total int
        Page int
        Size int
    } {lessons, total, page, size})
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
        http.Redirect(w, r, "/worker/addlesson.html", http.StatusSeeOther)
        return
    }
    if len(introduction) < 2 {
        log.Println("introduction is too short")
        http.Redirect(w, r, "/worker/addlesson.html", http.StatusSeeOther)
        return
    }
    if len(source) < 1 {
        log.Println("source is empty")
        http.Redirect(w, r, "/worker/addlesson.html", http.StatusSeeOther)
        return
    }
    type_, err := strconv.Atoi(r.PostFormValue("type"))
    if err != nil {
        type_ = 0
    }
    ticket, err := strconv.Atoi(r.PostFormValue("ticket"))
    if err != nil {
        ticket = 0
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
        http.Redirect(w, r, "/worker/addlesson.html", http.StatusSeeOther)
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
        Ticket: ticket,
        Level: level}
    lesson.Id, err = strconv.Atoi(r.PostFormValue("id"))
    if err != nil {
        db.AddLesson(ctx, lesson)
    } else {
        db.SetLesson(ctx, lesson)
    }
    http.Redirect(w, r, "/lessons.html", http.StatusSeeOther)
}
