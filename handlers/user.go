package handlers

import (
	"log"
	"net/http"
	"strconv"
    // "context"

	db "../database"

	"github.com/gorilla/sessions"
)

func init() {
    SecretRoutes["/user"] = User
    SecretRoutes["/users"] = Users
    SecretRoutes["/likes"] = Likes
    SecretRoutes["/inlikes"] = InLikes
    SecretRoutes["/follow"] = Follow
    SecretRoutes["/unfollow"] = UnFollow
    SecretRoutes["/mining"] = Mining
    SecretRoutes["/genmimes"] = GenMimes
    SecretRoutes["/getmimes"] = GetMimes
}

func User(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
    ctx := r.Context()
    id := s.Values["userid"].(int)
    user, err := db.GetUserById(ctx, id)
    if err != nil {
        log.Panic(err)
    }
    _id := r.URL.Query().Get("id")
    if _id != "" && user.Level >= LevelWorker {
        if id, err := strconv.Atoi(_id); err == nil {
            user, err = db.GetUserById(ctx, id)
            if err != nil {
                log.Panic(err)
            }
        }
    }
    Jsonify(w, user)
}

func Users(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
    ctx := r.Context()
    id := s.Values["userid"].(int)
    user, err := db.GetUserById(ctx, id)
    if err != nil {
        log.Panic(err)
    }
    if user.Level < LevelWorker {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil {
        page = 0
    }
    size, err := strconv.Atoi(r.URL.Query().Get("size"))
    if err != nil {
        size = 1
    }
    users, err := db.GetUsersByPage(ctx, page, size)
    if err != nil {
        log.Panic(err)
    }
    Jsonify(w, users)
}

func Likes(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
    ctx := r.Context()
    id := s.Values["userid"].(int)
    courses, err := db.GetLikes(ctx, id)
    if err != nil {
        log.Panic(err)
    }
    Jsonify(w, courses)
}

func InLikes(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
    ctx := r.Context()
    id := s.Values["userid"].(int)
    course, err := strconv.Atoi(r.URL.Query().Get("course"))
    if err != nil {
        log.Println(err)
        return
    }
    flag, err := db.IsFollowed(ctx, id, course)
    if err != nil {
        log.Panic(err)
    }
    Jsonify(w, flag)
}

func Follow(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
    ctx := r.Context()
    user := s.Values["userid"].(int)
    course, err := strconv.Atoi(r.URL.Query().Get("course"))
    if err != nil {
        log.Println(err)
        return
    }
    if err = db.Follow(ctx, user, course); err != nil {
        log.Panic(err)
    }
}

func UnFollow(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
    ctx := r.Context()
    user := s.Values["userid"].(int)
    course, err := strconv.Atoi(r.URL.Query().Get("course"))
    if err != nil {
        log.Println(err)
        return
    }
    if err = db.UnFollow(ctx, user, course); err != nil {
        log.Panic(err)
    }
}

func Mining(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
    ctx := r.Context()
    user, err := Authority(r)
    if err != nil {
        log.Panic(err)
    }
    code := r.URL.Query().Get("code")
    if code == "" {
        http.Error(w, "code is empty", http.StatusBadRequest)
        return
    }
    count, err := db.Mining(ctx, user.Id, code)
    if err != nil {
        log.Println(err)
    }
    if count == 0 {
        return
    }
    user.Ticket += count
    err = db.SetUserTicket(ctx, user)
    if err != nil {
        log.Println(err)
    }
    http.Redirect(w, r, "/user.html", http.StatusSeeOther)
}

func GenMimes(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
    ctx := r.Context()
    user, err := Authority(r)
    if err != nil {
        log.Panic(err)
    }
    if user.Level < LevelWorker {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    total, err := strconv.Atoi(r.URL.Query().Get("total"))
    if err != nil {
        http.Error(w, "total is not a integer", http.StatusBadRequest)
        return
    }
    count, err := strconv.Atoi(r.URL.Query().Get("count"))
    if err != nil {
        http.Error(w, "count is not a integer", http.StatusBadRequest)
        return
    }
    if total < 1 || count < 1 {
        return
    }
    err = db.GenMime(ctx, total, count)
    if err != nil {
        log.Println(err)
    }
}

func GetMimes(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
    ctx := r.Context()
    user, err := Authority(r)
    if err != nil {
        log.Panic(err)
    }
    if user.Level < LevelWorker {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    total := 0
    paramTotal := r.URL.Query().Get("total")
    if paramTotal != "" {
        total, err = strconv.Atoi(paramTotal)
    }
    if err != nil {
        http.Error(w, "total is not a integer", http.StatusBadRequest)
        return
    }
    count, err := strconv.Atoi(r.URL.Query().Get("count"))
    if err != nil {
        http.Error(w, "count is not a integer", http.StatusBadRequest)
        return
    }
    if count < 1 {
        return
    }
    codes, err := db.GetMimes(ctx, total, count)
    if err != nil {
        log.Println(err)
    }
    Jsonify(w, codes)
}
