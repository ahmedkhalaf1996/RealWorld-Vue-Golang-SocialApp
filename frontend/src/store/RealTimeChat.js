const RealTimeChat = {
    state: {
        ws: null,
        privateMessages: [],
        onlineFriends: [],
        userId: '',
        NumberOfMessgesReal: 0
    },
    getters: {
        Getuserid: (state) => () => {
            return state.userId
        },
        GetPrivateMessges: (state) => () => {
            return state.privateMessages
        },
        GetRealTimeNumberMessges: (state) => () => {
            return state.NumberOfMessgesReal
        },
        GetOnlinefriends: (state) => () => {
            return state.onlineFriends
        }
    },
    mutations: {
        SET_WS(state, ws) {
            state.ws = ws;
        },
        UpdateNumberOfMessages(state) {
            state.NumberOfMessgesReal = state.NumberOfMessgesReal + 1;
        },
        setOnlineUsers(state, onlineFriends) {
            state.onlineFriends = onlineFriends;
        },
        AddPrivateMessage(state, message) {
            state.privateMessages = message;
        },
        clearPrivateMessage(state) {
            state.privateMessages = [];
        },
        setUserId(state) {
            if (JSON.parse(localStorage.getItem('profile'))) {
                state.userId = JSON.parse(localStorage.getItem('profile'))?.result?._id;
            }
        }
    },
    actions: {
        async createChatConnection(context) {
            try {
                context.commit('setUserId');
                if (context.state.userId && context.state.ws == null) {
                    const uri = process.env.VUE_APP_RealTimeChatUrl
                    const ws = new WebSocket(`${uri}${context.state.userId}`)

                    ws.onopen = () => {
                        context.commit('SET_WS', ws)
                    }

                    ws.onmessage = (event) => {
                        const message = JSON.parse(event.data);
                        if (!message.onlineFriends) {
                            context.commit('UpdateNumberOfMessages')
                            context.commit('AddPrivateMessage', message)
                            console.log("store realtime", message)
                        } else {
                            const uniqueUsers = Array.from(new Set(message.onlineFriends));
                            context.commit('setOnlineUsers', uniqueUsers)
                        }
                    }
                }
            } catch (error) {
                console.log('E', error)
            }
        },
        async SendPrivateMessage(context, message){
            if(context.state.ws){
                return context.state.ws.send(JSON.stringify(message))
            }
        },
        async StopConnectionToChat(context) {
            try {
                if (context.state.ws){
                context.state.ws.close()
                context.commit('SET_WS', null);
                }

            } catch (error) {
                console.log("error", error)
            }
        }

    }


}


export default RealTimeChat;