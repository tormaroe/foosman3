import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import Admin from '../views/Admin.vue'
import TournamentList from '../components/admin/TournamentList.vue'
import TournamentAdmin from '../components/admin/Tournament.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'home',
    component: Home
  },
  {
    path: '/admin',
    component: Admin,
    children: [
      {
        path: '',
        component: TournamentList
      },
      {
        path: 'tournament/:id/',
        component: TournamentAdmin,
        props: true
      }
    ]
  }
]

const router = new VueRouter({
  routes
})

export default router
