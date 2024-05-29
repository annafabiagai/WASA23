<script>
export default {
	data() {
		return {
            errormsg: null,

			photoURL: "",
			liked: false,
            likesList: [],
            commentsList: [],
		}
	},

	props: ['IDphoto','IDuser','nickname','date','likesListParent','commentsListParent','isItMe'], 

	methods: {
		getPhoto() { // mi manca il metodo get photo nell'API, lillo poi l'ha commentato 
			// GET /photos/{pid}/
			this.photoURL = __API_URL__ + `/photos/${this.IDphoto}/`;
		},
		async deletePhoto() {
			try {
				// DELETE /photos/{pid}/
				await this.$axios.delete(`/photos/${this.IDphoto}/`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
				this.$emit("removePhoto", this.IDphoto);
			} catch (error) {
                const status = error.response.status;
        		const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
                alert(this.errormsg);
            }
		},
		visitAuthorProfile() {
            // /profiles/:username
			this.$router.push(`/profiles/${this.nickname}`);
		},
		async likeToggle() {
			try {
				if (!this.liked) {
					// PUT /likes/{pid}
                    await this.$axios.put(`/like/${this.IDphoto}`, null, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
					this.likesList.push({IDuser: localStorage.getItem('token'), nickname: localStorage.getItem('nickname')});
				} else {
					// DELETE /likes/{pid}
                    await this.$axios.delete(`/like/${this.IDphoto}`, {headers: {'Authorization': `${localStorage.getItem('token')}`}});
                    this.likesList = this.likesList.filter(user => user.IDuser != localStorage.getItem('token'));
				}
				this.liked = !this.liked;
			} catch (error) {
                const status = error.response.status;
        		const reason = error.response.data;
                this.errormsg = `Status ${status}: ${reason}`;
                alert(this.errormsg);
            }
    	},
        // on child event
		removeCommentFromList(IDcomment) {
			this.commentsList = this.commentsList.filter(comment => comment.IDcomment != IDcomment);
		},
		addCommentToList(comment){
			this.commentsList.unshift(comment); // at the beginning of the list
		},
        visitLiker(nickname) {
            if (nickname != this.$route.params.username) {
                document.querySelector('.modal-backdrop').remove();
                document.querySelector('.modal').remove();
                document.body.style.overflow = 'auto';
                this.$router.push(`/profiles/${nickname}`);
            }
        }
	},
	async mounted() {
        this.getPhoto()
        // it is a promise
        if (this.likesListParent != null) {
            this.likesList = this.likesListParent
        }
        if (this.commentsListParent != null) {
            this.commentsList = this.commentsListParent
        }
		this.liked = this.likesList.some(user => user.IDuser == localStorage.getItem('token'));
	},
}
</script>

<template>
	<div class="container-fluid mt-3 mb-5 ">

        <UsersModal
        :modalID="'likesModal'+IDphoto" 
		:usersList="likesList"
        @visitUser="visitLiker"
        />

        <CommentModal
        :modalID="'commentModal'+IDphoto" 
		:commentsList="commentsList" 
		:isItMe="isItMe" 
		:IDphoto="IDphoto"
		@removeComment="removeCommentFromList"
		@addComment="addCommentToList"
		/>

        <div class="d-flex flex-row justify-content-center">
            <div class="card my-card">
                <div class="d-flex justify-content-end">
                    <button v-if="isItMe" class="my-trnsp-btn my-dlt-btn me-2" @click="deletePhoto">
						<!--trash bin-->
						<i class="fa-solid fa-trash w-100 h-100"></i>
					</button>
                </div>
                <!--photo-->
                <div class="d-flex justify-content-center photo-background-color">
                    <img :src="photoURL" class="card-img-top img-fluid">
                </div>
                <div class="card-body">
                    <div class="container">
                        <div class="d-flex flex-row justify-content-end align-items-center mb-2">
                            <!--author-->
							<button class="my-trnsp-btn m-0 p-1 me-auto" @click="visitAuthorProfile">
                            	<i> From {{this.nickname}}</i>
							</button>
                            <!--like-->
                            <button class="my-trnsp-btn m-0 p-1 d-flex justify-content-center align-items-center">
                                <i @click="likeToggle" :class="'me-1 my-heart-color w-100 h-100 fa '+(liked ? 'fa-heart' : 'fa-heart-o')"></i>
                                <i data-bs-toggle="modal" :data-bs-target="'#likesModal'+IDphoto" class="my-comment-color ">
                                    {{likesList.length}}
                                </i>
                            </button>
                            <!--comment-->
                            <button class="my-trnsp-btn m-0 p-1  d-flex justify-content-center align-items-center" 
							data-bs-toggle="modal" :data-bs-target="'#commentModal'+IDphoto">
                                <i class="my-comment-color fa-regular fa-comment me-1"></i>
                                <i class="my-comment-color-2"> {{commentsList.length}}</i>
                            </button>
                        </div>
                        <div class="d-flex flex-row justify-content-start align-items-center ">
                            <p> Uploaded on {{date}}</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style>
.photo-background-color{
	background-color: grey;
}
.my-card{
	width: 27rem;
	border-color: black;
	border-width: thin;
}
.my-heart-color{
	color: grey;
}
.my-heart-color:hover{
	color: red;
}
.my-comment-color {
	color: grey;
}
.my-comment-color:hover{
	color: black;
}
.my-comment-color-2{
	color:grey
}
.my-dlt-btn{
	font-size: 19px;
}
.my-dlt-btn:hover{
	font-size: 19px;
	color: var(--color-red-danger);
}
</style>
