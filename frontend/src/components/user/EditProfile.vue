<template>
    <div class="row col-12 constrain">
        <div class="col-4 text-center">
            <q-avatar size="150px">
                <img :src="UserData?.imageUrl">
            </q-avatar>
        </div>
        <div class="col-8 text-left">
            <i>Edit Profile</i>
            <div class="text-h6 q-pa-lg" style="margin: auto;">
                <q-btn v-if="isSameUser" @click="Save" flat label="Save" />
            </div>

            <q-input dense v-model="UserData.name" autofocus placeholder="User Data  Title" />
            <div>
                <q-input 
                    v-model="UserData.bio"
                    placeholder="What's On Your Moind ? Your Bio"
                    type="textarea"
                />
            </div>
            <div class="q-pa-md">
                <q-file 
                    v-model="file"
                    lable="Pick Image"
                    filled 
                />
            </div>


        </div>
    </div>
</template>


<script>
import { mapActions } from 'vuex';
 export default {
    props:['userData', 'isSameUser'],
    data(){
        return {
            file:null,
            UserData:{...this.userData}
        }
    }, 
    watch:{
        file(){
            this.ConvertToBase64()
        }
    },
    methods:{
        ...mapActions(['UpdateUserData']),
        ConvertToBase64(){
            var reader = []
            reader = new FileReader();
            reader.readAsDataURL(this.file);

            new Promise(()=> {
                reader.onload = ()=> {
                    this.UserData.imageUrl = reader.result
                }
            })
        },
        async Save(){
            let userdata = {_id:this.UserData._id, name: this.UserData.name, bio:this.UserData.bio , imageUrl:this.UserData.imageUrl}

            const update = await this.UpdateUserData(userdata)

            if(update && update.data){

                this.$emit('update-user', {
                    data: update?.data
                })

                // console.log('new data', update)
                this.$emit('EditProfile')
            }
        }
      
    }, 
 }
</script>