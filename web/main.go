package main

import (
    "crypto/rand"
    "encoding/base64"
    "encoding/json"
    "github.com/codegangsta/martini"
    "github.com/public0821/dnserver/db"
    "io"
    "log"
    "net/http"
    "time"
)

var authedSession = make(map[string]bool)

func generateSessionid() (sessionid string, err error) {
    k := make([]byte, 32)
    if _, err = io.ReadFull(rand.Reader, k); err != nil {
        return
    }

    sessionid = base64.StdEncoding.EncodeToString(k)
    return
}

func Auth() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("sessionid")
        //no sessionid
        if err != nil {
            cookie := &http.Cookie{}
            cookie.Name = "sessionid"
            cookie.Expires = time.Now().Add(time.Second * 3600)
            cookie.Value, err = generateSessionid()
            if err != nil {
                log.Println(err)
                http.Error(w, "Server internal error", http.StatusInternalServerError)
                return
            }
            http.SetCookie(w, cookie)
            http.Redirect(w, r, "/static/html/login.html", http.StatusTemporaryRedirect)
            return
        }
        if r.URL.Path == "/login" {
            return
        }
        //session not authed
        sessionid := cookie.Value
        if authed, ok := authedSession[sessionid]; !(ok && authed) {
            if r.URL.Path != "/static/html/login.html" && r.URL.Path[:len("/static/html")] == "/static/html" {
                log.Println("1", r.URL.Path)
                http.Redirect(w, r, "/static/html/login.html", http.StatusTemporaryRedirect)
            }
            log.Println("2", r.URL.Path)
            return
        }
    }
}

func login(r *http.Request) (rcode int, result string) {
    if r.ContentLength > 1000 {
        return http.StatusNotAcceptable, "content length too long"
    }
    body := make([]byte, r.ContentLength)
    length, err := r.Body.Read(body)
    if int64(length) != r.ContentLength {
        return http.StatusInternalServerError, "reed body data error"
    }
    if err != nil {
        return http.StatusInternalServerError, err.Error()
    }
    var user db.User
    err = json.Unmarshal(body, &user)
    if err != nil {
        return http.StatusNotAcceptable, "content not json format"
    }
    log.Println(user)
    users, err := db.Query(&user, 0, 0)
    if err != nil {
        return http.StatusInternalServerError, err.Error()
    }
    if len(users) != 1 {
        log.Println("user info not correct in db")
        return http.StatusInternalServerError, ""
    }
    return http.StatusOK, ""
}

func changePassword(r *http.Request) (rcode int, result string) {
    if r.ContentLength > 1000 {
        return http.StatusNotAcceptable, "content length too long"
    }
    body := make([]byte, r.ContentLength)
    length, err := r.Body.Read(body)
    if int64(length) != r.ContentLength {
        return http.StatusInternalServerError, "reed body data error"
    }
    if err != nil {
        return http.StatusInternalServerError, err.Error()
    }
    var user db.User
    err = json.Unmarshal(body, &user)
    if err != nil {
        return http.StatusNotAcceptable, "content not json format"
    }
    log.Println(user)
    err = db.Modify(&user)
    if err != nil {
        return http.StatusInternalServerError, err.Error()
    }
    return http.StatusOK, string(body)
}

func main() {
    _, err := db.NewDBManager()
    if err != nil {
        log.Println(err)
        return
    }

    m := martini.Classic()

    m.Use(Auth())

    m.Get("/sysoption/:name", getSysOption)
    m.Get("/sysoption/", getAllSysOption)
    m.Post("/sysoption/", modifySysOption)

    m.Get("/forwardserver/", getForwardServer)
    m.Post("/forwardserver/", addForwardServer)
    m.Delete("/forwardserver/", deleteForwardServer)

    m.Get("/record", getRRecord)
    m.Get("/record/", getRRecord)
    m.Get("/record/count/", getRRecordCount)
    m.Post("/record/", addOrUpdateRRecord)
    m.Delete("/record/:id", deleteRRecord)

    m.Get("/login", login)
    m.Get("/user", func() string {
        return "enable recursion mode"
    })
    m.Get("/static/**", martini.Static("./static", martini.StaticOptions{Prefix: "/static/"}))
    http.ListenAndServe(":8080", m)
}
