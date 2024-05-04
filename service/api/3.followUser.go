package api

/*
go run ./cmd/webapi/
curl  -X PUT -H 'Authorization: 1' localhost:3000/following/2
*/

import (
	"net/http"
	"strconv"

	"github.com/annafabia03/WASA23/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var token uint64
	token, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)

	// Unauthorized check
	if err != nil {
		stringErr := "followUser: invalid authorization token"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}
	follower, present, err := rt.db.GetUserByID(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "followUser: authorization token not matching any existing user"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}

	var pathUid uint64
	pathUid, err = strconv.ParseUint(ps.ByName("userid"), 10, 64)

	// BadRequest check
	if err != nil {
		stringErr := "followUser: invalid path parameter user id"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	followed, present, err := rt.db.GetUserByID(pathUid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "followUser: path parameter user id not matching any existing user"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	if follower.ID == followed.ID {
		stringErr := "followUser: requesting user trying to follow himself"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	isFollowing, err := rt.db.CheckFollow(follower.ID, followed.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if isFollowing {
		stringErr := "followUser: requesting user already following user"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	// Forbidden check
	someoneIsBanned, err := rt.db.CheckBan(follower.ID, followed.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if someoneIsBanned {
		stringErr := "followUser: someone has banned the other"
		http.Error(w, stringErr, http.StatusForbidden)
		return
	}

	// database section
	err = rt.db.FollowUser(follower.ID, followed.ID)

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
