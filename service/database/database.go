/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/

package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	// USER TABLE

	CreateUser(nickname string) (dbUser User, err error) //create a new user in TABLE USER

	SetMyNickname(dbUser User) (err error) //update nickname

	GetUserByNickname(nickname string) (dbUser User, present bool, err error)

	GetUserByID(ID uint64) (dbUser User, present bool, err error)

	SearchUser(nicknameToSearch string) (userList []User, err error)

	GetProfile(ID uint64) (dbProfile Profile, err error)

	GetMyStream(requestingUserID uint64) (stream []Photo, err error)

	/// FOLLOW TABLE

	FollowUser(followerID uint64, followedID uint64) (err error)

	UnfollowUser(followerID uint64, followedID uint64) (err error)

	CheckFollow(followerID uint64, followedID uint64) (isFollowing bool, err error)

	GetFollowersList(followedID uint64) (followersList []User, err error)

	GetFollowingsList(followerID uint64) (followingsList []User, err error)

	RemoveFollowBothDirections(user1ID uint64, user2ID uint64) (err error)

	// BAN TABLE

	BanUser(bannerID uint64, bannedID uint64) (err error)

	UnbanUser(bannerID uint64, bannedID uint64) (err error)

	CheckBan(bannerID uint64, bannedID uint64) (isBanned bool, err error)

	CascadeBanBothDirections(user1ID uint64, user2ID uint64) (err error)

	/// PHOTO TABLE

	CreatePhoto(photo Photo) (dbPhoto Photo, err error)

	DeletePhoto(ID uint64) (err error)

	GetPhotoByID(ID uint64) (dbPhoto Photo, present bool, err error)

	GetPhotosList(ownerID uint64) (photosList []Photo, err error)

	/// LIKE TABLE

	LikePhoto(likerID uint64, photoID uint64) (err error)

	UnlikePhoto(likerID uint64, photoID uint64) (err error)

	GetLikesList(photoID uint64) (likesList []User, err error)

	RemoveLikesBothDirections(user1ID uint64, user2ID uint64) (err error)

	/// COMMENT TABLE

	CommentPhoto(comment Comment) (dbComment Comment, err error)

	UncommentPhoto(ID uint64) (err error)

	GetCommentsList(photoID uint64) (commentsList []Comment, err error)

	RemoveCommentsBothDirections(user1ID uint64, user2ID uint64) (err error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = createDatabase(db)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

func createDatabase(db *sql.DB) error {

	tables := [6]string{

		`CREATE TABLE IF NOT EXISTS user (
			userID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			nickname VARCHAR(16) NOT NULL UNIQUE
			);`,

		`CREATE TABLE IF NOT EXISTS follow(
			followerID INTEGER NOT NULL REFERENCES users (userID),
			followedID INTEGER NOT NULL REFERENCES users (userID),
			PRIMARY KEY (followerID, followedID)
			);`,

		`CREATE TABLE IF NOT EXISTS bann (
			bannerID INTEGER NOT NULL REFERENCES users (userID),
			bannedID INTEGER NOT NULL REFERENCES users (userID),
			PRIMARY KEY (bannerID, bannedID)
			);`,

		`CREATE TABLE IF NOT EXISTS photo (
			photoID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			ownerID INTEGER NOT NULL REFERENCES users (userID),
			format VARCHAR(3) NOT NULL,
			date TEXT NOT NULL
			);`,

		`CREATE TABLE IF NOT EXISTS like (
			likerID INTEGER NOT NULL REFERENCES users (userID),
			photoID INTEGER NOT NULL REFERENCES photos (photoID),
			PRIMARY KEY (likerID, photoID)
			);`,

		`CREATE TABLE IF NOT EXISTS comment (
			commentID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			photoID INTEGER NOT NULL REFERENCES photos (photoID),
			ownerID INTEGER NOT NULL REFERENCES users (userID),
			nickname VARCHAR(16) NOT NULL REFERENCES users (username),
			commentText TEXT NOT NULL,
			date TEXT NOT NULL
			);`,
	}

	for t := 0; t < len(tables); t++ {
		sqlStmt := tables[t]
		_, err := db.Exec(sqlStmt)

		if err != nil {
			return err
		}
	}

	return nil

}
