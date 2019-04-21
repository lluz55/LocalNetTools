<template>
  <div class="home">
    <UserDetails msg="Welcome to Your Vue.js + TypeScript App"/>
    <ShutDownPC/>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import UserDetails from '@/components/UserDetails.vue'; // @ is an alias to /src
import ShutDownPC from '@/components/ShutDownPC.vue';

@Component({
  components: {
    UserDetails,
    ShutDownPC,
  },
})
export default class Home extends Vue {
  
}
</script>
