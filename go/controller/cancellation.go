// Package controller file enrol.go implements handler and helper functions for enrolment.
// I.e. for tabs enrol and edit.
package controller

import (
	"github.com/geobe/gostip/go/view"
	"html"
	"net/http"
	"github.com/geobe/gostip/go/model"
)

// ShowCancellation is handler to show the selected applicant from the
// search select element for cancellation tab. It returns
// an html page fragment that is inserted into the respective tab area.
func ShowCancellation(w http.ResponseWriter, r *http.Request) {
	if checkMethodAllowed(http.MethodPost, w, r) != nil {
		return
	}
	if err := parseSubmission(w, r); err != nil {
		return
	}
	flag := html.EscapeString(r.PostFormValue("flag"))
	app, err := fetchApplicant(w, r, "appid", flag != "")
	if err != nil {
		return
	}

	values := viewmodel{
	}
	setViewModel(app, values)
	view.Views().ExecuteTemplate(w, "work_cancellation", values)
}

// SubmitCancellation is handler that accepts form submissions from the cancellation tab.
// Only http POST method is accepted.
func SubmitCancelation(w http.ResponseWriter, r *http.Request) {
	if err := checkMethodAllowed(http.MethodPost, w, r); err == nil {
		processCancellation(w, r, true)
	}
}

// processCancellation deletes or undeletes an applicant in the database
func processCancellation(w http.ResponseWriter, r *http.Request, enrol bool) {
	if checkMethodAllowed(http.MethodPost, w, r) != nil {
		return
	}
	if err := parseSubmission(w, r); err != nil {
		return
	}
	flag := html.EscapeString(r.PostFormValue("flag"))
	undo := flag != ""
	app, err := fetchApplicant(w, r, "appid", undo)
	if err == nil {
		if undo {
			app.DeletedAt = nil
			model.Db().Unscoped().Save(&app)
		} else {
			model.Db().Delete(&app)
		}
		w.WriteHeader(http.StatusNoContent)
	}
}