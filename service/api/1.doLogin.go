package api

/*
curl -X POST 'http://localhost:3000/session' -H 'Content-Type: application/json' -d '{"Nickname": "Nico"}'
*/

import (
	"encoding/json"
	"net/http"

	"github.com/annafabia03/WASA23/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user User
	err := json.NewDecoder(r.Body).Decode(&user) // Reads the JSON data from the request body (r.Body)
	//and decodes it into the user variable.
	if err != nil {
		stringErr := "doLogIn: invalid JSON object"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	//If there's an error decoding the JSON  it returns a 400 Bad Request error.

	if !user.HasValidUsername() { //Checks if the username provided in the request is valid. If it's not valid, it returns a 400 Bad Request error.
		stringErr := "doLogIn: invalid nickname"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	dbUser, present, err := rt.db.GetUserByNickname(user.Nickname) //Attempts to retrieve a user from the database based on the provided username.
	//It returns three values: the user data from the database, a boolean indicating
	//whether the user was found, and an error if any occurred.

	if err != nil { //If there's an error retrieving the user from the database, it returns a 500 Internal Server Error.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//If the user is present in the database (present == true), it returns the user data as JSON with status code 200 OK.

	if present {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		user.FromDatabase(dbUser)
		_ = json.NewEncoder(w).Encode(user)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		dbUser, err = rt.db.CreateUser(user.Nickname) //If the user is not present in the database (present == false),
		//it creates a new user in the database using rt.db.CreateUser(user.Nickname).
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user.FromDatabase(dbUser)
		_ = json.NewEncoder(w).Encode(user)
		return
	}

}
