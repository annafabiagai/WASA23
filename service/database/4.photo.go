package database

import "database/sql"

func (db *appdbimpl) CreatePhoto(photo Photo) (dbPhoto Photo, err error) {

	query := "INSERT INTO photo (ownerID, format, date) VALUES (?,?,?);"

	sqlResult, err := db.c.Exec(query, photo.OwnerID, photo.Format, photo.Date)
	if err != nil {
		return
	}
	dbPhoto = photo
	photoID, err := sqlResult.LastInsertId()
	dbPhoto.ID = uint64(photoID)
	return
}

func (db *appdbimpl) DeletePhoto(ID uint64) (err error) {

	query := "DELETE FROM photo WHERE photoID = ?;"

	_, err = db.c.Exec(query, ID)
	return err
}

func (db *appdbimpl) GetPhotoByID(ID uint64) (dbPhoto Photo, present bool, err error) {

	query := "SELECT * FROM photo WHERE photoID = ?;"

	row := db.c.QueryRow(query, ID)
	err = row.Scan(&dbPhoto.ID, &dbPhoto.OwnerID, &dbPhoto.Format, &dbPhoto.Date)
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

func (db *appdbimpl) GetPhotosList(ownerID uint64) (photosList []Photo, err error) {

	query := "SELECT * FROM photo WHERE ownerID = ? ORDER BY date DESC;"

	rows, err := db.c.Query(query, ownerID)
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
		photosList = append(photosList, photo)
	}

	err = rows.Err()
	return
}
