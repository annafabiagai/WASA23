package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/annafabia03/WASA23/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var token uint64
	token, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)

	// Unauthorized check
	if err != nil {
		stringErr := "getUserId: invalid authorization token"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}
	_, present, err := rt.db.SearchUserByID(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "getUserId: authorization token not matching any existing user"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}

	username := ps.ByName("username")
	if (len(username) < 3 && len(username) > 16){
		w.WriteHeader(http.StatusBadRequest)
		return
	} 
		
	/*

	// BadRequest check
	if err != nil {
		stringErr := "getUserId: invalid path parameter username"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	*/
	requestedUser, present, err := rt.db.SearchUserByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "getUserId: path parameter username not matching any existing user"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(requestedUser.ID)
}
