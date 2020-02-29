<template>
  <div v-if="tournament">
    <div><h1>{{tournament.name}}</h1></div>

    <table class="pure-table pure-table-horizontal" style="width:99%;margin-bottom:15px;">
      <thead>
        <tr>
          <th colspan="2" style="text-align:center">
            Matches in progress
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="m in inProgress" :key="m.id">
          <td>
            {{ m.team1Name }} vs {{ m.team2Name }}
          </td>
          <td>
            {{ m.table }}
          </td>
        </tr>
      </tbody>
    </table>

    <table class="pure-table pure-table-horizontal" style="width:99%;margin-bottom:15px;">
      <thead>
        <tr>
          <th style="text-align:center">
            Upcoming matches
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="m in scheduled" :key="m.id">
          <td>
            {{ m.team1Name }} vs {{ m.team2Name }}
          </td>
        </tr>
      </tbody>
    </table>

    <template v-for="g in tournament.groups">
      <table class="pure-table pure-table-horizontal" style="width:99%;margin-bottom:15px;" :key="g.id">
        <thead>
          <tr>
            <th>
              <a :href="'#/group/' + g.id">{{g.name}}</a>
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
      <router-link :to="'/dashboard/' + tournament.id">Open dashboard</router-link>
    </p>
    <p>
      <router-link :to="'/register-results/' + tournament.id">Register results</router-link>
    </p>
  </div>
</template>

<script>
export default {
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
      scheduled: []
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
      const tournamentRequest = this.axios.get(`http://localhost:1323/tournaments/${this.id}`)
      const scoresRequest = this.axios.get(`http://localhost:1323/tournaments/${this.id}/scores`)
      const inProgressRequest = this.axios.get(`http://localhost:1323/tournaments/${this.id}/matches/in-progress`)
      const scheduledRequest = this.axios.get(`http://localhost:1323/tournaments/${this.id}/matches/scheduled`)
      this.axios
        .all([tournamentRequest, scoresRequest, inProgressRequest, scheduledRequest])
        .then(this.axios.spread(function (tournamentRes, scoresRes, inProgressRes, scheduledRes) {
          self.tournament = tournamentRes.data
          self.inProgress = inProgressRes.data
          self.scheduled = scheduledRes.data

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
      return this.tournament.teams.filter(t => t.groupId === gID)
    }
  }
}
</script>

<style scoped>
th {
  font-size:12px;
}
td {
  font-size:12px;
}
</style>
