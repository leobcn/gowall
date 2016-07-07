package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"gopkg.in/mgo.v2"
	"strings"
)

type responseAdmin struct {
	Response
	Admin
}

func renderAdministrators(c *gin.Context) {
	query := bson.M{}

	name, ok := c.GetQuery("name")
	if ok && len(name) != 0 {
		query["name"] = bson.RegEx{
			Pattern: `^.*?` + name + `.*$`,
			Options: "i",
		}
	}

	var results []Admin

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)

	Result := getData(c, collection.Find(query), &results)

	Results, err := json.Marshal(Result)
	if err != nil {
		panic(err)
	}
	if XHR(c) {
		handleXHR(c, Results)
		return
	}

	c.Set("Results", template.JS(string(Results)))
	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

func createAdministrator(c *gin.Context) {
	response := responseAdmin{}
	response.BindContext(c)
	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf("root")
	if !ok {
		response.Errors = append(response.Errors, "You may not create administrators")
		response.Fail()
		return
	}

	response.Admin.DecodeRequest(c)
	if len(response.Name.Full) == 0 {
		response.Errors = append(response.Errors, "A name is required")
	}

	if response.HasErrors() {
		response.Fail()
		return
	}

	// handleName
	response.Name.Full = slugifyName(response.Name.Full)

	// duplicateAdministrator
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)
	err := collection.Find(bson.M{"name.full": response.Name.Full}).One(nil)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		response.Errors = append(response.Errors, "That administrator already exists.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}

	// handleName
	name := strings.Split(response.Name.Full, " ")
	response.Name.First = name[0]
	if len(name) == 2 {
		response.Name.Last = name[1]
	}
	if len(name) == 3 {
		response.Name.Middle = name[2]
	}

	// createAdministrator
	response.Admin.ID = bson.NewObjectId()
	err = collection.Insert(response.Admin) // todo I think mgo's behavior isn't expected
	if err != nil {
		panic(err)
		return
	}
	response.Success = true
	c.JSON(http.StatusOK, gin.H{"record": response, "success": true})
}

func readAdministrator(c *gin.Context) {
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)
	admin := Admin{}
	err := collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(&admin)
	if err != nil {
		if err == mgo.ErrNotFound {
			Status404Render(c)
			return
		}
		panic(err)
	}
	json, err := json.Marshal(admin)
	if err != nil {
		panic(err)
	}
	if XHR(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Data(http.StatusOK, "application/json; charset=utf-8", json)
		return
	}

	c.Set("Record", template.JS(getEscapedString(string(json))))
	c.HTML(http.StatusOK, "/admin/administrators/details/", c.Keys)
}

func updateAdministrator(c *gin.Context) {
	response := responseAdmin{}
	response.BindContext(c)

	err := json.NewDecoder(c.Request.Body).Decode(&response.Admin.Name)
	if err != nil {
		panic(err)
	}
	// clean errors from client
	response.CleanErrors()

	if len(response.Name.First) == 0 {
		response.Errors = append(response.Errors, "A name is required")
	}

	if len(response.Name.Last) == 0 {
		response.Errors = append(response.Errors, "A lastname is required")
	}

	if response.HasErrors() {
		response.Fail()
		return
	}

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)

	// patchAdministrator
	err = collection.UpdateId(bson.ObjectIdHex(c.Param("id")), response.Admin)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail()
		return
	}

	response.Finish()
}

func updateAdministratorPermissions(c *gin.Context) {
	response := responseAdmin{}
	response.BindContext(c)
	//TODO here is js bug
	//TODO there are not clear logic with populate of groups
	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf("root")
	if !ok {
		response.Errors = append(response.Errors, "You may not change the permissions of admin groups.")
		response.Fail()
		return
	}

	response.Admin.DecodeRequest(c)

	if len(response.Permissions) == 0 {
		response.Errors = append(response.Errors, "required")
	}

	if response.HasErrors() {
		response.Fail()
		return
	}

	//patchAdminGroup
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)

	err := collection.UpdateId(bson.ObjectIdHex(c.Param("id")), response.Admin)
	if err != nil {
		println(err.Error())
		panic(err)
	}

	response.Finish()
}

func deleteAdministrator(c *gin.Context) {
	response := Response{}
	response.BindContext(c)

	// validate
	if ok := getAdmin(c).IsMemberOf("root"); !ok {
		response.Errors = append(response.Errors, "You may not delete administrators.")
		response.Fail()
		return
	}

	// deleteUser
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)
	err := collection.RemoveId(bson.ObjectIdHex(c.Param("id")))
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail()
		return
	}

	response.Finish()
}
