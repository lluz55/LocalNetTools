<template>
<div>
    <div id="container">
        <div id="title">
            Shutdown PC
        </div>
        <div id="content">
            <div id="delayLabel">
                Delay <em>(Minutes)</em>
            </div>
            <div id="delayValue">
                <input v-model="delay" type="number" name="delay" id="delay">
            </div>
        </div>
        <div id="actions">
            <button class="btn btn-cancel" @click="cancel">Cancel</button>
            <button class="btn" @click="schedule">Schedule</button>
        </div>

    </div>

        <transition name="fade-down">
            <div v-if="shuttingDownPC" class="shutDownNotice">
                <div class="warning">
                    This PC will be shut down in {{remaningTime}} minute{{ remaningTime > 1 ? 's':''}}
                </div>
            </div>
        </transition>
        <transition name="fade-down">
            <div v-if="shutdownCanceled" class="shutDownNotice">
                <div class="warning">
                    The shut down was canceld
                </div>
            </div>
        </transition>
</div>
</template>


<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import axios from 'axios'

@Component
export default class ShutDownPC extends Vue {
    delay = 1
    shuttingDownPC = false
    shutdownCanceled = false
    remaningTime = 1

    private urlApi = window.location.href.replace("/#/", "/") + "api/"

    mounted () {
        console.log(this.urlApi)

        // Check if has shut down active
        axios.get(this.urlApi + "shutdown/timeleft").then(resp => {
            let data = resp.data
            if(!data.error) {
                this.remaningTime = (data.message - (data.message % 60)) / 60
                this.shuttingDownPC = true
            }
        })
    }

    schedule(){
        axios.get(this.urlApi + "shutdown/" + this.delay).then(resp => {
            let data = resp.data
            if(!data.error) {
                this.shuttingDownPC = true
            }
        })
    }

    cancel(){
        axios.get(this.urlApi + "shutdown/c").then(resp => {
            let data = resp.data
            if(!data.error) {
                this.shuttingDownPC = false 
                setTimeout(() => {
                    this.shutdownCanceled = true
                    setTimeout(() => {
                        this.shutdownCanceled = false
                    }, 1500);
                }, 500);
            }
        })
    }

}
</script>


<style scoped>
    #container {
        margin: 0 auto;
        padding: 10px;
        margin-top: 30px;
        border: 1px solid #f5f5f5;
        max-width: 450px;
        box-shadow: 0 1px 1px 0 rgba(0,0,0,.2);
        border-radius: 5px;
        background-color: #fff;
    }

    #title {
        font-size: 1.2em;
        font-weight: 400;
        margin: 10px 0;
    }

    #content {
        margin-top: 20px;
        font-size: .9em;
    }

    #content:after {
        clear: both;
        content: '';
        display: table;
    }

    #delayLabel {
        float: left;
        padding-top: 5px;
    }

    #delayValue {
        float: right;
    }

    #delay {
        border-radius: 5px;
        border: 1px solid #d5d5d5;
        padding: 5px 10px;
        text-align: right;
        font-weight: 600;
        max-width: 70px;
    }

    #delay:focus {
        background-color: #f5f5f5;
        
    }

    #actions {
        margin-top: 10px;
        text-align: right;
    }

    .btn {
        margin: 5px;
        display: inline-block;
        padding: 5px 10px;
        border: 1px solid #c5c5c5;
        box-shadow: 0 1px 1px 0 #d5d5d5;
        border-radius: 5px;
        background-color: #fff;
        color: #3C77CE;
        user-select: none;
    }

    .btn:hover {
        cursor: pointer;
    }

    .btn:active {
        background-color: #CEDEF6;

    }

    .btn-cancel {
        color: #F78080;
    }

   .shutDownNotice {
        margin-top: 15px;
    }

    .fade-down-enter-active {
        transition: all .35s ease;
    }

    .fade-down-enter, .fade-down-leave-to {
        opacity: 0;
        transform: translateY(-20px);
    }
</style>
