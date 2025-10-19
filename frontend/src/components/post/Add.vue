<template>
    <q-page-sticky 
      :position="$q.screen.lt.sm ? 'bottom-right': 'bottom-left'" 
      v-show="GetUserData()?.result"
      :offset="$q.screen.lt.sm ? [18, 18] :[18, 60]"
      >
        <div class="q-pa-md q-gutter-sm">
            <q-btn 
             :label="$q.screen.gt.xs ? 'Create Post': ''" 
             style="cursor: pointer;" 
             icon="eva-plus-circle-outline" 
             color="primary" 
             :round="$q.screen.lt.sm"
             :size="$q.screen.lt.sm ? 'lg': 'md'"
             @click="persistent = true"
              />

            <!-- // mobile optimized dialog  -->
            <q-dialog 
             v-model="persistent" 
             persistent 
             transition-show="slide-up" 
             transition-hide="slide-down"
             :maximized="$q.screen.lt.sm"
             :position="$q.screen.lt.sm ? 'bottom': 'standard'"
             >
                <q-card :style="$q.screen.lt.sm ? 'min-height: 80vh; border-radius: 16px 16px 0 0;':'min-width: 350px;' " >
                    <q-card-section class="row items-center q-pb-none" v-if="$q.screen.lt.sm">
                        <div class="text-h6">Create Post</div>
                        <q-space />
                    </q-card-section>

                    <!-- descopt header  -->
                    <q-card-section v-else>
                        <div class="text-h6">Create Post</div>
                    </q-card-section>

                    <q-card-section class="q-pt-none">
                        <q-input 
                            dense 
                            v-model="post.title" 
                            autofocus 
                            placeholder="Post Title" 
                            class="q-pa-md"
                            />
                        <div class="q-pa-md" >
                            <q-input
                                v-model="post.message"
                                placeholder="What's on your mind?"
                                type="textarea"
                                :rows="$q.screen.lt.sm ? 4 : 3"
                                autogrow
                            />
                        </div>
                        <div class="q-pa-md">
                            <q-file                            
                                v-model="file"
                                label="Pick Image"
                                filled 
                                accept="image/*"

                                :style="$q.screen.lt.sm ? 'width: 100%;' : 'max-width: 400px;'" >
                            <template v-slot:prepend>
                                <q-icon name="eva-camera-outline" />
                            </template>

                            </q-file>
                        </div>

                        <div class="q-gutter-sm row items-start" v-if="post.selectedFile">
                            <q-img 
                                :src="post.selectedFile"
                                spinner-coler="red"
                                :style="$q.screen.lt.sm ? 'height: 200px; max-width: 100%;':
                                'height: 140px; max-width: 150px;'"
                                class="rounded-borders"
                            />
                        </div>
                    </q-card-section>

                    <q-card-actions v-if="$q.screen.lt.sm" class="q-pa-md">
                        <q-btn 
                         flat 
                         label="Cancel"
                         v-close-popup
                         class="col"
                         size="md"
                        />
                       <q-btn 
                       unelevated
                         label="Craete"
                         color="primary"
                         v-close-popup
                         class="col"
                         @click="CreatePost"
                         size="md"
                        />
                    </q-card-actions>
                    
                    <q-card-actions v-else align="right" class="text-primary">
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
                    this.resetForm();

                    this.$emit('Created')
                
                        this.$q.notify({
                        icon: 'eva-alert-circle-outline',
                        type: 'positive',
                        message: `Post Created Successfully`
                    })

                }

            }
        },
        resetForm(){
            this.post = {
                title: '',
                message: '',
                name: '',
                selectedFile: null
            }
            this.$nextTick(()=> {
                this.file = null
            })
        },
        ConvertToBase64(){
            if(!this.file || !(this.file instanceof File || this.file instanceof Blob)){
                this.post.selectedFile = null;
                return;
            }
            const reader = new FileReader();
            reader.onload = () => {
                this.post.selectedFile = reader.result;
            }

            reader.onerror = (error)=> {
                console.error('Error reading file', error);
                this.$q.notify({
                    icon: 'evva-alert-circle-outline',
                    type:'negativve',
                    message:'Error reading file'
                });
                this.post.selectedFile = null;
            }

            reader.readAsDataURL(this.file)
        }
    },



}
</script>


