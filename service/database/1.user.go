package database

import "database/sql"

func (db *appdbimpl) CreateUser(nickname string) (dbUser User, err error) {

	query := "INSERT INTO user (nickname) VALUES (?);"

	sqlResult, err := db.c.Exec(query, nickname)
	if err != nil {
		return
	}
	dbUser.Nickname = nickname
	userID, err := sqlResult.LastInsertId()
	dbUser.ID = uint64(userID)
	return
}

func (db *appdbimpl) SetMyNickname(dbUser User) (err error) {

	query := "UPDATE user SET nickname = ? WHERE userID = ?"
	_, err = db.c.Exec(query, dbUser.Nickname, dbUser.ID)
	if err != nil {
		return
	}
	return
}

func (db *appdbimpl) GetUserByNickname(nickname string) (dbUser User, present bool, err error) {
	query := "SELECT * FROM user WHERE nickname = ?;"

	err = db.c.QueryRow(query, nickname).Scan(&dbUser.ID, &dbUser.Nickname)
	if err != nil && err != sql.ErrNoRows {
		return
	} else if err == sql.ErrNoRows {
		err = nil
		return
	} else {
		err = nil
		present = true
		return
	}
}

func (db *appdbimpl) GetUserByID(ID uint64) (dbUser User, present bool, err error) {

	query := "SELECT * FROM user WHERE userID = ?;"

	err = db.c.QueryRow(query, ID).Scan(&dbUser.ID, &dbUser.Nickname)
	if err != nil && err != sql.ErrNoRows {
		return
	} else if err == sql.ErrNoRows {
		err = nil
		return
	} else {
		err = nil
		present = true
		return
	}
}

func (db *appdbimpl) SearchUser(nicknameToSearch string) (userList []User, err error) {

	query := "SELECT * FROM user WHERE nickname LIKE ?"

	rows, err := db.c.Query(query, nicknameToSearch+"%")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Nickname)
		if err != nil {
			return
		}
		userList = append(userList, user)
	}

	err = rows.Err()
	return
}

func (db *appdbimpl) GetProfile(ID uint64) (dbProfile Profile, err error) {

	queries := [5]string{
		"SELECT nickname FROM user WHERE userID = ?;",
		"SELECT COUNT(*) FROM photo WHERE ownerID = ?;",
		"SELECT COUNT(*) FROM follow WHERE followedID = ?;",
		"SELECT COUNT(*) FROM follow WHERE followerID = ?;",
		"SELECT COUNT(*) FROM like WHERE likerID =?;",
	}

	err = db.c.QueryRow(queries[0], ID).Scan(&dbProfile.Nickname)
	if err != nil {
		return
	}
	err = db.c.QueryRow(queries[1], ID).Scan(&dbProfile.PhotosCount)
	if err != nil {
		return
	}
	err = db.c.QueryRow(queries[2], ID).Scan(&dbProfile.Followers)
	if err != nil {
		return
	}
	err = db.c.QueryRow(queries[3], ID).Scan(&dbProfile.Followings)
	if err != nil {
		return
	}
	err = db.c.QueryRow(queries[4], ID).Scan(&dbProfile.LikesCount)
	return

}

func (db *appdbimpl) GetMyStream(requestingUserID uint64) (stream []Photo, err error) {

	query := `
		SELECT photoID, ownerID, format, date FROM photos
		INNER JOIN following ON ownerID = followedID
		WHERE followerID = ?
		ORDER BY date DESC
		LIMIT 50;`

	rows, err := db.c.Query(query, requestingUserID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var photo Photo
		err = rows.Scan(&photo.ID, &photo.OwnerID, &photo.Format, &photo.Date)
		if err != nil {
			return
		}
		owner, _, _ := db.GetUserByID(photo.OwnerID)
		photo.OwnerID = owner.ID
		stream = append(stream, photo)
	}

	err = rows.Err()
	return
}
