<template>
  <div>
    <div v-if="editMatch" class="edit-container">
      <p><b>EDIT MATCH</b></p>
      <p>
        {{ editMatch.Team1.name }} vs. {{ editMatch.Team2.name }}
      </p>
      <div>
        <button @click="cancelEdit" class="pure-button" title="Cancel">Cancel</button>
        &nbsp;
        <button @click="unplayMatch" class="pure-button" style="background-color:red;color:white" title="Unplay">Unplay match</button>
      </div>
    </div>
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
          <th>&nbsp;</th>
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
          <td>
            <button v-if="m.state === 3" @click="edit(m)" class="pure-button" style="background-color:orange" title="Fix match">Fix</button>
          </td>
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
      matches: [],
      editMatch: null
    }
  },
  watch: {
    tournamentId: {
      immediate: true,
      async handler (value) {
        await this.loadMatches(value)
      }
    }
  },
  methods: {
    loadMatches: async function (id) {
      const res = await this.axios.get(`http://localhost:1323/tournaments/${id}/matches`)
      this.matches = res.data
    },
    edit: function (m) {
      this.editMatch = m
    },
    cancelEdit: function () {
      this.editMatch = null
    },
    unplayMatch: async function () {
      if (window.confirm('Are you sure?')) {
        await this.axios.post(`http://localhost:1323/matches/${this.editMatch.id}/reset`)
        this.cancelEdit()
        await this.loadMatches(this.tournamentId)
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

<style scoped>
.edit-container {
  border: solid 1px silver;
  margin-bottom: 10px;
  padding:12px;
}
</style>
