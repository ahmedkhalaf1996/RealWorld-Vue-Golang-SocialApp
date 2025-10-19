<template>
    <q-page class="constrain q-pa-md">
        <div class="row q-col-gutter-lg">
            <div class="col-2 col-xs-none"></div>
            <div class="col-8 col-xs-12 text-center">
                <div class="q-pa-md">
                    <q-btn-toggle v-model="model"
                        toggle-color="primary"
                        :options="[
                            {label: 'Posts', value: 'Posts'},
                            {label: 'Users', value: 'Users'},

                        ]"
                    />
                </div>

                <div class="q-pa-md">
                    <q-list bordered>

                        <q-separator />

                        <div v-if="model == 'Users'" >
                            <div v-for="data in Users" :key="data._id">
                                <q-item clickable v-ripple @click="GoUser(data._id)" >
                                    <q-item-section avatar>
                                        <q-avatar>
                                            <img v-if="data?.imageUrl" :src="data?.imageUrl" />
                                            <img v-else src="https://cdn-icons-png.flaticon.com/512/3237/3237472.png" />
                                        </q-avatar>
                                    </q-item-section>
                                    <q-item-section>{{ data?.name }}</q-item-section>
                                </q-item>
                            

                            </div>
                        </div>

                        <div v-if="model == 'Posts'">
                            <div v-for="post in Posts" :key="post._id">
                                <q-item clickable v-ripple @click="GoPost(post._id)">
                                    <q-item-section thumbnail>
                                        <img :src="post.selectedFile">
                                    </q-item-section>
                                    <q-item-section>{{ post?.title }}</q-item-section>
                                    <q-item-section>{{ post?.message }}</q-item-section>
                                
                                </q-item>
                            </div>
                        </div>




                    </q-list>

                </div>

            </div>
            <div class="col-2  col-xs-none"></div>
        </div>
    </q-page>
    
</template>

<script>
import { mapActions } from 'vuex';
export default {
    name:'SearchCompoenent',
    data(){
        return {
            value:true,
            model:'Posts',
            Users:[],
            Posts:[]
        }
    },
    watch:{
        $route(){
            //
            this.GetData()
        }
    },
    methods:{
        ...mapActions(['getPostsUsersBySearch']),
        async GetData(){
            const AllData = await this.getPostsUsersBySearch(String(this.$route.query.search))
            console.log("s", AllData),
            this.Users = AllData.user;
            this.Posts = AllData.posts;
        },
        GoUser(id){
            this.$router.push({path: `/Profile/${id}`})
        },
        GoPost(id){
            this.$router.push({path: `/PostDeatils/${id}`})
        }
    },
    mounted(){
        this.GetData()
    }
}
</script>







