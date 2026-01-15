
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
            if (context.state.userId) {
                // Close existing connection if any
                if (context.state.ws) {
                    context.state.ws.close();
                    context.commit('SET_WS', null);
                }
                
                // const uri = process.env.VUE_APP_RealTimeChatUrl
                // const ws = new WebSocket(`${uri}${context.state.userId}`)
                
                const baseUri = process.env.VUE_APP_RealTimeChatUrl
                const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
                const host = window.location.host;
                const wsUrl = `${protocol}//${host}${baseUri}/${context.state.userId}`  
                
                const ws = new WebSocket(wsUrl)
                


                ws.onopen = () => {
                    console.log('WebSocket connected for user:', context.state.userId);
                    context.commit('SET_WS', ws);
                }
                
                ws.onmessage = (event) => {
                    const message = JSON.parse(event.data);
                    if (!message.onlineFriends) {
                        context.commit('UpdateNumberOfMessages')
                        context.commit('AddPrivateMessage', message)
                        console.log("store realtime", message)
                    } else {
                        const uniqueUsers = Array.from(new Set(message.onlineFriends));
                        console.log("OnlineUsers", uniqueUsers)
                        context.commit('setOnlineUsers', uniqueUsers)
                    }
                }
                
                ws.onclose = (event) => {

                 if (JSON.parse(localStorage.getItem('profile'))) {

                    console.log('WebSocket disconnected:', event.code, event.reason);
                    context.commit('SET_WS', null);
                    // Attempt reconnection after 3 seconds
                    setTimeout(() => {
                        context.dispatch('createChatConnection');
                    }, 3000);
                }
                }
                
                ws.onerror = (error) => {
                    console.error('WebSocket error:', error);
                }
            }
        } catch (error) {

            console.log('Error creating WebSocket connection:', error);
        if (JSON.parse(localStorage.getItem('profile'))) {
            // make suere user still logedin
            // Retry after 5 seconds
            setTimeout(() => {
                context.dispatch('createChatConnection');
            }, 5000);
        }
        }
    },


        async SendPrivateMessage(context, message){
            if(context.state.ws){
                return context.state.ws.send(JSON.stringify(message))
            }
        },
        async StopConnectionToChat(context) {
            let success = false
            try {
                if (context.state.ws){
                context.state.ws.close()
                context.commit('SET_WS', null);
                success = true
                }

            } catch (error) {
                console.log("error", error)
                success = false
            } finally {
                context.commit("SET_WS", null);
            }
            return success
        }

    }


}


export default RealTimeChat;