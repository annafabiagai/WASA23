package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/annafabia03/WASA23/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) searchUserByUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var token uint64
	token, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)

	// Unauthorized check
	if err != nil {
		stringErr := "searchUserByUsername: invalid authorization token"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}
	_, present, err := rt.db.SearchUserByID(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "searchUserByUsername: authorization token not matching any existing user"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}

	usernameToSearch := r.URL.Query().Get("username")

	// database section
	usersList, err := rt.db.SearchUser(usernameToSearch)

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(usersList) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	for _, userToSearch := range usersList {
		var banned_check, err = rt.db.CheckBan(userToSearch.ID, token)
		// InternalServerError check
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if banned_check {
			http.Error(w, "User has been banned from this service.", http.StatusForbidden)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(usersList)
}
