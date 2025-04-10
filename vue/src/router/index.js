import Vue from 'vue'
import Router from 'vue-router'
import Index from '@/views/Index.vue'
import Share from '@/views/Share.vue'
import NotFound from '@/views/NotFound.vue'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'Index',
      component: Index
    },
    {
      path: '/share/:roomID',
      name: 'Share',
      component: Share
    },
    {
      path: '*',
      name: 'NotFound',
      component: NotFound
    }
  ]
})
