package database

func (db *appdbimpl) LikePhoto(likerID uint64, photoID uint64) (err error) {

	query := "INSERT INTO like (likerID, photoID) VALUES (?,?);"

	_, err = db.c.Exec(query, likerID, photoID)
	return
}

func (db *appdbimpl) UnlikePhoto(likerID uint64, photoID uint64) (err error) {

	query := "DELETE FROM like WHERE (likerID = ? AND photoID = ?);"

	_, err = db.c.Exec(query, likerID, photoID)
	return
}

func (db *appdbimpl) GetLikesList(photoID uint64) (likesList []User, err error) {

	query := "SELECT likerID FROM like WHERE photoID = ?;"

	rows, err := db.c.Query(query, photoID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID)
		if err != nil {
			return
		}
		user, _, err = db.GetUserByID(user.ID)
		if err != nil {
			return
		}
		likesList = append(likesList, user)
	}

	err = rows.Err()
	return
}

func (db *appdbimpl) RemoveLikesBothDirections(user1ID uint64, user2ID uint64) (err error) {

	query := `
		DELETE FROM like WHERE 
		likerID = ? 
		AND photoID IN (SELECT photoID FROM photo WHERE ownerID = ?);`

	_, err = db.c.Exec(query, user1ID, user2ID)
	if err != nil {
		return err
	}
	_, err = db.c.Exec(query, user2ID, user1ID)
	return
}
