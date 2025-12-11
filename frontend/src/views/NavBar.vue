<template>
 <q-header class="bg-white text-grey-10" bordered>
  <q-toolbar class="constrain">
    <q-btn flat to="/" class="gt-xs">
        <q-icon left size="3em" name="eva-camera-outline" />
        <q-toolbar-title class="text grand-hotel text-bold"> Home</q-toolbar-title>
    </q-btn>

    <q-separator class="large-screen-only gt-xs" vertical spaced />

    <q-toolbar-title class="text-center gt-xs">
        <q-input bottom-slots class="nuks" label="search" @keyup.enter="GoSearch($event)">
        </q-input>
    </q-toolbar-title>

    <q-btn round 
    v-show="GetUserData()?.result"
    @click="GoToChat"
    :icon="unReadedMessages > 0 ? 'eva-message-square-outline': 'eva-message-square'"
    :color="unReadedMessages > 0 ? 'primary':'dark'"
    class="gt-xs"
    >
     <q-badge v-if="unReadedMessages >0" color="negative" floating  rounded :label="unReadedMessages"/>
    </q-btn>

   
    <q-btn round 
     v-show="GetUserData()?.result"
     @click="GoToNotification"
     :icon="notificationNum > 0 ? 'eva-bell-outline':'eva-bell'"
     :color="notificationNum > 0 ? 'primary': 'dark'"
     class="gt-xs"
    >
     <q-badge v-if="notificationNum > 0 " floating color="negative" rounded :label="notificationNum"/>
    </q-btn>
 

    <q-btn v-show="GetUserData()?.result" round class="gt-xs">
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
    <!-- mobile menue buttom  -->
    <q-btn 
      flat
      round
      icon="eva-menu-outline"
      @click="mobileMenuOpen = !mobileMenuOpen"
      class="lt-sm"
    />
    <q-btn flat to="/" class="lt-sm" >
      <q-icon left size="2em" name="eva-camera-outline" />
      <q-toolbar-title class="text grand-hotel text-bold">Home</q-toolbar-title>
    </q-btn>

    <q-space class="lt-sm" />

    <q-btn v-show="GetUserData()?.result" round class="lt-sm" >
      <q-avatar size="32px" v-if="GetUserData()?.result?.imageUrl">
        <img :src="GetUserData()?.result?.imageUrl" >
      </q-avatar>
      <q-avatar size="32px" v-else>
        <img src="https://cdn-icons-png.flaticon.com/512/3237/3237472.png">
      </q-avatar>
          <q-badge
      v-if="(notificationNum + unReadedMessages) > 0"
      color="negative"
      floating
      rounded
      :label="notificationNum + unReadedMessages"
    />

    <q-menu>
      <q-list style="min-width: 120px;">
        <q-item clickable v-close-popup>
          <q-item-section @click="Profile">Profile</q-item-section>
        </q-item>
        <q-separator />
        <q-item clickable v-close-popup>
          <q-item-section @click="LogUserOut">Logout</q-item-section>
        </q-item>
      </q-list>
    </q-menu>
    </q-btn>

  </q-toolbar>

  <q-drawer
  v-model="mobileMenuOpen"
  overlay
  :width="200"
  bordered
  class="bg-grey-1"
  >

  <!-- search  -->
   <q-item>
    <q-item-section>
      <q-input
      dense
      outlined
      label="Search"
      @keyup.enter="GoSearch($event); mobileMenuOpen =false"
      />
    </q-item-section>
   </q-item>

   <q-separator />

   <!-- profile  -->
  <q-item clickable v-ripple @click="handleProfileClick" v-if="GetUserData()?.result">
    <q-item-section avatar>
      <q-icon name="eva-person-outline" />
    </q-item-section>
    <q-item-section>Profile</q-item-section>
  </q-item>

   <!-- message with bage  -->
    <q-item clickable v-ripple @click="handleChatClick" v-if="GetUserData()?.result">
      <q-item-section avatar>
        <dev class="mobile-icon-badge">
          <q-icon name="eva-message-square-outline" />
          <q-badge
            v-if="unReadedMessages > 0"
            color="negative"
            :label="unReadedMessages"
            />
        </dev>
      </q-item-section>
      <q-item-section>Messages</q-item-section>
    </q-item>

    <!-- notification with bade  -->
    <q-item clickable v-ripple @click="handleNotificationClick" v-if="GetUserData()?.result">
      <q-item-section avatar>
        <dev class="mobile-icon-badge">
          <q-icon name="eva-bell-outline" />
          <q-badge
            v-if="notificationNum > 0"
            color="negative"
            :label="notificationNum"
            />
        </dev>
      </q-item-section>
      <q-item-section>Notifications</q-item-section>
    </q-item>

    <q-separator />

    <!-- logout  -->
    <q-item clickable v-ripple @click="handleLogoutClick" v-if="GetUserData()?.result">
      <q-item-section avatar>
      <q-icon name="eva-log-out-outline" />
      </q-item-section>
      <q-item-section>Logout</q-item-section>
    </q-item>

    <!-- login for non auth users  -->
     <q-item clickable v-ripple to="/Auth" v-else>
      <q-item-section avatar>
        <q-icon name="eva-log-in-outline" />
      </q-item-section>
      <q-item-section>Login</q-item-section>
     </q-item>
</q-drawer>


 </q-header>


</template>
  
  <script>
  import { mapGetters, mapMutations, mapActions, mapState } from 'vuex';
  export default {
    name: 'NavBar',
    data (){
        return {
          notificationNum:0,
          unReadedMessages:0,
          mobileMenuOpen: false,
          // userData: null,
        }
    },
    computed: {
      ...mapGetters(['GetUserData']),
      ...mapState(['RealTimeNotify','RealTimeChat'])
    },
    watch:{
      "RealTimeNotify.notifyideslistNumber": async function () {
        this.UNreadedNotifyCount();
      },
      "RealTimeChat.privateMessages": async function(){
         this.unReadedMessages = this.unReadedMessages + 1;
       },
      $route: async function () {
        this.UNreadedNotifyCount();
        this.unreadMessageCount();
      }
    },

    methods: {
      ...mapMutations(['SetData']),
      ...mapActions(['logout','GetUnReadedNotifyNum', 'GetUnreadedMessageNum', 
        'StopConnectionToNotify',
        'StopConnectionToChat'
      ]),
      handleProfileClick(){
        this.mobileMenuOpen = false;
        this.Profile();
      },
      handleChatClick(){
        this.mobileMenuOpen = false;
        this.GoToChat();
      },
      handleNotificationClick(){
        this.mobileMenuOpen = false;
        this.GoToNotification();
      },
      handleLogoutClick(){
        this.mobileMenuOpen = false;
        this.LogUserOut();
      },
      GoSearch(e) {
            console.log("go", e.target.value)
            this.$router.push({path: `/Search`, query: { search: e.target.value }})
        },
      Profile() {
        let id = this.GetUserData()?.result?._id;
        this.$router.push(`/Profile/${id}`)
      },
      async LogUserOut(){
        await this.StopConnectionToNotify()

        var wait = await this.StopConnectionToChat()
        console.log(wait)
        this.logout()
        setTimeout(() => {
        this.$router.push(`/Auth`)          
        }, 1500);
      },
      GoToNotification(){
        this.$router.push('/Notification')
      },
      GoToChat(){
        this.$router.push('/Chat')
      },
      async UNreadedNotifyCount(){
                this.NotifyList = await this.GetUnReadedNotifyNum(this.GetUserData()?.result._id)
        let numofunreadednot = 0;
        this.NotifyList.forEach(el =>{
             if (! el.isreded ) {
                        numofunreadednot++;
             }
            })
        this.notificationNum = numofunreadednot;
      },
      async unreadMessageCount(){
    const  {total} = await this.GetUnreadedMessageNum(this.GetUserData()?.result._id)
    this.unReadedMessages = total;
      }
    },
    async mounted(){
    this.SetData();
    // getnot number
    await this.UNreadedNotifyCount();
    // get chat messages numbers
    const  {total} = await this.GetUnreadedMessageNum(this.GetUserData()?.result._id)
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


.mobile-icon-badge
    position: relative
    display: inline-block

    .q-badge
      position: absolute
      top: -6px 
      right:  -6px
      min-width: 18px
      height: 18px 
      font-size: 11px

</style>