import Vue from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import VueAxios from 'vue-axios'
import VueDraggable from 'vue-draggable'
import './lib/filters'

Vue.config.productionTip = false

axios.interceptors.response.use(
  function (response) { return response },
  function (error) {
    if (error.response) {
      alert(error.response.data.message)
    }
  })
Vue.use(VueAxios, axios)

Vue.use(VueDraggable)

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')