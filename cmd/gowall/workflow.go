package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"regexp"
	"encoding/json"
)

func handleXHR(c *gin.Context, data []byte) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Data(http.StatusOK, "application/json; charset=utf-8", data)
}

type Response struct {
	Success bool `json:"success" bson:"-"`
	Errors  []string `json:"errors" bson:"-"`
	ErrFor  map[string]string `json:"errfor" bson:"-"`

	Username    string  `json:"username"`
	Email   string  `json:"email"`
	Password string  `json:"password"`
}

func (r *Response)HasErrors() bool {
	return len(r.ErrFor) != 0 || len(r.Errors) != 0
}

func (r *Response)Fail(c *gin.Context) {
	r.Success = false
	c.JSON(http.StatusOK, r)
}

func (r *Response) Recover(c *gin.Context) {}

func (r *Response) ValidateUsername() {
	r.Username = strings.ToLower(r.Username)
	if len(r.Username) == 0 {
		r.ErrFor["username"] = "required"
	} else {
		ok, err := regexp.MatchString(`^[a-zA-Z0-9\-\_]+$`, r.Username)
		if err != nil {
			println(err.Error())
		}
		if !ok {
			r.ErrFor["username"] = `only use letters, numbers, \'-\', \'_\'`
		}
	}
}

func (r *Response) ValidateEmail() {
	r.Email = strings.ToLower(r.Email)
	if len(r.Email) == 0 {
		r.ErrFor["email"] = "required"
	} else {
		ok, err := regexp.MatchString(`^[a-zA-Z0-9\-\_\.\+]+@[a-zA-Z0-9\-\_\.]+\.[a-zA-Z0-9\-\_]+$`, r.Email)
		if err != nil {
			println(err.Error())
		}
		if !ok {
			r.ErrFor["email"] = `invalid email format`
		}
	}
}

func (r *Response) ValidatePassword() {
	if len(r.Password) == 0 {
		r.ErrFor["password"] = "required"
	} else {
		if len(r.Password) < 8 {
			r.ErrFor["password"] = `too weak password`
		}
	}
}

func (r *Response) CleanErrors() {
	r.Errors = []string{}
	r.ErrFor = make(map[string]string)
}

func (r *Response) DecodeRequest(c *gin.Context) {
	err := json.NewDecoder(c.Request.Body).Decode(r)
	if err != nil {
		panic(err)
	}
	// clean errors from client
	r.CleanErrors()
	return
}

func (r *Response) Finish(c *gin.Context) {
	r.Success = true
	c.JSON(http.StatusOK, r)
}
