
<script>


// get user profile
export default {
    data: function() {
        return {
            errormsg: null,

            // getUserProfile
            username: "",
            photosCount: 0,
            followersCount: 0,
            followingCount: 0,
            isItMe: false,
            doIFollowUser: false,
            isInMyBannedList: false,
            meBanned: false,

            // getPhotosList
            photosList: [],

            // getFollowersList
            followersList: [],

            // getFollowingsList
            followingsList: [],

            userExists: false,
            userID: 0,
        }
    },
    watch: {
        // property to watch
        pathUsername(newUName, oldUName) {
            if (newUName !== oldUName){
                this.getUserProfile()
            }
        }
    },
    computed: {
        pathUsername() {
            return this.$route.params.username
        },
    },
    methods: {
        async getUserProfile() {
            if (this.$route.params.username === undefined) {
                return
            }
            try {
                // /profiles/username
                // getUserId
                // GET /user/{username}
                let username = this.$route.params.username;
                let response = await this.$axios.get(`/user/${username}`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                this.IDuser = response.data;
                // GET /users/{uid}/
                response = await this.$axios.get(`/users/${this.IDuser}/`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                // console.log(response)
                let profile = response.data;
                this.nickname = profile.nickname;
                this.photosCount = profile.photosCount;
                this.followersCount = profile.followersList;
                this.followingCount = profile.followingList;
                this.isItMe = profile.isItMe;
                this.doIFollowUser = profile.doIFollowUser;
                this.isInMyBannedList = profile.isInMyBannedList;
                this.amIBanned = profile.meBanned;
                this.userExists = true;
                if (!this.isInMyBannedList && !this.meBanned) {
                    await this.getPhotosList();
                    this.getFollowersList();
                    this.getFollowingsList();
                }
            } catch (error) {
                const status = error.response.status;
                const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
            }
        },
        async followBtn() {
            try {
                if (this.doIFollowUser) { 
                     // DELETE /following/{uid}
                    await this.$axios.delete(`/following/${this.IDuser}`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                    this.getUserProfile();
                } else {
                    // PUT /following/{uid}
                    await this.$axios.put(`/following/${this.IDuser}`, null, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                    this.getUserProfile();
                }
                this.doIFollowUser = !this.doIFollowUser
            } catch (error) {
                const status = error.response.status;
                const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
            }
        },
        async banBtn() {
            try {
                if (this.isInMyBannedList) {
                    // DELETE /banned/{uid}
                    await this.$axios.delete(`/banned/${this.IDuser}`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                    this.getUserProfile();
                } else {
                    // PUT /banned/{uid}
                    await this.$axios.put(`/banned/${this.IDuser}`, null, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                    this.getUserProfile();
                }
            } catch (error) {
                const status = error.response.status;
                const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
            }
        },
        async uploadPhoto() {
            try {
                let file = document.getElementById('fileUploader').files[0];
                const reader = new FileReader();
                reader.readAsArrayBuffer(file); // stored in result attribute
                reader.onload = async () => {
                    // POST /photos/
                    let response = await this.$axios.post('/photos/', reader.result, {headers: {'Authorization': `${localStorage.getItem('token')}`, 'Content-Type': 'image/*'}});
                    this.photosList.unshift(response.data); // at the beginning of the list
                    this.photosCount += 1;
                }
            } catch (error) {
                const status = error.response.status;
                const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
            }
        },
        async getPhotosList() {
            try {
                // GET /users/{uid}/photos/
                let response = await this.$axios.get(`/users/${this.IDuser}/photos/`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                this.photosList = response.data === null ? [] : response.data;
            } catch (error) {
                const status = error.response.status;
                const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
            }
        },
        async getFollowersList() {
            try {
                // GET /users/{uid}/followers/
                let response = await this.$axios.get(`/users/${this.IDuser}/followers/`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                this.followersList = response.data === null ? [] : response.data;
            } catch (error) {
                const status = error.response.status;
                const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
            }
        },
        async getFollowingsList() {
            try {
                // GET /users/{uid}/followings/
                let response = await this.$axios.get(`/users/${this.IDuser}/followings/`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                this.followingsList = response.data === null ? [] : response.data;
            } catch (error) {
                const status = error.response.status;
                const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
            }
        },
        // on child event
        removePhotoFromList(IDphoto){
            this.photosList = this.photosList.filter(photo => photo.IDphoto != IDphoto);
            this.photosCount -= 1;
        },
        visitUser(nickname) {
            if (nickname != this.$route.params.nickname) {
                this.$router.push(`/profiles/${nickname}`);
            }
        }
    },
    mounted() {
        this.getUserProfile();
    }
}
</script>

<template>

    <UsersModal
    :modalID="'usersModalFollowers'" 
    :usersList="followersList"
    @visitUser="visitUser"
    />

    <UsersModal
    :modalID="'usersModalFollowing'" 
    :usersList="followingsList"
    @visitUser="visitUser"
    />

    <div class="container-fluid" v-if="userExists && !meBanned">
        <div class="row">
            <div class="col-12 d-flex justify-content-center">
                <div class="card w-50 container-fluid">
                    <div class="row">
                        <div class="col">
                            <div class="card-body d-flex justify-content-between align-items-center">
                                <h5 class="card-title p-0 me-auto mt-auto">@{{nickname}}</h5>

                                <button v-if="!isItMe && !isInMyBannedList" @click="followBtn" class="btn btn-success ms-2">
                                    {{doIFollowUser ? "Unfollow" : "Follow"}}
                                </button>

                                <button v-if="!isItMe" @click="banBtn" class="btn btn-danger ms-2">
                                    {{isInMyBannedList ? "Unban" : "Ban"}}
                                </button>
                            </div>
                        </div>
                    </div>

                    <div v-if="!isInMyBannedList" class="row mt-1 mb-1">
                        <button class="col-4 d-flex justify-content-center btn-foll">
                            <h6 class="ms-3 p-0 ">Posts: {{photosCount}}</h6>
                        </button>
                    
                        <button class="col-4 d-flex justify-content-center btn-foll">
                            <h6 data-bs-toggle="modal" :data-bs-target="'#usersModalFollowers'">
                                Followers: {{followersList.length}}
                            </h6>
                        </button>
                    
                        <button class="col-4 d-flex justify-content-center btn-foll">
                            <h6 data-bs-toggle="modal" :data-bs-target="'#usersModalFollowing'">
                                Following: {{followingsList.length}}
                            </h6>
                        </button>
                    </div>
                </div>
            </div>
        </div>


        <div class="row">
            <div class="container-fluid mt-3">
                <div class="row ">
                    <div class="col-12 d-flex justify-content-center">
                        <h2>Posts</h2>
                        <input id="fileUploader" type="file" class="profile-file-upload" @change="uploadPhoto" accept=".jpg, .png">
                        <label v-if="isItMe" class="btn my-btn-add-photo ms-2 d-flex align-items-center" for="fileUploader"> Add </label>
                    </div>
                </div>
                <div class="row ">
                    <div class="col-3"></div>
                    <div class="col-6">
                        <hr class="border border-dark">
                    </div>
                    <div class="col-3"></div>
                </div>
            </div>
        </div>

        <div class="row">
            <div class="col">
                <div v-if="!isInMyBannedList && photosCount>0">
                    <Photo v-for="photo in photosList"
                    :key="photo.IDphoto"
                    :IDphoto="photo.IDphoto"
                    :IDuser="photo.IDuser"
                    :nickname="this.nickname"
                    :date="photo.date"
                    :likesListParent="photo.likeList"
                    :commentsListParent="photo.commentsList"
                    :isItMe="isItMe"
                    @removePhoto="removePhotoFromList"
                    />
                </div>
                
                <div v-if="!isInMyBannedList && photosCount==0" class="mt-5 ">
                    <h2 class="d-flex justify-content-center" style="color: white;">No posts yet</h2>
                </div>
            </div>
        </div>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    
</template>

<style>
.profile-file-upload{
    display: none;
}
.my-btn-add-photo{
    background-color: green;
    border-color: grey;
}
.my-btn-add-photo:hover{
    color: white;
    background-color: green;
    border-color: grey;
}
.btn-foll{
    background-color: transparent;
    border: none;
    padding: 5px;
}
</style>

