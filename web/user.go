package web

import (
    "github.com/emicklei/go-restful"
    "github.com/public0821/dnserver/db"
    "net/http"
)

type UserResource struct {
}

func (r UserResource) auth(req *restful.Request, resp *restful.Response) {
    user := new(db.User)
    err := req.ReadEntity(user)
    if err != nil { // bad request
        resp.WriteErrorString(http.StatusBadRequest, err.Error())
        return
    }
    //resp.WriteEntity(db.User{Id: 1, Name: "test"})
}

func (r UserResource) updateOne(req *restful.Request, resp *restful.Response) {
    updatedUser := new(db.User)
    err := req.ReadEntity(updatedUser)
    if err != nil { // bad request
        resp.WriteErrorString(http.StatusBadRequest, err.Error())
        return
    }
    dbmap, err := db.OpenDbmap()
    if err != nil {
        resp.WriteErrorString(http.StatusInternalServerError, err.Error())
        return
    }
    defer db.CloseDbmap(dbmap)

    //dbmap.
}

func (r UserResource) Register(root string) {
    ws := new(restful.WebService)
    ws.Path(root + "/users")
    ws.Consumes(restful.MIME_JSON)
    ws.Produces(restful.MIME_JSON)

    ws.Route(ws.POST("/{id}").To(r.updateOne).
        Doc("update user password").
        Param(ws.PathParameter("id", "identifier of the product").DataType("string")))

    ws.Route(ws.POST("/auth").To(r.auth).
        Doc("update or create a product").
        Param(ws.BodyParameter("User", "a User (JSON)").DataType("db.User")))

    restful.Add(ws)
}
