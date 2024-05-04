package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	rt.router.POST("/session", rt.wrap(rt.doLogin)) // OKAY

	rt.router.PUT("/user", rt.wrap(rt.setMyUserName)) // OKAY

	rt.router.GET("/users/:userid", rt.wrap(rt.getUserProfile)) //OKAY

	rt.router.GET("/users/", rt.wrap(rt.searchNickname))

	rt.router.GET("/user/:nickname", rt.wrap(rt.getIDuser))

	rt.router.PUT("/banned/:userid", rt.wrap(rt.banUser)) //OKAY

	rt.router.DELETE("/banned/:userid", rt.wrap(rt.unbanUser)) //OKAY

	rt.router.PUT("/following/:userid", rt.wrap(rt.followUser)) //OKAY

	rt.router.DELETE("/following/:userid", rt.wrap(rt.unfollowUser)) //OKAY

	rt.router.GET("/users/:userid/followings", rt.wrap(rt.getFollowingsList))

	rt.router.GET("/users/:userid/followers", rt.wrap(rt.getFollowersList)) //OKAY

	rt.router.GET("/home_page", rt.wrap(rt.getMyStream))

	rt.router.POST("/photos", rt.wrap(rt.uploadPhoto)) //OKAY

	rt.router.DELETE("/photos/:photoid", rt.wrap(rt.deletePhoto)) //OKAY

	rt.router.POST("/photos/:photoid/comments/", rt.wrap(rt.commentPhoto)) //OKAY

	rt.router.DELETE("/photos/:photoid/comments/:commentid", rt.wrap(rt.uncommentPhoto)) //OKAY

	rt.router.GET("/users/:userid/photos/", rt.wrap(rt.getPhotosList)) //OKAY

	rt.router.PUT("/like/:photoid", rt.wrap(rt.likePhoto)) //OKAY

	rt.router.DELETE("/like/:photoid", rt.wrap(rt.unlikePhoto)) //OKAY

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
