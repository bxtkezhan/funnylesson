package handlers

import (
    "fmt"
    "log"
    "time"
    "context"
    "net/http"
    "path/filepath"
    "os"
    "io"
    "strconv"

    db "../database"

	"github.com/gorilla/sessions"
    "github.com/google/uuid"
)

func init() {
    Routes["/course"] = Course
    Routes["/courses"] = Courses
    Routes["/contents"] = Contents
    SecretRoutes["/addcourse"] = AddCourse
    SecretRoutes["/addcontent"] = AddContent
}

func Course(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil {
        id = 1
    }
    course, err := db.GetCourseById(ctx, id)
    if err != nil {
        log.Panic(err)
    }
    level, _ := Authority(r)
    if level < course.Level {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    Jsonify(w, course)
}

func Courses(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil {
        page = 0
    }
    size, err := strconv.Atoi(r.URL.Query().Get("size"))
    if err != nil {
        size = 1
    }
    category := r.URL.Query().Get("category")
    courses, err := db.GetCoursesByPage(ctx, page, size, category)
    if err != nil {
        log.Panic(err)
    }
    total, err := db.GetCoursesPageTotal(ctx, size, category)
    if err != nil {
        log.Panic(err)
    }
    Jsonify(w, struct{
        Courses []*db.CourseSummary
        Total int
        Page int
        Size int
    } {courses, total, page, size})
}

func AddCourse(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
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
    category := r.PostFormValue("category")
    if len(title) < 1 {
        log.Println("title is empty")
        http.Redirect(w, r, "/addcourse.html", http.StatusSeeOther)
        return
    }
    if len(introduction) < 2 {
        log.Println("introduction is too short")
        http.Redirect(w, r, "/addcourse.html", http.StatusSeeOther)
        return
    }
    if len(category) < 1 {
        log.Println("category is empty")
        http.Redirect(w, r, "/addcourse.html", http.StatusSeeOther)
        return
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
        http.Redirect(w, r, "/addcourse.html", http.StatusSeeOther)
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
        log.Panic(err == http.ErrMissingFile)
    }
    course := &db.Course{
        Time: time.Now().Unix(),
        Title: title,
        Introduction: introduction,
        Keywords: keywords,
        Image: filename,
        Category: category,
        Owner: id,
        Total: 0,
        Level: level}
    course.Id, err = strconv.Atoi(r.PostFormValue("id"))
    if err != nil {
        db.AddCourse(ctx, course)
    } else {
        db.SetCourse(ctx, course)
    }
    http.Redirect(w, r, "/addcourse.html", http.StatusSeeOther)
}

func Contents(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil {
        id = 1
    }
    lessons, err := db.GetContents(ctx, id)
    if err != nil {
        log.Panic(err)
    }
    Jsonify(w, lessons)
}

func AddContent(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
    ctx := r.Context()
    id := s.Values["userid"].(int)
    ctx_, cancel := context.WithCancel(ctx)
    user, err := db.GetUserById(ctx_, id)
    defer cancel()
    if err != nil {
        log.Panic(err)
    }
    if user.Level < LevelWorker {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    course, err := strconv.Atoi(r.PostFormValue("course"))
    if err != nil {
        http.Error(w, "Type of course is not int", http.StatusBadRequest)
        return
    }
    lesson, err := strconv.Atoi(r.PostFormValue("lesson"))
    if err != nil {
        http.Error(w, "Type of lesson is not int", http.StatusBadRequest)
        return
    }
    index, err := strconv.Atoi(r.PostFormValue("index"))
    if err != nil {
        index = -1
    }
    isExist, err := db.IsExistContent(ctx_, course, lesson)
    if err != nil {
        log.Panic(err)
    }
    if isExist {
        http.Error(w, "Lesson is exist", http.StatusBadRequest)
        return
    }
    if index < 0 {
        maxIndex, err := db.GetMaxIndex(ctx_, course)
        if err != nil {
            log.Panic(err)
        }
        index = maxIndex + 1
    }
    cancel()
    db.AddContent(ctx, course, lesson, index)
    http.Redirect(w, r,
        fmt.Sprintf("/course.html?id=%d", course),
        http.StatusSeeOther)
}
