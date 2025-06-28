import store from "@/store";


export default function NumAuthGarud(to, from, next){
    if(store.state.Auth.authData){
        next('/')
    } else {
        next()
    }
}




