<template>
<q-page class="constrain q-pa-md" >
 <div class="row q-col-gutter-lg constrain">
  <!-- error State  -->
  <div v-if="error" class="col-12 q-pa-lg text-center text-negative">
    <q-icon name="eva-aler-circle-outline" size="48px" />
    <q-btn
    @click="resetAndLoadProfile"
    color="primary"
    outline
    class="q-mt-md"
    >
    Try Again 
    </q-btn>
  </div>
 
  <!-- profile content  -->
   <template v-else>
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
      v-else-if="EditMode && userData" 
      />
  <div class="col-12" v-if="userData">
    <q-separator inset />
  </div>

  <!-- loading skelton f rinital posts  -->
  <div v-if="!load && !error" class="col-12">
    <div class="col-12 col-sm-6 col-md-4" v-for="i in 3" :key="i">
      <q-card>
        <q-skeleton height="200px" square />
        <q-card-action>
          <q-skeleton type="text" class="text-h6" />
          <q-skeleton type="text" class="text-subtitle2" width="50%" />
          
        </q-card-action>
      </q-card>
    </div>
  </div>

 
  <!-- posts grid  -->
  <div
   v-else-if="load && userPosts.length> 0"
   class="col-12 col-sm-6 col-md-4"
   v-for="post in userPosts"
   :key="post._id"
  >
    <Post :post="post" />
  </div>

  <!-- loading indicator for more posts  -->
  <div v-if="loadingMore" class="col-12 q-pa-lg text-center">
    <q-spinner-hourglass color="primary" size="3em" />
    <div class="q-mt-md text-grey-7">
      Loading more Posts...
    </div>
  </div>

  <!-- end of posts indicator  -->
   <div v-if="hasReachedend && userPosts.length > 0 && load" class="col-12 q-pa-md text-center text-grey-6">
    <q-icon name="eva-inbox-outline" size="24px" />
    <div class="q-mt-sm">No More Posts</div>
   </div>

   <!-- No Posts Message  -->
   <div v-if="load & userPosts.length === 0 && userData && !error" class="col-12 q-pa-lg text-center text-grey-6">
        <q-icon name="eva-inbox-outline" size="48px" />
    <div class="q-mt-md text-h6">No Posts Yet</div>
    <div class="text-body2">
      {{  isSameUser ? "You have't": "This user hasn't"  }} share any posts.
    </div>
   </div>




   </template>
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
        userData:[],
        isSameUser: false,
        EditMode: false,
        load: false,
        loadingMore: false,
        currentPage:1,
        maxPages:0,
        hasReachedend:false,
        error:null,
        scrollListener: null,
      }
    },
    watch:{
      $route: {
        handler(newRoute, oldRoute){
          if(newRoute.params.id !== oldRoute?.params.id){
            this.resetAndLoadProfile()
          }
        },
        immediate:true
      }
    },
    mounted(){
        console.log("userid", this.$route.params.id)
        this.SetData();
        this.resetAndLoadProfile()
    },
    computed: {
      ...mapGetters(['GetUserData','GetUserPosts']),
      userPosts(){
        console.log("USr POST", this.GetUserPosts)
        return this.GetUserPosts || [];
      }
    },
    methods: {
      ...mapMutations(['SetData']),
      ...mapActions(['GetUserByID', 'ResetUserPosts']),
      // Reset and load profile data 
      async resetAndLoadProfile(){
        try {
          this.error = null;
          this.ResetUserPosts();
          this.load = false;
          this.loadingMore = false;
          history.currentPage = 1;
          this.hasReachedend = false;
          this.userData = null;

          this.removeScrollListener();

          // Load inital data
          await this.GetAll(false);

          //  add scorll listern after sucessful laod
          if(!this.error){
            this.addScrollListener();
          }

        } catch (error) {
          console.error('Faild to reset and load profile', error);
          this.error = "Faild to Load profile . Please Try again."
        }
      },
      // Get All User data & posts
      async GetAll(append = false){
        try {
          const LogedUserID = this.GetUserData()?.result?._id
          const profileid =  this.$route.params.id

          if(!profileid || profileid === 'undefined' || profileid.trim() === ''){
            throw new Error('Pifile Id is Missing or invalid')
          }

          const data = await this.GetUserByID({
            id:profileid,
            page: this.currentPage,
            append: append
          })

          if(!append && data?.user){
            this.userData = data.user;
            this.isSameUser = String(LogedUserID) === String(profileid)
            console.log("Is same user:", this.isSameUser)
          }

          this.maxPages = data?.numberOfPages || 0;
          this.hasReachedend = this.currentPage >= this.maxPages;

          this.load = true;
          this.error = null;
          // console.log("User Profile Check Redis ", data)
          return data;



        } catch (error) {
          console.error("Error in GetAll:", error);
          this.error = 'Server Is currentlly unavailable.'
          this.load = true;
          throw error;
        }
      },
      async loadMorePosts(){
        if(this.loadingMore || this.hasReachedend || this.error){
          return;
        }
        if(this.currentPage < this.maxPages){
          this.loadingMore = true;
          this.currentPage++;

          try {
            await this.GetAll(true);
            await new Promise(resolve => setTimeout(resolve, 500));
          } catch (error) {
            console.error('Error loading more posts:', error);
            this.currentPage--;
          } finally {
            this.loadingMore = false;
          }
        }
      },
      handleScroll(){
        const scrllTop = window.pageXOffset || document.documentElement.scrollTop;
        const windowHeight = window.innerHeight;
        const documentHeigt = document.documentElement.scrollHeight;

        if(scrllTop + windowHeight >= documentHeigt - 200) {
          this.loadMorePosts();
        }
      },
      addScrollListener(){
        if(!this.scrollListener){
          this.scrollListener = this.handleScroll.bind(this);
          window.addEventListener('scroll', this.scrollListener)
        }
      },
      removeScrollListener(){
        if(this.scrollListener){
          window.removeEventListener('scroll', this.scrollListener);
          this.scrollListener = null;
        }
      },
      updateUserLocal(updatedData){
        if(updatedData?.data){
           this.userData = updatedData.data
        }
      }
    },
    beforeUnmount(){
      this.removeScrollListener();
    },
    components:{ShowProfile, EditProfile, Post}
  }
  </script>

<style scoped >
.p-page {
  scroll-behavior: smooth;
}

.text-negative {
  color: #c10015;
}
</style>
