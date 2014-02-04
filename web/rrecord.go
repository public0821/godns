package main

import (
	"github.com/codegangsta/martini"
	"github.com/public0821/dnserver/db"
	//"github.com/codegangsta/martini-contrib/sessions"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func getRRecordCount() (rcode int, result string) {
	var rrcord db.RRecord
	count, err := db.Count(&rrcord)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, fmt.Sprintf("{\"count\":%d}", count)

}

func deleteRRecord(params martini.Params) (rcode int, result string) {
	var rrcord db.RRecord
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return http.StatusNotAcceptable, err.Error()
	}
	rrcord.Id = int64(id)
	err = db.Delete(&rrcord)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, ""
}

func getRRecord(r *http.Request) (rcode int, result string) {
	start := 0
	offset := 0
	var record db.RRecord
	record.Name = r.FormValue("name")
	tempType, _ := strconv.Atoi(r.FormValue("type"))
	record.Type = uint16(tempType)
	record.Value = r.FormValue("value")
	startStr := r.FormValue("start")
	offsetStr := r.FormValue("offset")
	if len(startStr) != 0 {
		number, err := strconv.Atoi(startStr)
		if err != nil {
			return http.StatusNotAcceptable, err.Error()
		}
		start = number
	}
	if len(offsetStr) != 0 {
		number, err := strconv.Atoi(offsetStr)
		if err != nil {
			return http.StatusNotAcceptable, err.Error()
		}
		offset = number
	}
	rrcords, err := db.Query(&record, start, offset)
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

func addOrUpdateRRecord(r *http.Request) (rcode int, result string) {
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
	var rrcord db.RRecord
	err = json.Unmarshal(body, &rrcord)
	if err != nil {
		log.Println(string(body))
		return http.StatusNotAcceptable, "content not json format"
	}
	log.Println(rrcord)
	if rrcord.Id == 0 {
		err = db.Add(&rrcord)
	} else {
		err = db.Modify(&rrcord)
	}
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, string(body)
}
