<template>
 <q-header class="bg-white text-grey-10" borderd>
  <q-toolbar class="constrain x">
    <q-btn flat to="/">
        <q-icon left size="3em" name="eva-camera-outline" />
        <q-toolbar-title class="text grand-hotel text-bold"> Home</q-toolbar-title>
    </q-btn>

    <q-separator class="large-screen-only" vertical spaced />

    <q-toolbar-title class="text-center">
        <q-input bottom-slots class="nuks" label="search" @keyup.enter="GoSearch($event)">
        </q-input>
    </q-toolbar-title>

    <q-btn round 
    v-show="GetUserData()?.result"
    @click="GoToChat"
    :icon="unReadedMessages > 0 ? 'eva-message-square-outline': 'eva-message-square'"
    :color="unReadedMessages > 0 ? 'primary':'dark'"
    >
     <q-badge v-if="unReadedMessages >0" color="negative" floating  rounded :label="unReadedMessages"/>
    </q-btn>

   
    <q-btn round 
     v-show="GetUserData()?.result"
     @click="GoToNotification"
     :icon="notificationNum > 0 ? 'eva-bell-outline':'eva-bell'"
     :color="notificationNum > 0 ? 'primary': 'dark'"
    >
     <q-badge v-if="notificationNum > 0 " floating color="negative" rounded :label="notificationNum"/>
    </q-btn>
 

    <q-btn v-show="GetUserData()?.result" round>
     <q-avatar size="42px" v-if="GetUserData()?.result?.imageUrl">
        <img :src="GetUserData()?.result?.imageUrl" >
     </q-avatar>
     <q-avatar size="42px" v-else>
        <img src="https://cdn-icons-png.flaticon.com/512/3237/3237472.png" >
     </q-avatar>
     <q-menu>
        <q-list style="min-width: 100px">
            <q-item clickable v-close-popup>
                <q-item-section @click="Profile" >Profile</q-item-section>
            </q-item>
            <q-separator />
            <q-item clickable v-close-popup>
                <q-item-section @click="LogUserOut">Logout</q-item-section>
            </q-item>
        </q-list>
     </q-menu>
    </q-btn>

  </q-toolbar>
 </q-header>


</template>
  
  <script>
  import { mapGetters, mapMutations, mapActions } from 'vuex';
  export default {
    name: 'NavBar',
    data (){
        return {
          notificationNum:0,
          unReadedMessages:0,
          // userData: null,
        }
    },
    computed: {
      ...mapGetters(['GetUserData'])
    },
    methods: {
      ...mapMutations(['SetData']),
      ...mapActions(['logout','GetUnReadedNotifyNum', 'GetUnreadedMessageNum']),
      GoSearch(e) {
            console.log("go", e.target.value)
            this.$router.push({path: `/Search`, query: { search: e.target.value }})
        },
      Profile() {
        let id = this.GetUserData()?.result?._id;
        this.$router.push(`/Profile/${id}`)
      },
      LogUserOut(){
        this.logout(),
        this.$router.push(`/Auth`)
      },
      GoToNotification(){
        this.$router.push('/Notification')
      },
      GoToChat(){
        this.$router.push('/Chat')
      }
    },
    async mounted(){
    this.SetData();
    // getnot number
        this.NotifyList = await this.GetUnReadedNotifyNum(this.GetUserData().result._id)
        let numofunreadednot = 0;
        this.NotifyList.forEach(el =>{
             if (! el.isreded ) {
                        numofunreadednot++;
             }
            })
        this.notificationNum = numofunreadednot;
    // get chat messages numbers
    const  {total} = await this.GetUnreadedMessageNum(this.GetUserData().result._id)
    this.unReadedMessages = total;
  },
  }
  </script>

<style lang="sass">
.nuks 
  width: 250px
  text-align: center
  display: inline-block !important

.q-toolbar-title
  display: flex 
  align-items: center 

.q-btn 
  margin-left: 10px
</style>