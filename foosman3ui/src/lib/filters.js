import Vue from 'vue'

Vue.filter('tournament-state', n => {
  switch (n) {
  case 0: return 'NEW'
  case 1: return 'GROUP PLAY STARTED'
  case 2: return 'GROUP PLAY ENDED'
  case 3: return 'ELIMINATION PLAY STARTED'
  case 4: return 'COMPLETE'
  default: return n
  }
})
