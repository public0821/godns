package web

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/public0821/dnserver/db"
	//"fmt"
	"log"
	"net/http"
	//"strconv"
)

func deleteForwardServer(params martini.Params, r *http.Request) (rcode int, result string) {
	var fs db.ForwardServer
	fs.Ip = params["ip"]
	err := db.Delete(&fs)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, ""
}

func getForwardServer() (rcode int, result string) {
	rrcords, err := db.Query(&db.ForwardServer{}, 0, 0)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	if len(rrcords) == 0 {
		return http.StatusOK, "[]"
	}
	data, err := json.Marshal(rrcords)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, string(data)
}

func addForwardServer(r *http.Request) (rcode int, result string) {
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
	var fs db.ForwardServer
	err = json.Unmarshal(body, &fs)
	if err != nil {
		return http.StatusNotAcceptable, "content not correct json format"
	}
	log.Println(fs)
	err = db.Add(&fs)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, ""
}
