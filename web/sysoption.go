package web

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/public0821/dnserver/db"
	"log"
	"net/http"
	//"strconv"
)

func getAllSysOption() (rcode int, result string) {
	key := db.SysOption{}
	modes, err := db.Query(&key, 0, 0)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	if len(modes) == 0 {
		return http.StatusOK, "[]"
	}
	data, err := json.Marshal(modes[0])
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, string(data)
}
func getSysOption(params martini.Params) (rcode int, result string) {
	key := db.SysOption{Name: params["name"]}
	modes, err := db.Query(&key, 0, 0)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	if len(modes) != 1 {
		return http.StatusInternalServerError, fmt.Sprintf("no %s configuration in db", params["name"])
	}
	data, err := json.Marshal(modes[0])
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, string(data)
}

func modifySysOption(r *http.Request) (rcode int, result string) {

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
	var sysoption db.SysOption
	err = json.Unmarshal(body, &sysoption)
	if err != nil {
		return http.StatusNotAcceptable, "content not correct json format"
	}
	log.Println(sysoption)
	err = db.Modify(&sysoption)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, ""
}
