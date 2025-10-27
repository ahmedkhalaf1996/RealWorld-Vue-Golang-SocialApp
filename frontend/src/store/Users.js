import * as api from '@/api/index.js'

const Users = {
    state: {
        User: null,
        UserPosts:[],
        PostsPagination:{
            currentPage:1,
            numberOfPages:0
        }
    },
    getters: {
        GetUser: (state) => () => {
            return state.User
        },
        GetUserFollowersFollowing: async () => {
            const userd = JSON.parse(localStorage.getItem('profile'));
            var followers = userd.result.followers || [];
            var following = userd.result.following || [];
            
            const combinedArray = [...followers, ...following];
            const uniqueArray = Array.from(new Set(combinedArray));

            var userdata = [];
            for(const uid of uniqueArray){
                const {data } = await api.fetchUserProfile(uid);
                var user = {"_id": data.user._id, "name": data.user.name, "imageUrl": data.user.imageUrl};
                userdata.push(user)
            }
            return userdata;
        },
        GetUserPosts:(state)=> {
            return state.UserPosts
        },
        GetPostsPagination:(state)=>{
            return state.PostsPagination
        }
    },
    mutations: {
        UserData(state, payload){
            state.User = payload?.data
        },
        SetUserPosts(state, payload){
            if(payload.append){
                const existingIds = new Set(state.UserPosts.map(post => post._id))
                const newPosts = payload.posts.filter(post => !existingIds.has(post._id))
                state.UserPosts = [...state.UserPosts, ...newPosts]
            } else {
                state.UserPosts = payload.posts || []
            }
        }, 
        SetPostsPagination(state, payload){
            state.PostsPagination = {
                currentPage: payload.currentPage || 1,
                numberOfPages: payload.numberOfPages || 0
            }
        },
        ResetUserPosts(state){
            state.UserPosts = []
            state.PostsPagination = {
                currentPage: 1,
                numberOfPages: 0
            }
        },
        UpdateUserFollowStatus(state, {userId, isFollowing, followersCount}) {
            if(state.User && state.User._id === userId){
                state.User.isFollowing = isFollowing
                state.User.followersCount = followersCount
            }
        }

    },
    actions: {
        // getuserbyid
        async GetUserByID(context, {id, page=1, append = false}) {
            try {
                if(!id || id === 'undefined' || id.trim() === ''){
                    throw new Error('User ID is Requreid')
                }

                const { data } = await api.fetchUserProfile(id, page);
                console.log("Store Get User Profile Check Cach", data?.cached)
                if(!data){
                    throw new Error('No data receved From API')
                }

                if(!append && data.user){
                    context.commit('UserData')
                }

                const posts = Array.isArray(data.posts) ? data.posts : []
                context.commit('SetUserPosts', {
                    posts:posts,
                    append:append
                })

                context.commit('SetPostsPagination', {
                    currentPage: data.currentPage || page,  
                    numberOfPages: data.numberOfPages || 0
                })

                // res 
                return {
                    user: data.user,
                    posts: posts,
                    currentPage: data.currentPage || page,
                    numberOfPages: data.numberOfPages || 0
                }

            } catch (error) {
              console.error('Get userby id Erorr store', error);
              if(error.response?.status === 502){
                console.error('Server is currently unavilable. ')
              }
            }
        },
        // reset user posts
        ResetUserPosts({commit}) {
            commit('ResetUserPosts')
        },
        // update user data
        async UpdateUserData(context, userData) {
            try {
                const {data} = await api.UpdateUser(userData);

                context.commit('UserData', data.user)

                return data;
            } catch (error) {
                console.log(error);
                return error;
            }
        },
        // following user
        async FollowUser(context, ProfileID) {
            try {
                const {data} = await api.following(ProfileID )

                return data
            } catch (error) {
                console.log(error)
                return error
            }
        },
        async GetTheUserSug(context, id){
            try {
                const {data} = await api.getSugUser(id)
                return data
            } catch (error) {
                console.log(error)
                return error
            }
        }
    }
}



export default Users



