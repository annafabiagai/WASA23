package api

/*
go run ./cmd/webapi/
curl -X PUT -H 'Authorization: 1' localhost:3000/banned/2
*/

import (
	"net/http"
	"strconv"

	"github.com/annafabia03/WASA23/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var token uint64
	token, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)

	// Unauthorized check
	if err != nil {
		stringErr := "banUser: invalid authorization token"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}
	banner, present, err := rt.db.GetUserByID(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "banUser: authorization token not matching any existing user"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}

	var pathUid uint64
	pathUid, err = strconv.ParseUint(ps.ByName("userid"), 10, 64)

	// BadRequest check
	if err != nil {
		stringErr := "banUser: invalid path parameter user id"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	banned, present, err := rt.db.GetUserByID(pathUid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "banUser: path parameter user id not matching any existing user"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	if banner.ID == banned.ID {
		stringErr := "banUser: requesting user trying to ban himself"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	isBanned, err := rt.db.CheckBan(banner.ID, banned.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if isBanned {
		stringErr := "banUser: requesting user already banned user"
		http.Error(w, stringErr, http.StatusForbidden)
		return
	}

	// database section
	err = rt.db.BanUser(banner.ID, banned.ID)

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// cascade ban
	err = rt.db.CascadeBanBothDirections(banner.ID, banned.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
