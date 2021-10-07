package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	db "../database"

	"github.com/gorilla/sessions"
)

func init() {
    Routes["/signup"] = SignUp
    Routes["/login"] = Login
    SecretRoutes["/logout"] = Logout
}

func SignUp(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    email := strings.TrimSpace(r.PostFormValue("email"))
    username := strings.TrimSpace(r.PostFormValue("username"))
    password := strings.TrimSpace(r.PostFormValue("password"))
    if len(email) == 0 {
        log.Println("the email is empty")
        http.Redirect(w, r, "/signup.html", http.StatusSeeOther)
        return
    }
    _ctx, cancel := context.WithCancel(ctx)
    emailValided := db.ValidEmail(_ctx, email)
    cancel()
    if  !emailValided {
        log.Println("the email is already be used")
        http.Redirect(w, r, "/signup.html", http.StatusSeeOther)
        return
    }
    if len(username) < 2 {
        log.Println("username is too short")
        http.Redirect(w, r, "/signup.html", http.StatusSeeOther)
        return
    }
    if len(password) < 8 {
        log.Println("password is too short")
        http.Redirect(w, r, "/signup.html", http.StatusSeeOther)
        return
    }
    user := &db.User{
        Email: email,
        Username: username,
        Password: db.HashString(password),
        Introduction: "empty",
        Image: "default.png",
        Time: time.Now().Unix(),
        Level: LevelDefault}
    db.AddUser(ctx, user)
    http.Redirect(w, r, "/login.html", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    email := r.PostFormValue("email")
    password := r.PostFormValue("password")
    if id := db.ValidPassword(ctx, email, password); id != 0 {
        session, err := Store.Get(r, "cookie-fl")
        if err == nil {
            session.Values["userid"] = id
            session.Save(r, w)
        }
        http.Redirect(w, r, "/user.html", http.StatusSeeOther)
    } else {
        http.Redirect(w, r, "/login.html", http.StatusSeeOther)
    }
}

func Logout(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
    s.Values["userid"] = 0
    s.Save(r, w)
    http.Redirect(w, r, "/login.html?from=logout", http.StatusSeeOther)
}
