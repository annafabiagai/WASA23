package database

func (db *appdbimpl) FollowUser(followerID uint64, followedID uint64) (err error) {

	query := "INSERT INTO follow (followerID, followedID) VALUES (?,?);"

	_, err = db.c.Exec(query, followerID, followedID)
	return err
}

func (db *appdbimpl) UnfollowUser(followerID uint64, followedID uint64) (err error) {

	query := "DELETE FROM follow WHERE (followerID = ? AND followedID = ?);"

	_, err = db.c.Exec(query, followerID, followedID)
	return err
}

func (db *appdbimpl) CheckFollow(followerID uint64, followedID uint64) (isFollowing bool, err error) {

	query := "SELECT EXISTS (SELECT '_' FROM follow WHERE followerID = ? AND followedID = ?);"

	err = db.c.QueryRow(query, followerID, followedID).Scan(&isFollowing)
	return
}

func (db *appdbimpl) GetFollowersList(followedID uint64) (followersList []User, err error) {

	query := "SELECT followerID FROM follow WHERE followedID = ?;"

	rows, err := db.c.Query(query, followedID)
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
		followersList = append(followersList, user)
	}

	err = rows.Err()
	return
}

func (db *appdbimpl) GetFollowingsList(followerID uint64) (followingsList []User, err error) {

	query := "SELECT followedID FROM follow WHERE followerID = ?;"

	rows, err := db.c.Query(query, followerID)
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
		followingsList = append(followingsList, user)
	}

	err = rows.Err()
	return
}

func (db *appdbimpl) RemoveFollowBothDirections(user1ID uint64, user2ID uint64) (err error) {
	err = db.UnfollowUser(user1ID, user2ID)
	if err != nil {
		return err
	}
	err = db.UnfollowUser(user2ID, user1ID)
	return err
}
