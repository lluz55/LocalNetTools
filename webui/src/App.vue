<template>
  <div id="app">
    <router-view/>
  </div>
</template>

<style lang="scss">
* {
    font-family: 'Droid Sans Mono';
    outline: none;
}

body {
  
  background-color: #f4f8f9;
}
#app {
  font-family: 'Droid Sans Mono', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #050505;
}

ul {
  margin: 0;
  padding: 0;
  list-style: none;
}

</style>


<script lang="ts">
require('typeface-droid-sans-mono')

import Vue from 'vue'
export default Vue.extend({
  
})
</script>
