<template>
  <div class="q-pa-md" >
    <q-item :to="{path:`/Profile/${GetUserData()?.result?._id}`}"
            v-if="GetUserData()?.result" 
            clickable v-ripple >
        <q-item-section avatar>
            <q-icon color="primary" name="eva-person-outline" />
        </q-item-section>    
        <q-item-section>Profile</q-item-section>       
    </q-item>
    <q-item  to="/Auth" v-else   clickable v-ripple >
        <q-item-section avatar>
            <q-icon color="primary" name="eva-person-outline" />
        </q-item-section>    
        <q-item-section>Authentication</q-item-section>       
    </q-item>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
export default {
  name:'SideBar',
  computed:{
    ...mapGetters(['GetUserData']),
  }
}

</script>