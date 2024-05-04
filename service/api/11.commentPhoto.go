package api

/*
go run ./cmd/webapi/
curl -X POST -H 'Content-Type: text/plain' -H 'Authorization: 2' -d "che figata" http://localhost:3000/photos/3/comments/
*/

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/annafabia03/WASA23/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var token uint64
	token, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)

	// Unauthorized check
	if err != nil {
		stringErr := "commentPhoto: invalid authorization token"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}
	owner, present, err := rt.db.GetUserByID(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "commentPhoto: authorization token not matching any existing user"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}

	var pathPid uint64
	pathPid, err = strconv.ParseUint(ps.ByName("photoid"), 10, 64)

	// BadRequest check
	if err != nil {
		stringErr := "commentPhoto: invalid path parameter photoid"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	photo, present, err := rt.db.GetPhotoByID(pathPid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "commentPhoto: path parameter photo id not matching any existing photo"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	// Forbidden check
	someoneIsBanned, err := rt.db.CheckBan(owner.ID, photo.OwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if someoneIsBanned {
		stringErr := "commentPhoto: someone has banned the other"
		http.Error(w, stringErr, http.StatusForbidden)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		stringErr := "commentPhoto: invalid request body"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	commentText := string(body)
	comment := Comment{
		PhotoID:  photo.ID,
		OwnerID:  owner.ID,
		Nickname: owner.Nickname,
		Text:     commentText,
		Date:     time.Now().Format("2006-01-02 15:04:05"),
	}

	// database section
	dbComment, err := rt.db.CommentPhoto(comment.ToDatabase())
	comment.FromDatabase(dbComment)

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(comment)
}
