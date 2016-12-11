// Package controller file registration.go implements all handlers
// and helper functions for applicant registration.
package controller

import (
	"github.com/geobe/gostip/go/model"
	"github.com/geobe/gostip/go/view"
	"html"
	"html/template"
	"net/http"
)

// ShowRegistration is handler for registration form get requests:
// On GET, show an empty registration form, on POST save data into db
// and show a confirmation message.
func ShowRegistration(w http.ResponseWriter, r *http.Request) {
	db := model.Db()
	if err := parseSubmission(w, r); err != nil {
		return
	}
	values := viewmodel{
		// HTMLAttr unescapes string for use as HTML Attribute
		"disabled":     template.HTMLAttr(""),
		"action":       "register/submit",
		"oblasts":      model.Oblasts(),
		"district":     0,
		"buttons":      true,
		"displaystyle": "block",
	}
	if err := checkMethodAllowed(http.MethodGet, w, r); err == nil {
		var appId uint
		if v, ok := r.Form["appid"]; ok {
			appId = atouint(html.EscapeString(v[0]))
		} else {
			appId = 0
		}
		if appId > 0 {
			var app model.Applicant
			db.Preload("Data").Preload("Data.Oblast").First(&app, appId)
			if app.ID == appId {
				setViewModel(app, values)
			}
		}
		view.Views().ExecuteTemplate(w, "registration", values)
	}

}

// SubmitRegistration is handler for registration form post requests:
// On POST save data into db and show a confirmation message.
func SubmitRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		db := model.Db()
		if err := parseSubmission(w, r); err != nil {
			return
		}
		values := viewmodel{
			// HTMLAttr unescapes string for use as HTML Attribute
			"disabled":     template.HTMLAttr("disabled"),
			"action":       "register",
			"oblasts":      model.Oblasts(),
			"buttons":      false,
			"thankyou":     true,
			"displaystyle": "none",
		}
		var app model.Applicant
		appId := atoint(html.EscapeString(r.PostFormValue("appid")))
		if appId > 0 {
			var err error
			app, err = retrieveApplicant(appId, w)
			// update registration
			if err != nil {
				return
			}
			setApplicantData(&app, r, false)
			db.Save(&app)
		} else {
			// new registration
			app = model.Applicant{}
			appdata := model.ApplicantData{}
			app.Data = appdata
			setApplicantData(&app, r, false)
			db.Create(&app)
		}
		setViewModel(app, values)
		view.Views().ExecuteTemplate(w, "registration", values)
	}

}
