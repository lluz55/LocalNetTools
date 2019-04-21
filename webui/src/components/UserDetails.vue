<template>
  <div class="hello">
    <ul id="container">
      <li>
        <div id="username" class="label-light-color"> Computer of: <strong class="label-color">{{ username }}</strong></div>
      </li>
      <li>
        <strong id="address" class="label-color">{{urlClean}}</strong>
      </li>
      <li>
        <div id="addressLabel" class="label-light-color">
          Use this address from other device to access this page
        </div>
      </li>
    </ul>
    
    
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import axios from 'axios'

@Component
export default class UserDetails extends Vue {
  private urlApi = window.location.href.replace("/#/", "/") + "api/"
  public username = "Lucas Luz"

  mounted () {
    console.log(this.urlApi)
    axios.get(this.urlApi + 'getuserdetail/').then( resp => {
      let data = resp.data
      console.log(data)
      if(!data.error) {
        this.username = data.message
      } 
    }).catch(err=> {
      console.log(err)
    })
  }

  get urlClean() {
    return this.urlApi.replace("/api/", "")
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">

  #container {
    margin: 0 auto;
    margin-top: 30px;
    max-width: 480px;
  }

  #username {
    font-size: 1.1em;
    margin-bottom: 30px;
  }


  #addressLabel {
    font-size: .7em;
  }

  #address {   
    font-size: 1.3em;
  }

  .label-color {
    color: #555;

  }

  .label-light-color {
    color: #777;

  }
</style>
