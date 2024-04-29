package api

import (
	"github.com/annafabia03/WASA23/service/database"
)

type Profile struct {
	Nickname      string           `json:"nickname"`
	LikesCount    uint64           `json:"likesCount"`
	PhotosCount   uint64           `json:"photosCount"`
	Followers     []database.User  `json:"followersList"`
	Followings    []database.User  `json:"followingsList"`
	InBannedList  bool             `json:"inMyBannedList"`
	MeBanned      bool             `json:"meBanned"`
	IsItMe        bool             `json:"isItMe"`
	DoIFollowUser bool             `json:"doIFollowUser"`
	Photos        []database.Photo `json:"photoList"`
}

func (p *Profile) ToDatabase() database.Profile {
	return database.Profile{
		Nickname:      p.Nickname,
		LikesCount:    p.LikesCount,
		PhotosCount:   p.PhotosCount,
		Followers:     p.Followers,
		Followings:    p.Followings,
		InBannedList:  p.InBannedList,
		MeBanned:      p.MeBanned,
		IsItMe:        p.IsItMe,
		DoIFollowUser: p.DoIFollowUser,
		Photos:        p.Photos,
	}
}

func (p *Profile) FromDatabase(profile database.Profile) {
	p.Nickname = profile.Nickname
	p.LikesCount = profile.LikesCount
	p.PhotosCount = profile.PhotosCount
	p.Followers = profile.Followers
	p.Followings = profile.Followings
	p.Photos = profile.Photos
}
