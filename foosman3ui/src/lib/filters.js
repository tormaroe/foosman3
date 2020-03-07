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

Vue.filter('match-state', n => {
  switch (n) {
  case 0: return 'Planned'
  case 1: return 'Scheduled'
  case 2: return 'In progress'
  case 3: return 'Played'
  default: return n
  }
})

Vue.filter('playoff-tier', n => {
  switch (parseInt(n)) {
  case 1: return 'Final'
  case 2: return 'Semifinal'
  case 4: return 'Quarterfinal'
  case 8: return 'Last sixteen'
  case 16: return 'Last thirty two'
  default: return 'Unknown: ' + n
  }
})
