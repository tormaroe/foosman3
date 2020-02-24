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
            <th style="text-align:center">
              <a :href="'#/group/' + g.id">{{g.name}}</a>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in groupTeams(g.id)" :key="t.id">
            <td>
              <a :href="'#/team/' + t.id">{{t.name}}</a>
            </td>
          </tr>
        </tbody>
      </table>
    </template>
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
      const tournamentRequest = this.axios.get(`http://localhost:1323/tournaments/${this.id}`)
      const inProgressRequest = this.axios.get(`http://localhost:1323/tournaments/${this.id}/matches/in-progress`)
      const scheduledRequest = this.axios.get(`http://localhost:1323/tournaments/${this.id}/matches/scheduled`)
      this.axios.all([tournamentRequest, inProgressRequest, scheduledRequest]).then(this.axios.spread(function (tournamentRes, inProgressRes, scheduledRes) {
        self.tournament = tournamentRes.data
        self.inProgress = inProgressRes.data
        self.scheduled = scheduledRes.data
      }))
    },
    groupTeams: function (gID) {
      return this.tournament.teams.filter(t => t.groupId === gID)
    }
  }
}
</script>
