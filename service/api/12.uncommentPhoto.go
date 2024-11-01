package api

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
	requestingUser, present, err := rt.db.SearchUserByID(token)
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
	pathPid, err = strconv.ParseUint(ps.ByName("pid"), 10, 64)

	// BadRequest check
	if err != nil {
		stringErr := "uncommentPhoto: invalid path parameter pid"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	dbPhoto, present, err := rt.db.SearchPhotoByID(pathPid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "uncommentPhoto: path parameter pid not matching any existing photo"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	var pathCid uint64
	pathCid, err = strconv.ParseUint(ps.ByName("cid"), 10, 64)

	// BadRequest check
	if err != nil {
		stringErr := "uncommentPhoto: invalid path parameter cid"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	comment, present, err := rt.db.SearchCommentByID(pathCid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "uncommentPhoto: path parameter cid not matching any existing comment"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	// Forbidden check
	isAuthor := requestingUser.ID == dbPhoto.AuthorID
	if requestingUser.ID != comment.AuthorID && !isAuthor {
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
