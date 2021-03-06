<template>
  <div v-if="tournament">
    <div><h1><soccer-icon /> {{tournament.name}}</h1></div>

    <table v-show="inProgress.length > 0" class="pure-table pure-table-horizontal" style="width:99%;margin-bottom:15px;">
      <thead>
        <tr>
          <th colspan="3" style="text-align:center">
            <play-icon /> Matches in progress
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="m in inProgress" :key="m.id">
          <td>
            {{ m.table }}
          </td>
          <td>
            <a :href="'#/team/' + m.team1Id">{{ m.team1Name }}</a>
          </td>
          <td>
            <a :href="'#/team/' + m.team2Id">{{ m.team2Name }}</a>
          </td>
        </tr>
      </tbody>
    </table>

    <table v-show="scheduled.length > 0" class="pure-table pure-table-horizontal" style="width:99%;margin-bottom:15px;">
      <thead>
        <tr>
          <th colspan=3 style="text-align:center">
            <alarm-icon /> Upcoming matches
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="m in scheduled" :key="m.id">
          <td>
            <a :href="'#/team/' + m.team1Id">{{ m.team1Name }}</a>
          </td>
          <td style="text-align:center">vs.</td>
          <td style="text-align:right">
            <a :href="'#/team/' + m.team2Id">{{ m.team2Name }}</a>
          </td>
        </tr>
      </tbody>
    </table>

    <template v-for="(matches, tier) in elimination">
      <table :key="'elim' + tier" class="pure-table pure-table-horizontal" style="width:99%;margin-bottom:15px;">
        <thead>
          <tr>
            <th colspan=2 style="text-align:center">
              {{ tier | playoff-tier }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="m in matches" :key="m.id">
            <td :class="elimClass(m, m.team1_id)">
              <a :href="'#/team/' + m.team1_id">{{ m.Team1.name }}</a>
            </td>
            <td :class="elimClass(m, m.team2_id)">
              <a :href="'#/team/' + m.team2_id">{{ m.Team2.name }}</a>
            </td>
          </tr>
        </tbody>
      </table>
    </template>

    <template v-for="g in tournament.groups">
      <table class="pure-table pure-table-horizontal" style="width:99%;margin-bottom:15px;" :key="g.id">
        <thead>
          <tr>
            <th>
              {{g.name}}
            </th>
            <th style="text-align:right">Ma</th>
            <th style="text-align:right">Wi</th>
            <th style="text-align:right">Dr</th>
            <th style="text-align:right">Lo</th>
            <th style="text-align:right">Pt</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in groupTeams(g.id)" :key="t.id">
            <td>
              <a :href="'#/team/' + t.id">{{t.name}}</a>
            </td>
            <td style="text-align:right">{{t.stats.PlayedCount}}</td>
            <td style="text-align:right">{{t.stats.Wins}}</td>
            <td style="text-align:right">{{t.stats.Draws}}</td>
            <td style="text-align:right">{{t.stats.Losses}}</td>
            <td style="text-align:right">{{t.stats.Points}}</td>
          </tr>
        </tbody>
      </table>
    </template>
    <p>
      <router-link :to="'/dashboard/' + tournament.id"><view-dashboard-icon /> Open dashboard</router-link>
    </p>
    <p>
      <router-link :to="'/register-results/' + tournament.id"><counter-icon /> Register results</router-link>
    </p>
  </div>
</template>

<script>
import PlayIcon from 'vue-material-design-icons/Play.vue'
import AlarmIcon from 'vue-material-design-icons/Alarm.vue'
import SoccerIcon from 'vue-material-design-icons/Soccer.vue'
import ViewDashboardIcon from 'vue-material-design-icons/ViewDashboard.vue'
import CounterIcon from 'vue-material-design-icons/Counter.vue'

export default {
  components: {
    PlayIcon,
    AlarmIcon,
    SoccerIcon,
    ViewDashboardIcon,
    CounterIcon
  },
  props: {
    id: {
      type: [Number, String],
      required: true
    }
  },
  data: function () {
    return {
      tournament: null,
      inProgress: [],
      scheduled: [],
      elimination: []
    }
  },
  watch: {
    id: async function () {
      await this.load()
    }
  },
  mounted: function () {
    this.load()
  },
  methods: {
    load: function () {
      let self = this
      const _ = this._
      const tournamentRequest = this.axios.get(`tournaments/${this.id}`)
      const scoresRequest = this.axios.get(`tournaments/${this.id}/scores`)
      const inProgressRequest = this.axios.get(`tournaments/${this.id}/matches/in-progress`)
      const scheduledRequest = this.axios.get(`tournaments/${this.id}/matches/scheduled`)
      const elimRequest = this.axios.get(`tournaments/${this.id}/elimination-matches`)
      this.axios
        .all([tournamentRequest, scoresRequest, inProgressRequest, scheduledRequest, elimRequest])
        .then(this.axios.spread(function (tournamentRes, scoresRes, inProgressRes, scheduledRes, elimRes) {
          self.tournament = tournamentRes.data
          self.inProgress = inProgressRes.data
          self.scheduled = scheduledRes.data
          self.elimination = self._.groupBy(elimRes.data, 'playoff_tier')

          const scores = scoresRes.data
          self.tournament.teams.forEach(t => {
            _.assignIn(t, {
              stats: _.find(scores, ['TeamID', t.id])
            })
          })
        }))
    },
    groupTeams: function (gID) {
      // TODO: Order teams by Points and then Wins
      return this._.orderBy(
        this.tournament.teams.filter(t => t.groupId === gID),
        ['stats.Points', 'stats.Wins'],
        ['desc', 'desc']
      )
    },
    elimClass: function (m, teamId) {
      const res = m.MatchResults[m.MatchResults[0].TeamID === teamId ? 0 : 1]
      return {
        winner: res.Win > 0,
        eliminated: res.Loss > 0
      }
    }
  }
}
</script>

<style scoped>
h1 {
  font-size:18px;
}
th {
  font-size:12px;
}
td {
  font-size:12px;
}
.eliminated {
  text-decoration: line-through;
}
.winner {
  font-weight: bold;
}
</style>
