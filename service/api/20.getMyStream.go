package api

/*
go run ./cmd/webapi/
curl -X GET -H 'Authorization: 1' localhost:3000/home_pagecurl -X GET -H 'Authorization: 1' localhost:3000/home_page

non funziona: RITORNA SEMPRE NULL
*/

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/annafabia03/WASA23/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var token uint64
	token, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)

	// Unauthorized check
	if err != nil {
		stringErr := "getMyStream: invalid authorization token"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}
	requestingUser, present, err := rt.db.GetUserByID(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "getMyStream: authorization token not matching any existing user"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}

	// database section
	stream, err := rt.db.GetMyStream(requestingUser.ID)

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for i, photo := range stream {
		likesList, err := rt.db.GetLikesList(photo.ID)
		stream[i].LikeList = likesList
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		commentsList, err := rt.db.GetCommentsList(photo.ID)
		stream[i].CommentList = commentsList
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stream)
}
