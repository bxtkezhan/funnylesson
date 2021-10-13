package handlers

import (
    "log"
    "os"
    "net/http"
    "strings"
    "encoding/json"
    "context"

    db "../database"

    "github.com/gorilla/sessions"
)

func Setup(){
    http.Handle("/", http.FileServer(http.Dir("./www")))
    http.HandleFunc("/img/", FileServer)
    http.HandleFunc("/api/", Handler)
}

type SecretFunc func(
    w http.ResponseWriter, r *http.Request, s *sessions.Session)

const (
    LevelVip     = 1<<iota
    LevelWorker  = 1<<iota
    LevelAdmin   = 1<<iota
    LevelDefault = 0

    SecretKey = "super-secret-key"
)

var (
    Routes       = make(map[string]http.HandlerFunc)
    SecretRoutes = make(map[string]SecretFunc)
    Store        = sessions.NewCookieStore([]byte(SecretKey))
)

func Handler(w http.ResponseWriter, r *http.Request) {
    path := strings.TrimPrefix(r.URL.Path, "/api")
    w.Header().Set("Content-Type", "application/json")
    if handler, found := Routes[path]; found {
        handler(w, r)
    } else if handler, found := SecretRoutes[path]; found {
        Secret(w, r, handler)
    } else {
        http.NotFound(w, r)
        log.Println(r.URL.Path)
    }
}

func Secret(w http.ResponseWriter, r *http.Request, handler SecretFunc) {
    session, err := Store.Get(r, "cookie-fl")
    if err == nil {
        if id, ok := session.Values["userid"].(int); !ok || id == 0 {
            http.Error(w, "Forbidden", http.StatusForbidden)
        } else {
            handler(w, r, session)
        }
    } else {
        log.Panic(err)
    }
}

func FileServer(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path[1:]
    if info, err := os.Stat(path); err == nil && !info.IsDir() {
        http.ServeFile(w, r, path)
    } else {
        http.NotFound(w, r)
    }
}

func Jsonify(w http.ResponseWriter, v interface{}) {
    err := json.NewEncoder(w).Encode(v)
    if err != nil {
        log.Panic(err)
    }
}

func Authority(r *http.Request) (*db.User, error) {
    ctx := r.Context()
    session, err := Store.Get(r, "cookie-fl")
    if err != nil {
        return nil, err
    }
    id, found := session.Values["userid"].(int)
    if !found || id == 0 {
        return &db.User{}, nil
    }
    ctx_, cancel := context.WithCancel(ctx)
    user, err := db.GetUserById(ctx_, id)
    cancel()
    if err != nil {
        return nil, err
    }
    return user, nil
}
