package api

import (
	"github.com/annafabia03/WASA23/service/database"
)

type Comment struct {
	OwnerID  uint64 `json:"IDuser"`
	Nickname string `json:"ownerNickname"`
	ID       uint64 `json:"IDcomment"`
	PhotoID  uint64 `json:"IDphoto"`
	Text     string `json:"commentText"`
	Date     string `json:"date"`
}

func (c *Comment) ToDatabase() database.Comment {
	return database.Comment{
		OwnerID:  c.OwnerID,
		Nickname: c.Nickname,
		ID:       c.ID,
		PhotoID:  c.PhotoID,
		Text:     c.Text,
		Date:     c.Date,
	}
}

func (c *Comment) FromDatabase(comment database.Comment) {
	c.OwnerID = comment.OwnerID
	c.Nickname = comment.Nickname
	c.ID = comment.ID
	c.PhotoID = comment.PhotoID
	c.Text = comment.Text
	c.Date = comment.Date
}
