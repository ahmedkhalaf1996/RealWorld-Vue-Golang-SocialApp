<template>
<q-page class="constrain q-pa-md" >
 <div class="row q-col-gutter-lg constrain">
  <ShowProfile
     :userData="userData"
     :userPosts="userPosts"
     :isSameUser="isSameUser"
     @EditProfile="EditMode = !EditMode"
     @update-user="updateUserLocal"
      v-if="!EditMode" />

 <EditProfile
     :userData="userData"
     :isSameUser="isSameUser"
     @EditProfile="EditMode = !EditMode"
     @update-user="updateUserLocal"
      v-else />
  <div class="col-12">
    <q-separator inset />
  </div>
  <div class="col-4" v-for="post in userPosts" :key="post._id">
    <Post :post="post" />
  </div>
 </div>

</q-page>





</template>
  
  <script>
  // @ is an alias to /src
  import { mapGetters, mapMutations, mapActions } from 'vuex';
  import Post from '@/components/post/Post.vue';
  import ShowProfile from '@/components/user/ShowProfile.vue'
  import EditProfile from '@/components/user/EditProfile.vue';
  export default {
    name: 'ProfileView',
    data(){
      return {
        userPosts:[],
        userData:[],
        isSameUser: false,
        EditMode: false,
      }
    },
    watch:{
      $route(){
        this.GetAll()
      }
    },
    mounted(){
        console.log("userid", this.$route.params.id)
        this.SetData();
        this.GetAll()
    },
    created(){
      this.GetAll()
    },
    computed: {
      ...mapGetters(['GetUserData'])
    },
    methods: {
      ...mapMutations(['SetData']),
      ...mapActions(['GetUserByID']),
      // Get All User data & posts
      async GetAll(){
        const LogedUserID = this.GetUserData()?.result?._id
        console.log('LogedUser', LogedUserID)

        const profileid =  this.$route.params.id

        const data = await this.GetUserByID(profileid)

        this.userData = data?.user 
        this.userPosts = data?.posts

        this.isSameUser = String(LogedUserID) == String(profileid);

        console.log('is same', this.isSameUser)
      },
      updateUserLocal(updatedData){
        this.userData = updatedData.data
      }
    },
    components:{ShowProfile, EditProfile, Post}
  }
  </script>