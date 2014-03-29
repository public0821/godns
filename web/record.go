package web

import (
    "encoding/json"
    "github.com/emicklei/go-restful"
    "github.com/public0821/dnserver/db"
    "log"
    "net/http"
    "strconv"
)

type RecordResource struct {
}

func (r RecordResource) query(req *restful.Request, resp *restful.Response) {
    name := req.QueryParameter("name")
    rtype := req.QueryParameter("type")
    value := req.QueryParameter("value")
    start_str := req.QueryParameter("start")
    end_str := req.QueryParameter("end")
    log.Println(name, rtype, value, start_str, end_str)
    start := 0
    end := 0
    if len(start_str) != 0 {
        number, err := strconv.Atoi(start_str)
        if err != nil {
            resp.WriteErrorString(http.StatusBadRequest, err.Error())
            return
        }
        start = number
    }
    if len(end_str) != 0 {
        number, err := strconv.Atoi(end_str)
        if err != nil {
            resp.WriteErrorString(http.StatusBadRequest, err.Error())
            return
        }
        end = number
    }

    dbmap, err := db.OpenDbmap()
    if err != nil {
        resp.WriteErrorString(http.StatusInternalServerError, err.Error())
        return
    }
    defer db.CloseDbmap(dbmap)

    var records []db.Record
    sql := "select * from record"
    where := ""
    if len(name) > 0:
        where += "name like"
    if start >0 && end >0:
    _, err = dbmap.Select(&records, "select * from record where name like ? and type=? and value like ? limit ? offset ?", name, rtype, value, end-start+1, start)
    if err != nil { // bad request
        resp.WriteErrorString(http.StatusBadRequest, err.Error())
        return
    }
    resp.WriteEntity(records)

}

func (r RecordResource) addOne(req *restful.Request, resp *restful.Response) {
    log.Println("add one")
    record := new(db.Record)
    err := req.ReadEntity(record)
    if err != nil { // bad request
        resp.WriteErrorString(http.StatusBadRequest, err.Error())
        return
    }
    data, err := json.Marshal(record)
    log.Println(string(data))
    dbmap, err := db.OpenDbmap()
    if err != nil {
        resp.WriteErrorString(http.StatusInternalServerError, err.Error())
        return
    }
    defer db.CloseDbmap(dbmap)

    err = dbmap.Insert(record)
    if err != nil {
        resp.WriteErrorString(http.StatusInternalServerError, err.Error())
        return
    }
}

func (r RecordResource) updateOne(req *restful.Request, resp *restful.Response) {
    log.Println("update one")
    //updatedUser := new(db.User)
    //err := req.ReadEntity(updatedUser)
    //if err != nil { // bad request
    //resp.WriteErrorString(http.StatusBadRequest, err.Error())
    //return
    //}
    //dbmap, err := db.OpenDbmap()
    //if err != nil {
    //resp.WriteErrorString(http.StatusInternalServerError, err.Error())
    //return
    //}
    //defer db.CloseDbmap(dbmap)

    //dbmap.
}

func (r RecordResource) Register(root string) {
    ws := new(restful.WebService)
    ws.Path(root + "/records")
    ws.Consumes(restful.MIME_JSON)
    ws.Produces(restful.MIME_JSON)

    ws.Route(ws.GET("").To(r.query).
        Doc("query record").
        Param(ws.QueryParameter("name", "record name").DataType("string")).
        Param(ws.QueryParameter("type", "record type").DataType("uint16")).
        Param(ws.QueryParameter("value", "record value").DataType("string")).
        Param(ws.QueryParameter("start", "record value").DataType("uint")).
        Param(ws.QueryParameter("end", "record value").DataType("uint")))

    ws.Route(ws.POST("").To(r.addOne).
        Doc("create a record").
        Param(ws.BodyParameter("Record", "a Record (JSON)").DataType("db.Record")))
    ws.Route(ws.POST("/{id}").To(r.updateOne).
        Doc("update record").
        Param(ws.PathParameter("id", "identifier of the record").DataType("string")))

    //ws.Route(ws.DELETE("/{id}").To(r.deleteOne).
    //Doc("delete a record").
    //Param(ws.PathParameter("id", "identifier of the record").DataType("string")))

    //m.Get("/record", getRRecord)
    //m.Get("/record/", getRRecord)
    //m.Get("/record/count/", getRRecordCount)
    //m.Post("/record/", addOrUpdateRRecord)
    //m.Delete("/record/:id", deleteRRecord)
    restful.Add(ws)
}
