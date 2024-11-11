package main

import (
	"html/template"
	"io"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

var id = 0
type Contact struct {
	Id int
	Name string
	Email string
}

func newContact(name, email string) *Contact {
	id++
	return &Contact {
		Name: name,
		Email: email,
		Id: id,
	}
}

type Contacts = map[string]*Contact

type ContactData struct {
	Contacts Contacts
}

func newContactData() *ContactData {
	return &ContactData{
		Contacts: map[string]*Contact{
			"john@example.com" : newContact("John", "john@example.com"),
			"joe@example.com" :  newContact("Joe", "joe@example.com"),
			"king@example.com": newContact("King", "king@example.com"),
		},
	}
}

func (contactData *ContactData) hasEmail(email string) bool {
	_, exists := contactData.Contacts[email]
	if exists {
		return true
	}
	return false
}

type FormData = map[string]string

func newFormData() FormData{
	return FormData{
			"ErrorMessage" : "",
			"Name" : "",
			"Email" : "",
		}
}

type PageData struct{
	ContactData ContactData
	FormData FormData
}

func newPageData() PageData{
	return PageData{
		ContactData: *newContactData(),
		FormData: newFormData(),
	}
}

func main()  {
	e := echo.New()
	e.Use(middleware.Logger())

	e.File("/favicon.ico", "images/favicon.ico")
	e.Static("/images", "images")
	e.Static("/css", "css")

	e.Renderer = newTemplate()
	pageData := newPageData()

	e.GET("/", func (c echo.Context) error {
		return c.Render(200, "index", pageData)
	})

	e.GET("/favicon.ico", func(c echo.Context) error {
		return c.File("images/favicon.ico")
	})

	e.POST("/contacts", func (c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if pageData.ContactData.hasEmail(email){
			pageData.FormData["Name"] = name
			pageData.FormData["Email"] = email
			pageData.FormData["ErrorMessage"] = "Try a different email"

			return c.Render(422, "createContact", pageData.FormData)
		}

		contact := newContact(name, email)
		pageData.ContactData.Contacts[email] = contact

		c.Render(200,"createContact",newFormData())

		return c.Render(200, "oob-contact", contact)
	})

	e.DELETE("/contacts/:email", func(c echo.Context) error{
		time.Sleep(3*time.Second)
		email := c.Param("email")

		if !pageData.ContactData.hasEmail(email){
			return c.NoContent(200)
		}
		delete(pageData.ContactData.Contacts, email)
		return c.NoContent(200)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
