package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) LikePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println(45)
	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}
