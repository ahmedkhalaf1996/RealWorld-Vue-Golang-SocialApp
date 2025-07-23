const RealTimeNotify = {
    state: {
        ws: null,
        notifyideslistNumber: 0,
        notifyidData: null,
    },
    getters:{
        Getnotifyideslist:(state) => () => {
            return state.notifyideslistNumber
        },
    },
    mutations:{
        SET_WS(state, ws){
            state.ws = ws;
        },
        ADD_NOTIFICATION(state, notify){
            state.notifyideslistNumber = state.notifyideslistNumber +1;
            state.notifyidData = notify;

        }
    },
    actions:{
        async connectToNotify(context){
            if(JSON.parse(localStorage.getItem('profile')) && context.state.ws == null) {
                const Userid = JSON.parse(localStorage.getItem('profile')).result._id;
                const uri = process.env.VUE_APP_RealTimeNotificationUrl;
                const ws = new WebSocket(`${uri}${Userid}`)

                ws.onopen = ()=> {
                    context.commit('SET_WS', ws);
                }

                ws.onmessage = (event) => {
                    const Notify = JSON.parse(event.data);
                    context.commit('ADD_NOTIFICATION', Notify)
                }
            }
        },

        async StopConnectionToNotify(context){
            try {
               context.ws.disconnect()
                context.commit('SET_WS', null);
            } catch (error) {
                console.log("error", error)          
            }
        }

    },


}

export default RealTimeNotify;


