import { createStore } from 'vuex'
import Auth from './Auth'
import Users from './Users'
import Posts from './Posts'
import NotificationStore from './Notification'
import Chat from './Chat'
import RealTimeNotify from './RealTimeNotify'
export default createStore({

  modules: {
    Auth,
    Users,
    Posts,
    NotificationStore,
    Chat,
    RealTimeNotify
  }
})
