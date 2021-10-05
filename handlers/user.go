package handlers

import (
    "fmt"
	"log"
	"net/http"
	"strconv"

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
    http.Redirect(w, r,
        fmt.Sprintf("/course.html?id=%d", course),
        http.StatusSeeOther)
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
    http.Redirect(w, r,
        fmt.Sprintf("/course.html?id=%d", course),
        http.StatusSeeOther)
}
