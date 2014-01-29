package main

import (
    "crypto/rand"
    "encoding/base64"
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
            if r.URL.Path != "/static/html/login.html" {
                http.Redirect(w, r, "/static/html/login.html", http.StatusTemporaryRedirect)
            }
            return
        }
    }
}

func main() {
    _, err := db.NewDBManager()
    if err != nil {
        log.Println(err)
        return
    }

    m := martini.Classic()

    //m.Use(Auth())

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

    m.Get("/login", func() string {
        return "login"
    })
    m.Get("/user", func() string {
        return "enable recursion mode"
    })
    m.Get("/static/**", martini.Static("./static", martini.StaticOptions{Prefix: "/static/"}))
    http.ListenAndServe(":8080", m)
}
