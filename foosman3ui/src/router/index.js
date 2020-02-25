import Vue from 'vue'
import VueRouter from 'vue-router'

import Home from '../views/Home.vue'
import PublicTournamentList from '../components/PublicTournamentList.vue'
import PublicTournament from '../components/PublicTournament.vue'
import PublicGroup from '../components/PublicGroup.vue'
import PublicTeam from '../components/PublicTeam.vue'

import Dashboard from '../views/Dashboard.vue'

import Admin from '../views/Admin.vue'
import TournamentList from '../components/admin/TournamentList.vue'
import TournamentAdmin from '../components/admin/Tournament.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    component: Home,
    children: [
      {
        path: '',
        component: PublicTournamentList
      },
      {
        path: 'tournament/:id',
        component: PublicTournament,
        props: true
      },
      {
        path: 'group/:id',
        component: PublicGroup,
        props: true
      },
      {
        path: 'team/:id',
        component: PublicTeam,
        props: true
      }
    ]
  },
  {
    path: '/dashboard/:id',
    component: Dashboard,
    props: true
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
