<template>
  <div>
    <table class="pure-table pure-table-horizontal" style="width:99%">
      <thead>
        <tr>
          <th>Seq#</th>
          <th>Team 1</th>
          <th>Team 2</th>
          <th>Result</th>
          <th>Group</th>
          <th>Table</th>
          <th>State</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="m in matches" :key="m.id">
          <td>{{ m.Sequence }}</td>
          <td>{{ m.Team1.name }}</td>
          <td>{{ m.Team2.name }}</td>
          <td>{{ m | resultDescription }}</td>
          <td>{{ m.Group.name }}</td>
          <td>{{ m.table }}</td>
          <td>{{ m.state }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  props: {
    tournamentId: {
      type: [Number, String],
      required: true
    }
  },
  data: function () {
    return {
      matches: []
    }
  },
  watch: {
    tournamentId: {
      immediate: true,
      async handler (value) {
        const res = await this.axios.get(`http://localhost:1323/tournaments/${value}/matches`)
        this.matches = res.data
      }
    }
  },
  filters: {
    resultDescription: function (match) {
      if (match.state < 3) {
        return '..'
      }
      let it1 = 0
      let it2 = 1
      if (match.MatchResults[it1].TeamID !== match.team1_id) {
        it1 = 1
        it2 = 0
      }
      if (match.MatchResults[it1].Win > 0) {
        return 'Winner: Team 1'
      } else if (match.MatchResults[it2].Win > 0) {
        return 'Winner: Team 2'
      } else {
        return 'Draw'
      }
    }
  }
}
</script>
