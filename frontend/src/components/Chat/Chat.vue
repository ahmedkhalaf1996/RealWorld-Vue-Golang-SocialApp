<template>
    <q-page class="constrain q-pa-md">
        <div class="row q-col-gutter-lg">
            <div class="col-12 chat-container">
                <div class="user-list">
                    <div class="q-pa-md">
                        <q-toolbar class="bg-primary text-white shadow-1">
                            <q-toolbar-title>Following & Followers</q-toolbar-title>
                        </q-toolbar>

                        <q-list bordered>
                            <q-item
                                @click="selectUser(contact)"
                                v-for="contact in contacts" :key="contact._id" class="q-my-sm" clickable v-ripple>
                             <q-item-section avatar>
                                <q-avatar v-if="!contact.imageUrl" color="primary" text-color="white">
                                    {{  contact.name[0]  }}
                                </q-avatar>
                                <q-avatar v-else>
                                    <img :src="contact?.imageUrl">
                                </q-avatar>
                             </q-item-section>
                             <q-item-section>
                                <q-item-label>{{ contact.name }}</q-item-label>
                             </q-item-section>
                            
                            <q-item-section side  v-if="contact.isOnline">
                                <q-badge color="positive" rounded/>
                             </q-item-section>  

                             <q-item-section side v-if="contact.unReadedmessage && contact.unReadedmessage > 0" >
                                <q-badge color="negative" rounded :label="contact?.unReadedmessage" />
                             </q-item-section>


                            </q-item>
                        </q-list>
                    </div>
                </div>

                <!-- chat box  -->
                <div class="chat-messages" v-if="selectedUser != null" style="background: white;">
                    <div class="q-pa-md row justify-center"
                    style=" overflow-y: auto; max-height: 400px;"
                    ref="messageContainer"
                    @scroll="handleScroll"
                    >
                    <div v-for="msg in messageBetweenUsers" :key="msg._id" style="width: 100%;">
                        <q-chat-message
                            :name="msg.sender === MainUserData._id ? MainUserData.name : selectedUser.name"
                            :avatar="getAvatar(msg)"
                            :text="[msg.content]"
                            :sent="msg.sender === MainUserData._id ? true : false"
                            />
                    </div>
                </div>

                <q-separator spaced />
                <q-input outlined v-model="messaageToSend.text" @keyup.enter="Sendmessage" label="write message..">
                    <q-btn
                        v-if="messaageToSend.text != ''"
                        @click="Sendmessage"
                        flat
                        round 
                        color="primary"
                        icon="eva-arrow-right"
                        />
                </q-input>

              </div>

            </div>
        </div>
    </q-page>




</template>

<script>
import { mapGetters, mapActions, mapState } from 'vuex';
export default {
    name:'ChatComponent',
    data(){
        return {
            messaageToSend: {text: ''},
            contacts:[],
            messageBetweenUsers:[],
            messagelistnum:0,
            selectedUser: null,
            MainUserData:{},
            uniqueOnlineUsers:[],
        };
    },
    computed:{
        ...mapGetters(['GetUserFollowersFollowing','GetUserData']),
        ...mapState(["RealTimeChat"])
    },
    watch: {
    "RealTimeChat.onlineFriends": {
        handler: function (onlineFriendsArray) {
            console.log('Online friends changed:', onlineFriendsArray);
            if (Array.isArray(onlineFriendsArray)) {
                this.uniqueOnlineUsers = [...new Set(onlineFriendsArray)];
                this.updateOnlineList();
            }
        },
        deep: true,
        immediate: true
    },
    "RealTimeChat.privateMessages": function(message){
        if (!message || !message.sender) return;
        
        if (this.contacts.length > 0) {
            this.contacts.forEach((contact) => {
                if (contact._id == message.sender) {
                    contact.unReadedmessage = (contact.unReadedmessage || 0) + 1;
                }
            });
            
            // If viewing conversation with sender, add message
            if (this.selectedUser && this.selectedUser._id == message.sender) {
                this.messageBetweenUsers.push(message);
                this.$nextTick(() => {
                    this.scrollDownFunction();
                });
            }
        }
    }
},
    // watch: {
    //     "RealTimeChat.onlineFriends": function (online) {
    //         const onlineFriendsArray = Object.values(online);
    //         console.log('Online friend changed new val', onlineFriendsArray)
    //             this.uniqueOnlineUsers = Array.from(new Set(onlineFriendsArray));
    //             this.updateOnlineList();
    //     },
    //     "RealTimeChat.privateMessages": function(message){
    //         if(this.contacts.length > 0){
    //             this.contacts.forEach((contact)=> {
    //                 if(contact._id == message.sender) {
    //                     contact.unReadedmessage++;
    //                 }
    //             })
    //         if(this.selectedUser && this.selectedUser?._id == message.sender){
    //             this.messageBetweenUsers.push(message);
    //             setTimeout(() => {
    //                 this.scrollDownFunction();
    //             }, 100);
    //         }

    //         }
    //     }
    // },
    async mounted(){
        this.RefreshUserData();
        this.MainUserData = this.GetUserData().result;
        this.GetUsList();

        this.uniqueOnlineUsers = Array.from(new Set(Object.values(this.RealTimeChat.onlineFriends)));
        this.updateOnlineList();

    },
    methods:{
        ...mapActions([
            'GetUnreadedMessageNum',
            'GetChatMsgsBetweenTwoUsers',
            'SendMessage',
            'MarkMsgsAsReaded',
            'SendPrivateMessage',
            'RefreshUserData'
        ]),
        getAvatar(msg){
            const user = msg.sender === this.MainUserData._id ? this.MainUserData : this.selectedUser
            const avatar = Array.isArray(user.imageUrl) ? user.imageUrl[0]: user.imageUrl

            return avatar && avatar.trim() !== '' ? avatar : 'https://cdn-icons-png.flaticon.com/512/3237/3237472.png'
        },
        updateOnlineList(){
            this.contacts.forEach((contact)=> {
                if(this.uniqueOnlineUsers.includes(contact._id)) {
                    contact.isOnline = true;
                } else {
                    contact.isOnline = false;
                }
            })
        },
        handleScroll(){
            const container = this.$refs.messageContainer;
            if (container.scrollTop === 0){
                // scorelled to the top
                this.GetOldestMessgesBetweenUsers();
            }},

        async GetOldestMessgesBetweenUsers(){
            this.messagelistnum = this.messagelistnum +1;
            var firstuid = this.MainUserData._id
            var seconduid = this.selectedUser._id
            var from = this.messagelistnum;
            var ndata = {from, firstuid, seconduid};

            var {msgs} = await this.GetChatMsgsBetweenTwoUsers(ndata);
            this.messageBetweenUsers.unshift(...msgs);

        },
        scrollDownFunction(){
            const container = this.$refs.messageContainer;
            container.scrollTop = container.scrollHeight;
        },
        async CallMarkMsgAsReaded(user){
            var mainuid = this.MainUserData._id;
            var otheruid = user._id;
            var GetunReadedmessage = 0

            this.contacts.forEach(
                user => {
                    if(String(otheruid) == String(user._id)){
                        GetunReadedmessage = user.unReadedmessage
                    }
                }
            )

            var data = {mainuid, otheruid, GetunReadedmessage}
            var {isMarked} = await this.MarkMsgsAsReaded(data);

            if(isMarked){
                this.contacts.forEach(user => {
                    if(String(otheruid)== String(user._id)){
                        user.unReadedmessage = 0;
                    }
                })
            }
        },
        async GetUnreadedMsgList(){
            var {messages} = await this.GetUnreadedMessageNum(this.MainUserData._id);
            this.contacts.forEach(user => {
                messages.forEach(msg => {
                    if(String(msg.otherUserid) == String(user._id)){
                        user.unReadedmessage = Number(msg.numOfUnreadedMessages);
                    }
                })
            })
        },
        async GetUsList(){
            this.contacts = [];
            var glist = await this.GetUserFollowersFollowing;
            this.contacts = glist;
            if(this.contacts){
                this.GetUnreadedMsgList();
            }
            this.updateOnlineList();

        },
        async selectUser(user){
            this.selectedUser = null;
            this.messageBetweenUsers = [];

            this.selectedUser = user;
            this.messagelistnum = 0;
            var firstuid = this.MainUserData._id;
            var seconduid = user._id;
            var from = 0;
            var ndata = {from, firstuid, seconduid};
            var {msgs} = await this.GetChatMsgsBetweenTwoUsers(ndata);
            this.messageBetweenUsers.push(...msgs);
            setTimeout(() => {
                this.scrollDownFunction();
                this.CallMarkMsgAsReaded(user)
            }, 100);

        },
      Sendmessage() {
    if (!this.messaageToSend.text.trim()) {
        return;
    }
    
    var content = this.messaageToSend.text;
    var sender = this.MainUserData._id;
    var recever = this.selectedUser._id;

    var sdata = {
        content: content,
        sender: sender,
        recever: recever
    };
    
    // Always send through WebSocket if connected
    const isReceiverOnline = this.uniqueOnlineUsers.includes(recever);
    
    if (isReceiverOnline) {
        console.log('Sending message via WebSocket (receiver is online)');
        this.SendPrivateMessage(sdata).then((success) => {
            if (success) {
                // Message will be echoed back from server
                console.log('Message sent successfully');
            } else {
                // Fallback to HTTP if WebSocket fails
                console.log('WebSocket send failed, using HTTP fallback');
                this.SendMessage(sdata).then(() => {
                    this.messageBetweenUsers.push(sdata);
                    this.$nextTick(() => {
                        this.scrollDownFunction();
                    });
                });
            }
        });
    } else {
        console.log('Sending message via HTTP (receiver is offline)');
        this.SendMessage(sdata).then((success) => {
            if (success) {
                this.messageBetweenUsers.push(sdata);
                this.$nextTick(() => {
                    this.scrollDownFunction();
                });
            }
        });
    }

    this.messaageToSend.text = '';
}
 
    }
}

</script>

<style scoped>
.chat-container {
    display: flex;
}

.chat-messages {
    flex: 1;
    padding: 10px;
}


</style>




