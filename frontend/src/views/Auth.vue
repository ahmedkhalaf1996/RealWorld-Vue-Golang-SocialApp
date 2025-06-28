<template>
 <q-page class="constrain q-pa-md">
    <div class="row q-col-gutter-lg">
        <div class="col-5">
            <q-card class="my-card" style="width: 100%; padding: 10px;">
                <h1 class="text-h6 text-center">Signin</h1>
                <q-card-section>
                    <form @submit.prevent.stop="Login" class="q-gutter-md">
                        <q-input
                          filled 
                           v-model="Sin_data.email"
                           label="Your Email *"
                           hint="Email"
                           lazy-rules
                        />
                        <q-input
                           filled 
                           v-model="Sin_data.password"
                           label="Your Password *"
                           hint="password"
                           type="password"
                           lazy-rules
                        />  
                        <div>
                            <q-btn label="sigin in" type="submit" color="primary" />
                        </div>                      
                    </form>
                </q-card-section>
            </q-card>
        </div>
        <div class="col-7">
            <q-card class="my-card" style="width: 100%; padding: 10px;">
                <h1 class="text-h6 text-center">Signup | Craete New Account</h1>
                <q-card-section>
                    <form @submit.prevent.stop="Register" class="q-gutter-md">
                        <q-input
                          filled 
                           v-model="Sup_data.firstName"
                           label="Your first Name *"
                           hint="firstName"
                           lazy-rules
                        />
                        <q-input
                          filled 
                           v-model="Sup_data.lastName"
                           label="Your lastName *"
                           hint="lastName"
                           lazy-rules
                        />
                        <q-input
                          filled 
                           v-model="Sup_data.email"
                           label="Your Email *"
                           hint="Email"
                           lazy-rules
                        />
                        <q-input
                           filled 
                           v-model="Sup_data.password"
                           type="password"
                           label="Your Password *"
                           hint="password"
                           lazy-rules
                        />  
                        <div>
                            <q-btn label="Create New Account" type="submit" color="positive" />
                        </div>                      
                    </form>
                </q-card-section>
            </q-card>
        </div>
    </div>
 </q-page>
</template>
  
<script>

import {mapActions} from 'vuex'

export default {
  name: 'AuthView',
  data () {
    return {
        Sin_data:{
            email:'',
            password: '',
        },
        Sup_data:{
            email:'',
            password: '',
            firstName: '',
            lastName:'',
        }
    }
  },
  methods:{
    ...mapActions(['signin', 'signup']),
    async Login(){
        console.log("login in data", this.Sin_data)
        var validate = true 
        if (this.Sin_data.email == ''){
            validate = false 
            this.$q.notify({
                icon:'eva-alert-circle-outline',
                type:'negative',
                message:`Email is Required`
           })    
        } else if (this.Sin_data.password == ''){
            validate = false 
            this.$q.notify({
                icon:'eva-alert-circle-outline',
                type:'negative',
                message:`password is Required`
           })  
        }
        // sucess and ready to next
        if(validate){
            var formdata = {email:this.Sin_data.email, password: this.Sin_data.password};
            const data  = await this.signin(formdata);
            // console.log("data", data, 'message', message)
            // console.log("response data on Login", data)
            // console.log('data response', data.response.data.message )

            if(data?.response?.data?.message || data?.response?.data){
                this.$q.notify({
                icon:'eva-alert-circle-outline',
                type:'negative',
                message:`Erorr ${data.response.data.message ?  data?.response?.data.message :  data?.response?.data}`
           })  
          } else {
            this.$q.notify({
                icon:'eva-alert-circle-outline',
                type:'positive',
                message:`Successfully Sigin in`
           })  
           this.$router.push('/')
          }
        }

    },
    async Register(){
        console.log("Register in data", this.Sup_data)
        // validation
        var isVaidate = true;
        for (const key in this.Sup_data){
            const val = this.Sup_data[key];
            if(val === ''){
                this.$q.notify({
                icon:'eva-alert-circle-outline',
                type:'negative',
                message:`${key} is Required`
              });
              isVaidate = false
            }
        } // v end
        if(isVaidate){
            const data = await this.signup(this.Sup_data)

            console.log("data on Register", data)
            if(data?.response?.data?.message){
                this.$q.notify({
                icon:'eva-alert-circle-outline',
                type:'negative',
                message:`Erorr ${data.response.data.message}`
           })  
          } else {
            // meaning succusfully
            this.$q.notify({
                icon:'eva-alert-circle-outline',
                type:'positive',
                message:`Successfully Sigin up`
           })  
        //    this.$router.push('/')
          }
        }
    },
  }
}
</script>