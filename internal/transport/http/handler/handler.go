package handler

import (
	"html/template"
	"net/http"
)

type Handler struct {
}

type ErrorResponse struct {
	Message string
	Code    int
}

func errorResponse(w http.ResponseWriter, message string, status int) {
	response := ErrorResponse{Message: message, Code: status}
	html, err := template.ParseFiles("./ui/error.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err = html.Execute(w, response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
