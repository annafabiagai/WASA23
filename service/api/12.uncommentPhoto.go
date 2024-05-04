package api

/*
go run ./cmd/webapi/
curl -v \
	-X DELETE \
	-H 'Authorization: 1' \
	localhost:3000/photos/{1}/comments/{2}
*/

import (
	"net/http"
	"strconv"

	"github.com/annafabia03/WASA23/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var token uint64
	token, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)

	// Unauthorized check
	if err != nil {
		stringErr := "uncommentPhoto: invalid authorization token"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}
	requestingUser, present, err := rt.db.GetUserByID(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "uncommentPhoto: authorization token not matching any existing user"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}

	var pathPid uint64
	pathPid, err = strconv.ParseUint(ps.ByName("photoid"), 10, 64)

	// BadRequest check
	if err != nil {
		stringErr := "uncommentPhoto: invalid path parameter photoid"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	dbPhoto, present, err := rt.db.GetPhotoByID(pathPid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "uncommentPhoto: path parameter photoid not matching any existing photo"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	var pathCid uint64
	pathCid, err = strconv.ParseUint(ps.ByName("commentid"), 10, 64)

	// BadRequest check
	if err != nil {
		stringErr := "uncommentPhoto: invalid path parameter comment id"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	comment, present, err := rt.db.SearchCommentByID(pathCid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "uncommentPhoto: path parameter comment id not matching any existing comment"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	// Forbidden check
	isAuthor := requestingUser.ID == dbPhoto.OwnerID
	if requestingUser.ID != comment.OwnerID && !isAuthor {
		stringErr := "uncommentPhoto: requesting user not author of the comment"
		http.Error(w, stringErr, http.StatusForbidden)
		return
	}

	// database section
	err = rt.db.UncommentPhoto(comment.ID)

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
