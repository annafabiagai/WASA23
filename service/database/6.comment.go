package database

func (db *appdbimpl) CommentPhoto(comment Comment) (dbComment Comment, err error) {

	query := "INSERT INTO comment (photoID, ownerID, nickname, commentText, date) VALUES (?,?,?,?,?);"

	sqlResult, err := db.c.Exec(query, comment.PhotoID, comment.OwnerID, comment.Nickname, comment.Text, comment.Date)
	if err != nil {
		return
	}
	dbComment = comment
	commentID, err := sqlResult.LastInsertId()
	dbComment.ID = uint64(commentID)
	return
}

func (db *appdbimpl) UncommentPhoto(ID uint64) (err error) {

	query := "DELETE FROM comment WHERE commentID = ?;"

	_, err = db.c.Exec(query, ID)
	return
}

func (db *appdbimpl) GetCommentsList(photoID uint64) (commentsList []Comment, err error) {

	query := "SELECT * FROM comment WHERE photoID = ? ORDER BY date DESC;"

	rows, err := db.c.Query(query, photoID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.ID, &comment.PhotoID, &comment.OwnerID, &comment.Nickname, &comment.Text, &comment.Date)
		if err != nil {
			return
		}
		commentsList = append(commentsList, comment)
	}

	err = rows.Err()
	return
}

func (db *appdbimpl) RemoveCommentsBothDirections(user1ID uint64, user2ID uint64) (err error) {

	query := `
		DELETE FROM comment WHERE 
		photoID IN (SELECT photoID FROM photo WHERE ownerID = ?) 
		AND  ownerID = ?;`

	_, err = db.c.Exec(query, user1ID, user2ID)
	if err != nil {
		return err
	}
	_, err = db.c.Exec(query, user2ID, user1ID)
	return
}
