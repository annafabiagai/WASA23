package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/annafabia03/WASA23/service/album"
	"github.com/annafabia03/WASA23/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var token uint64
	token, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)

	// Unauthorized check
	if err != nil {
		stringErr := "uploadPhoto: invalid authorization token"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}
	owner, present, err := rt.db.GetUserByID(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "uploadPhoto: authorization token not matching any existing user"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}

	binaryData, err := io.ReadAll(r.Body)

	// BadRequest check
	if err != nil {
		stringErr := "uploadPhoto: invalid binary data"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	format := http.DetectContentType(binaryData)
	switch format {
	case "image/png":
		format = "png"
	case "image/jpeg":
		format = "jpg"
	default:
		stringErr := "uploadPhoto: binary data not png/jpg"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	binaryImage := binaryData
	photo := Photo{
		OwnerID: owner.ID,
		Format:  format,
		Date:    time.Now().Format("2006-01-02 15:04:05"),
	}

	// database section
	dbPhoto, err := rt.db.CreatePhoto(photo.ToDatabase())
	photo.FromDatabase(dbPhoto)

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fileSystem section
	err = album.CreatePhotoFile(photo.ToAlbum(), binaryImage)

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(photo)

}
