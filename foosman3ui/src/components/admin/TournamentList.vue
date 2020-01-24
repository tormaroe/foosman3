<template>
  <div>
    <form class="pure-form">
      <fieldset>
        <legend>Add a tournament</legend>
        <input type="text" placeholder="Name" v-model="newTournament.name">
        <button
          type="button"
          @click="addTournament"
          class="pure-button pure-button-primary">
            Add
          </button>
      </fieldset>
    </form>
    <table class="pure-table pure-table-horizontal">
      <thead>
        <tr>
          <th>Name</th>
          <th>State</th>
          <th>Table count</th>
          <th>&nbsp;</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="t in tournaments" :key="t.id">
          <td>{{ t.name }}</td>
          <td>{{ t.state | tournament-state }}</td>
          <td>{{ t.tableCount }}</td>
          <td>
            <router-link :to="'/admin/tournament/' + t.id">Manage</router-link>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  name: 'TournamentList',
  data: function () {
    return {
      tournaments: [],
      newTournament: {
        name: '',
        tableCount: 1
      }
    }
  },
  mounted: function () {
    this.load()
  },
  methods: {
    load: function () {
      const self = this
      this.axios.get('http://localhost:1323/tournaments').then(function (res) {
        self.tournaments = res.data
      })
    },
    addTournament: async function () {
      // TODO: Validation
      await this.axios.post('http://localhost:1323/tournaments', this.newTournament)
      this.newTournament.name = ''
      await this.load()
    }
  }
}
</script>

<style scoped>
input {
  margin-right: 2px;
}
</style>
