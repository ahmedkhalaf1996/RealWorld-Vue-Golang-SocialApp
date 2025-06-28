<template>
 <div class="row col-12 constrain">
   <div class="col-4 text-center">
    <q-avatar size="150px">
     <img v-if="userData?.imageUrl" :src="userData?.imageUrl" >
     <img v-else src="https://cdn-icons-png.flaticon.com/512/3237/3237472.png">
   </q-avatar>
   </div>
   <div class="col-8 text-left">
    <div class="text-h6 q-pa-lg" style="margin: auto;">
        {{ userData?.name }}
        <q-btn v-if="isSameUser" @click="Edit" flat label="Edit"/>

        <!-- follow un follow -->
         <q-btn v-if="!isSameUser && !isUserFollowing"
            @click="FollowOrUnFollow" flat style="color: #FF0080" label="Follow"/>

         <q-btn v-if="!isSameUser && isUserFollowing"
            @click="FollowOrUnFollow" flat class="primary" label="UN Follow"/>
    </div>
    <q-separator inset />
    <div class="text-subtitle1 q-pa-lg" style="margin: auto;">
        {{ userData.bio }}
        <div>
            <i>{{ userPosts.length }} Posts</i>
            <i>
                <i v-if="userData?.followers?.length > 0">
                    {{ userData?.followers?.length  }}</i>
                    followers
            </i>
            <i>
                <i v-if="userData?.following?.length > 0">
                    {{ userData?.following?.length  }}</i>
                    following
            </i>
        </div>
    </div>
   </div>
 </div>


</template>

<script>
import { mapActions } from 'vuex';
 export default {
    props:['userData','userPosts', 'isSameUser'],
    data(){
        return {isUserFollowing:false}
    }, 
    methods:{
        ...mapActions(['FollowUser', 'GetUserByID']),
        async checkUserFollowing(){
            const logeuid = JSON.parse(localStorage.getItem('profile'))?.result?._id
            // const id = this.userData?._id 

            const { user } = await this.GetUserByID(this.$route.params.id)

            if(user && user?.followers.find((id) => id == logeuid)){
                this.isUserFollowing = true 
            } else {
                this.isUserFollowing = false 
            }
            console.log("isUserFollo", this.isUserFollowing, "isSameUser", this.isSameUser)

        },
        async FollowOrUnFollow(){
            console.log("follow or un Follow user")
            let data = await this.FollowUser(this.userData._id)
            console.log("data show profile follow", data)
            if(data && data.FirstUser){

                this.$emit('update-user', {
                    data: data?.FirstUser
                })
                // change buttom
                this.checkUserFollowing()
            }
        },
        Edit(){
            this.$emit('EditProfile')
        },
    },
    mounted(){
        this.checkUserFollowing()
    }
 }
</script>