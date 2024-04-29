package api

import (
	"github.com/annafabia03/WASA23/service/album"
	"github.com/annafabia03/WASA23/service/database"
)

type Photo struct {
	ID          uint64             `json:"IDphoto"`
	OwnerID     uint64             `json:"IDuser"`
	Nickname    string             `json:"nickname"`
	Format      string             `json:"format"`
	CommentList []database.Comment `json:"commentList"`
	LikeList    []database.User    `json:"likeList"`
	Date        string             `json:"date"`
}

func (p *Photo) ToDatabase() database.Photo {
	return database.Photo{
		ID:          p.ID,
		OwnerID:     p.OwnerID,
		Nickname:    p.Nickname,
		Format:      p.Format,
		CommentList: p.CommentList,
		LikeList:    p.LikeList,
		Date:        p.Date,
	}
}

func (p *Photo) FromDatabase(photo database.Photo) {
	p.ID = photo.ID
	p.OwnerID = photo.OwnerID
	p.Nickname = photo.Nickname
	p.Format = photo.Format
	p.CommentList = photo.CommentList
	p.LikeList = photo.LikeList
	p.Date = photo.Date
}

func (p *Photo) ToAlbum() album.Photo {
	return album.Photo{
		ID:       p.ID,
		AuthorID: p.OwnerID,
		Format:   p.Format,
		Date:     p.Date,
	}
}

func (p *Photo) FromAlbum(photo album.Photo) {
	p.ID = photo.ID
	p.OwnerID = photo.AuthorID
	p.Format = photo.Format
	p.Date = photo.Date
}
