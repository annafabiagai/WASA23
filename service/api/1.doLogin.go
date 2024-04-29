package api

import (
	"encoding/json"
	"net/http"

	"github.com/annafabia03/WASA23/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		stringErr := "doLogIn: invalid JSON object"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	if !user.HasValidUsername() {
		stringErr := "doLogIn: invalid username"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	dbUser, present, err := rt.db.GetUserByNickname(user.Nickname)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if present {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		user.FromDatabase(dbUser)
		_ = json.NewEncoder(w).Encode(user)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		dbUser, err = rt.db.CreateUser(user.Nickname)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user.FromDatabase(dbUser)
		_ = json.NewEncoder(w).Encode(user)
		return
	}

}
