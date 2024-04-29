package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/annafabia03/WASA23/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var token uint64
	token, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)

	// Unauthorized check
	if err != nil {
		stringErr := "setMyUserName: invalid authorization token"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}
	dbUser, present, err := rt.db.GetUserByID(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "setMyUserName: authorization token not matching any existing user"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}

	var updatedUser User
	updatedUser.FromDatabase(dbUser)
	err = json.NewDecoder(r.Body).Decode(&updatedUser)

	// BadRequest check
	if err != nil {
		stringErr := "setMyUserName: invalid JSON object"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	if !updatedUser.HasValidUsername() {
		stringErr := "setMyUserName: invalid nickname"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	_, present, err = rt.db.GetUserByNickname(updatedUser.Nickname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if present {
		stringErr := "setMyUserName: nickname already exists"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	// database section
	err = rt.db.SetMyNickname(updatedUser.ToDatabase())

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(updatedUser)
}
