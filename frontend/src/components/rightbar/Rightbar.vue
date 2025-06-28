<template>
<div>
    <q-item class="RightBar"
        v-for="user in UsersData"
        :key="user._id"
        :to="`/Profile/${user?._id}`" >
    
    <q-item-section avatar>
        <q-avatar>
            <img v-if="user.imageUrl" :src="user.imageUrl" >
            <img v-else src="https://cdn-icons-png.flaticon.com/512/1077/1077063.png" >
            
        </q-avatar>
    </q-item-section>
    <q-item-section>
        <q-item-label class="text-bold">{{ user?.name }}</q-item-label>
        <q-item-label caption>
            {{ user?.bio }}
        </q-item-label>
    </q-item-section>
    </q-item>
</div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
export default{
    name:'RightBar',
    data(){
        return {
            UsersData:[]
        }
    },
    computed:{
        ...mapGetters(['GetUserData'])
    },
    methods:{
        ...mapActions(['GetTheUserSug'])
    },
    async mounted(){
        let logedinUser = this.GetUserData()?.result 
        if (logedinUser){
            const {users} = await this.GetTheUserSug(logedinUser?._id)
            console.log("Rightbar", users)
            this.UsersData = users
        }
    }
}
</script>

<style lang="sass" scoped>
.RightBar 
    cursor: pointer
</style>



