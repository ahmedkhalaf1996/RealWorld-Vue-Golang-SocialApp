<template>
    <q-page-sticky position="bottom-left"  v-show="GetUserData()?.result">
        <div class="q-pa-md q-gutter-sm">
            <q-btn label="Create Post" style="cursor: pointer;" icon="eva-plus-circle-outline" 
                    color="primary" @click="persistent = true" />

            <!-- // pupup  -->
            <q-dialog v-model="persistent" persistent transition-show="scale" transition-hide="scale">
                <q-card style="min-width: 350px;" >
                    <q-card-section>
                        <div class="text-h6">Create Post</div>
                    </q-card-section>

                    <q-card-section class="q-pt-none">
                        <q-input dense v-model="post.title" autofocus placeholder="Post Title"/>
                        <div class="q-pa-md" style="max-width: 300px;">
                            <q-input
                                v-model="post.message"
                                placeholder="What's on your mind?"
                                type="textarea"
                            />
                        </div>
                        <div class="q-pa-md">
                            <q-file 
                                v-model="file"
                                label="Pick Image"
                                filled 
                                style="max-width: 400px;"
                            />

                        </div>

                        <div class="q-gutter-sm row items-start">
                            <q-img 
                                :src="post.selectedFile"
                                spinner-coler="red"
                                style="height: 140px; max-width: 150px;"
                            />
                        </div>
                    </q-card-section>

                    <q-card-actions align="right" class="text-primary">
                        <q-btn flat label="Cancel" v-close-popup />
                        <q-btn flat label="Create" v-close-popup @click="CreatePost"/>
                        
                    </q-card-actions>
                </q-card>
            </q-dialog>
        </div>
    </q-page-sticky>  

</template>


<script>

import {mapActions, mapGetters} from 'vuex'

export default {
    name: 'AddComponent',
    data (){
      return {
        persistent: false,
        post: {title:'', message:'', name:'', selectedFile: null},
        file: null
      }        
    },
    watch:{
        file(){
            // convert fun
            this.ConvertToBase64()
        }
    },
    computed: {
        ...mapGetters(['GetUserData'])
    },
    methods: {
        ...mapActions(['createPost']),
        async CreatePost(){
            var name = JSON.parse(localStorage.getItem('profile'))?.result?.name;
            this.post.name = name;
            // validation
            var isValidate = true;
            for (const key in this.post){
                const val = this.post[key];
                if (val === ''){
                    this.$q.notify({
                        icon: 'eva-alert-circle-outline',
                        type: 'negative',
                        message: `${key} is Required`
                    })
                    isValidate = false
                }
            }
            // after validate
            if(isValidate){
                const data = await this.createPost(this.post);
                console.log('data', data)

                if (data) {
                    // console.log('data', data)
                    this.$emit('Created')
                
                        this.$q.notify({
                        icon: 'eva-alert-circle-outline',
                        type: 'positive',
                        message: `Post Created Successfully`
                    })

                }

            }
        },
        ConvertToBase64(){
            var reader = [];
            reader = new FileReader();
            reader.readAsDataURL(this.file);

            new Promise(()=> {
                reader.onload = ()=> {
                    this.post.selectedFile = reader.result
                }
            })
        }
    },



}
</script>


