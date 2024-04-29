package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	rt.router.POST("/session", rt.wrap(rt.doLogin))

	rt.router.PUT("/user", rt.wrap(rt.setMyUserName))

	rt.router.GET("/users/:userid", rt.wrap(rt.getUserProfile))

	rt.router.GET("/users/", rt.wrap(rt.searchUserByUsername))

	rt.router.GET("/user/:username", rt.wrap(rt.getUserId))

	rt.router.PUT("/banned/:userid", rt.wrap(rt.banUser))

	rt.router.DELETE("/banned/:userid", rt.wrap(rt.unbanUser))

	rt.router.PUT("/following/:userid", rt.wrap(rt.followUser))

	rt.router.DELETE("/following/:userid", rt.wrap(rt.unfollowUser))

	rt.router.GET("/users/:userid/followings", rt.wrap(rt.getFollowingsList))

	rt.router.GET("/users/:userid/followers", rt.wrap(rt.getFollowersList))

	rt.router.GET("/home_page", rt.wrap(rt.getMyStream))

	rt.router.POST("/user/:userid/photos", rt.wrap(rt.uploadPhoto))

	rt.router.GET("/user/:userid/photos/:photoid", rt.wrap(rt.getPhoto))

	rt.router.DELETE("/user/:userid/photos/:photoid", rt.wrap(rt.deletePhoto))

	rt.router.POST("/photos/:photoid/comments/", rt.wrap(rt.commentPhoto))

	rt.router.DELETE("/photos/:photoid/comments/:commentid", rt.wrap(rt.uncommentPhoto))

	rt.router.GET("/photos/:photoid/comments/", rt.wrap(rt.getCommentsList))

	rt.router.GET("/user/:userid/photos", rt.wrap(rt.getPhotosList))

	rt.router.PUT("/like/:photoid", rt.wrap(rt.likePhoto))

	rt.router.DELETE("/like/:photoid", rt.wrap(rt.unlikePhoto))

	///  NON HO IL PATH PER LIKE LIST

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
