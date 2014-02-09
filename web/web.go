package web

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/public0821/dnserver/db"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var authedSession = make(map[string]bool)

func generateSessionid() (sid string, err error) {
	log.Println("generateSessionid")
	k := make([]byte, 32)
	if _, err = io.ReadFull(rand.Reader, k); err != nil {
		return
	}

	sid = base64.StdEncoding.EncodeToString(k)
	return
}

func Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectStatus := http.StatusUnauthorized
		if strings.HasPrefix(r.URL.Path, "/static/") {
			redirectStatus = http.StatusSeeOther
		}
		//static file and not html
		if strings.HasPrefix(r.URL.Path, "/static/") && !strings.HasSuffix(r.URL.Path, ".html") {
			return
		}
		//log.Println("111", r.Header)
		cookie, err := r.Cookie("sid")
		//no sid
		if err != nil {
			log.Println("no sid")
			cookie := &http.Cookie{}
			cookie.Name = "sid"
			cookie.Path = "/"
			cookie.Expires = time.Now().Add(time.Second * 3600)
			cookie.Value, err = generateSessionid()
			if err != nil {
				log.Println(err)
				http.Error(w, "Server internal error", http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/static/html/login.html", redirectStatus)
			return
		}
		//session not authed
		sid := cookie.Value
		authed, ok := authedSession[sid]
		if !(ok && authed) {
			//authed and post to /user/login
			if r.URL.Path == "/user/login/" {
				return
			}
			//if not start with /static/html/login, redirected to /static/html/login.html
			if !strings.HasPrefix(r.URL.Path, "/static/html/login") {
				log.Println(sid)
				log.Println(authedSession)
				log.Println("not authed")
				http.Redirect(w, r, "/static/html/login.html", redirectStatus)
				return
			}
			return
		}
		//session  authed
		if authed && r.URL.Path == "/static/html/login.html" {
			http.Redirect(w, r, "/static/html/home.html", http.StatusSeeOther)
			return
		}
		return
	}
}

func logout(w http.ResponseWriter, r *http.Request) (rcode int, result string) {
	cookie, err := r.Cookie("sid")
	if err != nil {
		log.Println("no sid")
		http.Redirect(w, r, "/static/html/login.html", http.StatusSeeOther)
		return
	}
	//session not authed
	sid := cookie.Value
	if authed, ok := authedSession[sid]; !ok || !authed {
		log.Println("not authed ")
		http.Redirect(w, r, "/static/html/login.html", http.StatusSeeOther)
		return
	}
	authedSession[sid] = false
	http.Redirect(w, r, "/static/html/login.html", http.StatusSeeOther)
	return
}

func login(w http.ResponseWriter, r *http.Request) (rcode int, result string) {
	cookie, err := r.Cookie("sid")
	//log.Println("100", r.Header)
	//log.Println(cookie)
	//no sid
	if err != nil {
		log.Println("no sid")
		//r.s
		http.Redirect(w, r, "/static/html/login.html", http.StatusSeeOther)
		return
	}
	//session authed
	sid := cookie.Value
	if authed, ok := authedSession[sid]; ok && authed {
		log.Println("already authed ")
		http.Redirect(w, r, "/static/html/home.html", http.StatusSeeOther)
		return
	}
	var user db.User
	user.Name = r.FormValue("name")
	user.Pwd = r.FormValue("pwd")
	log.Println(user)
	if len(user.Name) == 0 || len(user.Pwd) == 0 {
		http.Redirect(w, r, "/static/html/login.html", http.StatusSeeOther)
		return
	}
	users, err := db.Query(&user, 0, 0)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	if len(users) == 0 {
		errinfo := "user or password info incorrect"
		log.Println(errinfo)
		http.Redirect(w, r, "/static/html/login_error.html", http.StatusSeeOther)
		return
		//return http.StatusNotAcceptable, errinfo
	}
	authedSession[sid] = true
	http.Redirect(w, r, "/static/html/home.html", http.StatusSeeOther)
	return
}

type ChangePwdParam struct {
	Old_user     string
	Old_password string
	New_user     string
	New_password string
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
	var param ChangePwdParam
	err = json.Unmarshal(body, &param)
	if err != nil {
		return http.StatusNotAcceptable, "content not json format"
	}
	log.Println(string(body))
	log.Println(param)
	if len(param.Old_user) == 0 || len(param.Old_password) == 0 || len(param.New_user) == 0 || len(param.New_password) == 0 {
		return http.StatusNotAcceptable, "param can't be empty"
	}
	var user db.User
	user.Name = param.Old_user
	user.Pwd = param.Old_password
	users, err := db.Query(&user, 0, 0)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "db error"
	}
	if len(users) != 1 {
		return http.StatusNotAcceptable, "old user or old password error"
	}
	newUser, _ := users[0].(db.User)
	newUser.Name = param.New_user
	newUser.Pwd = param.New_password
	err = db.Modify(&newUser)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, ""
}

func Start() {

	m := martini.Classic()

	m.Use(Auth())

	m.Get("/sysoption/:name", getSysOption)
	m.Get("/sysoption/", getAllSysOption)
	m.Post("/sysoption/", modifySysOption)

	m.Get("/forwardserver/", getForwardServer)
	m.Post("/forwardserver/", addForwardServer)
	m.Delete("/forwardserver/:ip", deleteForwardServer)

	m.Get("/record", getRRecord)
	m.Get("/record/", getRRecord)
	m.Get("/record/count/", getRRecordCount)
	m.Post("/record/", addOrUpdateRRecord)
	m.Delete("/record/:id", deleteRRecord)

	m.Post("/user/login/", login)
	m.Get("/user/logout/", logout)
	m.Post("/user/chpassword/", changePassword)
	m.Get("/static/**", martini.Static("./static", martini.StaticOptions{Prefix: "/static/"}))
	http.ListenAndServe(":8080", m)
}
