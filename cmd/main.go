package main

import (
	"html/template"
	"io"

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

type Contact struct {
	Name string
	Email string
}

func newContact(name, email string) *Contact {
	return &Contact {
		Name: name,
		Email: email,
	}
}

type Contacts = map[string]*Contact

type ContactData struct {
	Contacts Contacts
}

func newContactData() *ContactData {
	return &ContactData{
		Contacts: map[string]*Contact{
			"john@example.com": newContact("John", "john@example.com"),
			"joe@example.com":  newContact("Joe", "joe@example.com"),
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
			"errorMessage" : "",
			"name" : "",
			"email" : "",
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

	e.Renderer = newTemplate()
	pageData := newPageData()

	e.GET("/", func (c echo.Context) error {
		return c.Render(200, "index", pageData)
	})
	e.POST("/contacts", func (c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if pageData.ContactData.hasEmail(email){
			pageData.FormData["name"] = name
			pageData.FormData["email"] = email
			pageData.FormData["errorMessage"] = "Try a different email"

			return c.Render(422, "createContact", pageData.FormData)
		}
		contact := newContact(name, email)
		pageData.ContactData.Contacts[email] = contact

		c.Render(200,"createContact",newFormData())

		return c.Render(200, "oob-contact", contact)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
