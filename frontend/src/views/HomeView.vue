<template>
   <q-page class="constrain q-pa-md">
     <div class=" row q-col-gutter-lg">
      <div class="col-3">
        <SideBar />
      </div>
      <div v-if="!load" class="col-6  q-mx-auto">
        <div class="q-pa-md">
          <q-card>
            <q-item>
              <q-item-section avatar>
                <q-skeleton type="QAvatar"/>
              </q-item-section>
              <q-item-section>
                <q-item-label>
                  <q-skeleton type="text" />
                </q-item-label>
                <q-item-label caption>
                  <q-skeleton type="text" />
                </q-item-label>
              </q-item-section>
            </q-item>

            <q-skeleton height="200px" square />
            <q-card-actions class="q-gutter-md">
              <q-skeleton type="QBtn"/>
              <q-skeleton type="QBtn"/>
            </q-card-actions>
          </q-card>
        </div>
      </div>
      <div v-else class="col-6 q-mx-auto">
        <Post v-for="post in posts" :key="post._id" :post="post" />
         <!-- add posts here compnont  -->
      </div>
      <div class="col-3">
        <Rightbar />
      </div> 
    </div>
     <div class="q-pa-lg flex justify-center fixed-bottom">
      <Add @Created="GetAllPosts"/>
      <q-pagination 
        v-model="current"
        color="primary"
        :max="max"
        :max-pages="5"
        :ellipses="false"
        :boundary-numbers="false"
        />
     </div>
   </q-page>
</template>

<script>
import Add from '@/components/post/Add.vue'
import Post from '@/components/post/Post.vue'
import SideBar from '@/components/sideBar/SideBar.vue';
import Rightbar from '@/components/rightbar/Rightbar.vue';
import { mapActions } from 'vuex';
export default {
  name: 'HomeView',
  data(){
    return {
      current:1,
      max:0,
      posts:[],
      load:false
    }
  },
  watch:{
    current(){
      this.GetAllPosts();
    }
  },
  components: {
    Add,
    Post,
    SideBar,
    Rightbar,
  },
  methods:{
    ...mapActions(['getPosts']),
    async GetAllPosts(){
      console.log("Get All Posts Called")
      const data = await this.getPosts(this.current)
      console.log("post data", data)
      if(data?.data){
       this.max = data?.numberOfPages;
       this.posts = data?.data; 
      }

      if(data){
        this.load = true;
      }
    }
  },
  async mounted(){
    setTimeout(() => {
      this.GetAllPosts()
    }, 5000);
  },
}
</script>
