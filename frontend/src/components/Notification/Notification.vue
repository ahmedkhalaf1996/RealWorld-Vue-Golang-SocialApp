<template>
 <q-page class="constrain q-pa-md">
    <div class="row q-col-gutter-lg">
        <div class="col-12">

            <q-list bordered padding >

                <div v-for="notify in NotifyList" :key="notify._id">
                    <q-item clickable @click="MoveToThePath(notify)" :class="{'text-red': !notify.isreded}" >
                        <q-item-section top avatar>
                            <q-avatar v-if="notify?.user?.avatart">
                                <img :src="notify?.user?.avatart">
                            </q-avatar>
                            <q-avatar v-else>
                                <img src="https://cdn-icons-png.flaticon.com/512/3237/3237472.png">
                            </q-avatar>
                            
                        </q-item-section>

                        <q-item-section>
                            <q-item-label>{{ notify?.deatils }}</q-item-label>
                            <q-item-label>{{ notify?.user?.name }}</q-item-label>
                        </q-item-section>
                    </q-item>
                    <q-separator spaced />
                </div>


            </q-list>

        </div>
    </div>
 </q-page>


</template>

<script>
import { mapGetters, mapActions } from 'vuex';
// import {watch} from 'vue';

export default {
    name:'Notification-Component',
    data(){
        return {
            NotifyList:[]
        }
    },
    async mounted(){
        var id = this.GetUserData().result._id;
        this.NotifyList = await this.GetUnReadedNotifyNum(id)
        console.log("notifilist", this.NotifyList)
        // mark notification as readed
        setTimeout(() => {
            this.NotifyList.forEach(async el =>{
                if(!el.isreded){
                    await this.MarkNotifyAsReaded(id);
                    el.isreded = true;
                }
            })
        }, 500);
    },
    computed:{
        ...mapGetters(['GetUserData'])
    },
    methods:{
        ...mapActions(['GetUnReadedNotifyNum', 'MarkNotifyAsReaded']),

        MoveToThePath(notify){
            if(notify?.deatils.toString().includes("Post")){
                this.$router.push(`/PostDeatils/${notify.targetid}`);
            } else {
                this.$router.push(`/Profile/${notify.targetid}`);
            }
        }
    }

}

</script>