const RealTimeChat = {
    state: {
        ws: null,
        privateMessages: [],
        onlineFriends: [],
        userId: '',
        NumberOfMessgesReal: 0,
        reconnectAttempts: 0,
        maxReconnectAttempts: 5
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
            // Always create a new array to trigger reactivity
            state.onlineFriends = [...new Set(onlineFriends)];
            console.log('Updated online friends:', state.onlineFriends);
        },
        updateSingleUserStatus(state, {userID, online}) {
            if (online) {
                // Add user if not already in list
                if (!state.onlineFriends.includes(userID)) {
                    state.onlineFriends = [...state.onlineFriends, userID];
                }
            } else {
                // Remove user from list
                state.onlineFriends = state.onlineFriends.filter(id => id !== userID);
            }
            console.log(`User ${userID} is now ${online ? 'online' : 'offline'}`, state.onlineFriends);
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
        },
        incrementReconnectAttempts(state) {
            state.reconnectAttempts++;
        },
        resetReconnectAttempts(state) {
            state.reconnectAttempts = 0;
        }
    },
    actions: {
        async createChatConnection(context) {
            try {
                context.commit('setUserId');
                
                if (!context.state.userId) {
                    console.log('No user ID found, cannot connect');
                    return;
                }
                
                if (context.state.ws != null) {
                    console.log('WebSocket already connected');
                    return;
                }

                const uri = process.env.VUE_APP_RealTimeChatUrl
                const ws = new WebSocket(`${uri}${context.state.userId}`)

                ws.onopen = () => {
                    console.log('WebSocket connected successfully');
                    context.commit('SET_WS', ws);
                    context.commit('resetReconnectAttempts');
                }

                ws.onmessage = (event) => {
                    try {
                        const message = JSON.parse(event.data);
                        console.log('Received message:', message);

                        // Handle different message types
                        if (message.type === 'online_friends' || message.onlineFriends) {
                            const uniqueUsers = Array.from(new Set(message.onlineFriends || []));
                            console.log("Initial online users:", uniqueUsers);
                            context.commit('setOnlineUsers', uniqueUsers);
                        } 
                        else if (message.type === 'status_update') {
                            // Handle individual status update
                            console.log(`Status update: ${message.userID} is ${message.online ? 'online' : 'offline'}`);
                            context.commit('updateSingleUserStatus', {
                                userID: message.userID,
                                online: message.online
                            });
                            // Also update the full list if provided
                            if (message.onlineFriends) {
                                const uniqueUsers = Array.from(new Set(message.onlineFriends));
                                context.commit('setOnlineUsers', uniqueUsers);
                            }
                        }
                        else if (message.type === 'message' || (message.sender && message.content)) {
                            // Handle incoming message
                            context.commit('UpdateNumberOfMessages');
                            context.commit('AddPrivateMessage', message);
                            console.log("Received chat message:", message);
                        }
                        else if (message.type === 'message_sent') {
                            // Handle message sent confirmation
                            console.log("Message sent confirmation:", message);
                        }
                        else {
                            console.log('Unknown message type:', message);
                        }
                    } catch (error) {
                        console.error('Error parsing WebSocket message:', error);
                    }
                }

                ws.onerror = (error) => {
                    console.error('WebSocket error:', error);
                }

                ws.onclose = (event) => {
                    console.log('WebSocket closed:', event.code, event.reason);
                    context.commit('SET_WS', null);
                    
                    // Attempt reconnection
                    if (context.state.reconnectAttempts < context.state.maxReconnectAttempts) {
                        context.commit('incrementReconnectAttempts');
                        const delay = Math.min(1000 * Math.pow(2, context.state.reconnectAttempts), 30000);
                        console.log(`Reconnecting in ${delay}ms (attempt ${context.state.reconnectAttempts}/${context.state.maxReconnectAttempts})`);
                        setTimeout(() => {
                            context.dispatch('createChatConnection');
                        }, delay);
                    }
                }

            } catch (error) {
                console.error('Error creating WebSocket connection:', error);
            }
        },
        
        async SendPrivateMessage(context, message) {
            if (context.state.ws && context.state.ws.readyState === WebSocket.OPEN) {
                try {
                    context.state.ws.send(JSON.stringify(message));
                    console.log('Message sent:', message);
                    return true;
                } catch (error) {
                    console.error('Error sending message:', error);
                    return false;
                }
            } else {
                console.error('WebSocket is not open');
                return false;
            }
        },
        
        async StopConnectionToChat(context) {
            try {
                if (context.state.ws) {
                    context.state.ws.close();
                    context.commit('SET_WS', null);
                    context.commit('setOnlineUsers', []);
                    console.log('WebSocket connection closed');
                }
            } catch (error) {
                console.error("Error closing WebSocket:", error);
            }
        }
    }
}

export default RealTimeChat;

// const RealTimeChat = {
//     state: {
//         ws: null,
//         privateMessages: [],
//         onlineFriends: [],
//         userId: '',
//         NumberOfMessgesReal: 0
//     },
//     getters: {
//         Getuserid: (state) => () => {
//             return state.userId
//         },
//         GetPrivateMessges: (state) => () => {
//             return state.privateMessages
//         },
//         GetRealTimeNumberMessges: (state) => () => {
//             return state.NumberOfMessgesReal
//         },
//         GetOnlinefriends: (state) => () => {
//             return state.onlineFriends
//         }
//     },
//     mutations: {
//         SET_WS(state, ws) {
//             state.ws = ws;
//         },
//         UpdateNumberOfMessages(state) {
//             state.NumberOfMessgesReal = state.NumberOfMessgesReal + 1;
//         },
//         setOnlineUsers(state, onlineFriends) {
//             state.onlineFriends = onlineFriends;
//         },
//         AddPrivateMessage(state, message) {
//             state.privateMessages = message;
//         },
//         clearPrivateMessage(state) {
//             state.privateMessages = [];
//         },
//         setUserId(state) {
//             if (JSON.parse(localStorage.getItem('profile'))) {
//                 state.userId = JSON.parse(localStorage.getItem('profile'))?.result?._id;
//             }
//         }
//     },
//     actions: {
//      async createChatConnection(context) {
//         try {
//             context.commit('setUserId');
//             if (context.state.userId) {
//                 // Close existing connection if any
//                 if (context.state.ws) {
//                     context.state.ws.close();
//                     context.commit('SET_WS', null);
//                 }
                
//                 const uri = process.env.VUE_APP_RealTimeChatUrl
//                 const ws = new WebSocket(`${uri}${context.state.userId}`)
                
//                 ws.onopen = () => {
//                     console.log('WebSocket connected for user:', context.state.userId);
//                     context.commit('SET_WS', ws);
//                 }
                
//                 ws.onmessage = (event) => {
//                     const message = JSON.parse(event.data);
//                     if (!message.onlineFriends) {
//                         context.commit('UpdateNumberOfMessages')
//                         context.commit('AddPrivateMessage', message)
//                         console.log("store realtime", message)
//                     } else {
//                         const uniqueUsers = Array.from(new Set(message.onlineFriends));
//                         console.log("OnlineUsers", uniqueUsers)
//                         context.commit('setOnlineUsers', uniqueUsers)
//                     }
//                 }
                
//                 ws.onclose = (event) => {
//                     console.log('WebSocket disconnected:', event.code, event.reason);
//                     context.commit('SET_WS', null);
//                     // Attempt reconnection after 3 seconds
//                     setTimeout(() => {
//                         context.dispatch('createChatConnection');
//                     }, 3000);
//                 }
                
//                 ws.onerror = (error) => {
//                     console.error('WebSocket error:', error);
//                 }
//             }
//         } catch (error) {
//             console.log('Error creating WebSocket connection:', error);
//             // Retry after 5 seconds
//             setTimeout(() => {
//                 context.dispatch('createChatConnection');
//             }, 5000);
//         }
//     },


//         async SendPrivateMessage(context, message){
//             if(context.state.ws){
//                 return context.state.ws.send(JSON.stringify(message))
//             }
//         },
//         async StopConnectionToChat(context) {
//             try {
//                 if (context.state.ws){
//                 context.state.ws.close()
//                 context.commit('SET_WS', null);
//                 }

//             } catch (error) {
//                 console.log("error", error)
//             }
//         }

//     }


// }


// export default RealTimeChat;