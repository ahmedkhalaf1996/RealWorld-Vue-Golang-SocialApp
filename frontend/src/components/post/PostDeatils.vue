<template>
    <q-page class="constrain q-pa-md">
        <div class="row q-col-gutter-lg">
            <div class="col-3"></div>
            <div class="col-12 col-sm-6" v-if="post ">
                <div class="q-pa-md q-gutter-sm" v-if="!EditPost" >
                    <q-btn v-if="IsSameUser" color="primary" icon="eva-edit" @click="EditPost = !EditPost" label="Edit Post" />
                </div>
                <Post :post="post"
                 :EditPost="EditPost"
                  @changeEdit="EditPost = !EditPost"
                  @postUpdated="handlePostUpdate" />
            </div>
            <div class="col-3"></div>
        </div>
    </q-page>

</template>

<script>
import Post from './Post.vue';
import { mapActions } from 'vuex';
export default {
    name:'PostDeatils',
    data(){
        return {
            EditPost: false,
            post: null,
            IsSameUser:false,
        }
    },
    methods:{
        ...mapActions(['getPost']),
        handlePostUpdate(updatedPost) {
            this.post = updatedPost
        }
    },
    async mounted(){
        const {post} = await this.getPost(this.$route.params.id)
        this.post = post

        const logedInUser = JSON.parse(localStorage.getItem('profile'))
        const LogedInUserId = logedInUser?.result?._id;
        if(post?.creator == LogedInUserId){
            this.IsSameUser = true
        }
    },
    components:{
        Post
    }
}
</script>