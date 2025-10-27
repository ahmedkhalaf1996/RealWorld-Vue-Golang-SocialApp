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
            @click="FollowOrUnFollow" 
            flat
            :loading="followLoading"
            :disable="followLoading"
            style="color: #FF0080" label="Follow"/>

         <q-btn v-if="!isSameUser && isUserFollowing"
            @click="FollowOrUnFollow" flat
                       :loading="followLoading"
            :disable="followLoading"
            class="primary" label="UN Follow"/>
    </div>
    <q-separator inset />
    <div class="text-subtitle1 q-pa-lg" style="margin: auto;">
        {{ userData?.bio }}
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
import { mapActions, mapGetters } from 'vuex';
 export default {
    props:['userData','userPosts', 'isSameUser'],
    data(){
        return {isUserFollowing:false, 
            followLoading:false
        }
    }, 
    watch:{
        userData:{
            handler(newUserData){
                if(newUserData && !this.isSameUser){
                    this.checkUserFollowing();
                }
            },
            immediate:true
        },
        isSameUser:{
            handler(newValue){
                if(newValue){
                    this.isUserFollowing = false
                }
            },
            immediate: true
        }
    },
    computed:{
        ...mapGetters(['GetUserData'])
    },
    methods:{
        ...mapActions(['FollowUser']),
        async checkUserFollowing(){
            if(!this.userData || this.isSameUser){
                this.isUserFollowing = false;
                return;
            }

            const logeuid = this.GetUserData()?.result?._id

            if(!logeuid){
                this.isUserFollowing = false;
                return;
            }

            const followers = this.userData.followers || [];
            this.isUserFollowing = followers.some(followerId =>
                String(followerId) === String(logeuid)
            )

        },
        async FollowOrUnFollow(){
            if(this.isSameUser || this.followLoading || !this.userData?._id){
                return
            }

            try {

                this.followLoading = true;

                let data = await this.FollowUser(this.userData._id)
                if(data){
                    if(data.FirstUser){
                        this.$emit('update-user', {
                            data: data.FirstUser
                        })
                    }

                    this.isUserFollowing = !this.isUserFollowing;

                    this.$q.notify({
                        type:'positive',
                        message:this.isUserFollowing ? 'Now Following' : 'Unfollowed',
                        timeout: 1500
                    })
                }
            } catch (error) {
                console.error('Follwing error', error)
                 this.$q.notify({
                        type:'negative',
                        message:'Faild to update follow status',
                        timeout: 1500
                    })
            } finally {
                this.followLoading = false;
            }
        },
        Edit(){
            this.$emit('EditProfile')
        },
    },
    mounted(){
        if(this.userData && !this.isSameUser){
                    this.checkUserFollowing()
        }
    }
 }
</script>