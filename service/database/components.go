package database

type User struct {
	ID       uint64 `json:"IDuser"`
	Nickname string `json:"nickname"`
}

type Photo struct {
	ID          uint64    `json:"IDphoto"`
	OwnerID     uint64    `json:"IDuser"`
	Nickname    string    `json:"nickname"`
	Format      string    `json:"format"`
	CommentList []Comment `json:"commentList"`
	LikeList    []User    `json:"likeList"`
	Date        string    `json:"date"`
}

type Comment struct {
	OwnerID  uint64 `json:"IDuser"`
	Nickname string `json:"ownerNickname"`
	ID       uint64 `json:"IDcomment"`
	PhotoID  uint64 `json:"IDphoto"`
	Text     string `json:"commentText"`
	Date     string `json:"date"`
}

type Profile struct {
	Nickname     string  `json:"nickname"`
	LikesCount   uint64  `json:"likesCount"`
	PhotosCount  uint64  `json:"photosCount"`
	Followers    []User  `json:"followersList"`
	Followings   []User  `json:"followingsList"`
	InBannedList bool    `json:"inMyBannedList"`
	MeBanned     bool    `json:"meBanned"`
	Photos       []Photo `json:"photoList"`
}
