import Vue from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import VueAxios from 'vue-axios'
import VueDraggable from 'vue-draggable'
import VueLodash from 'vue-lodash'
import lodash from 'lodash'
import 'vue-material-design-icons/styles.css'
import './lib/filters'

Vue.config.productionTip = false

axios.defaults.baseURL = process.env.VUE_APP_API
axios.interceptors.response.use(
  function (response) { return response },
  function (error) {
    if (error.response) {
      alert(error.response.data.message)
    }
    return undefined
  })
Vue.use(VueAxios, axios)

Vue.use(VueDraggable)

Vue.use(VueLodash, { lodash: lodash })

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
