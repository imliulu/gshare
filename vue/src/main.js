import Vue from 'vue'
import App from './App.vue'
import router from './router'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import apiClient from "@/api";

Vue.config.productionTip = false

Vue.use(ElementUI)
Vue.prototype.$api = apiClient

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
