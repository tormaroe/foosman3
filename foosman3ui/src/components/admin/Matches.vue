<template>
  <div>
    <p>
      Match count: {{ matches.length }}
    </p>
    <button
      v-show="matches.length === 0"
      class="pure-button pure-button-primary"
      type="button"
      @click="generateMatches">
      Generate matches
    </button>
    <button
      v-show="matches.length > 0"
      class="pure-button pure-button-primary"
      type="button"
      @click="clearMatches">
      Clear matches
    </button>
    <table class="pure-table pure-table-horizontal" style="width:99%">
      <thead>
        <tr>
          <th>Group</th>
          <th>Team 1</th>
          <th>Team 2</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="m in matches" :key="m.GroupName">
          <td>{{ m.GroupName }}</td>
          <td>{{ m.Team1 }}</td>
          <td>{{ m.Team2 }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  props: {
    tournament: {
      type: Object,
      required: true
    }
  },
  data: function () {
    return {
      matches: []
    }
  },
  mounted: async function () {
    await this.load()
  },
  methods: {
    load: async function () {
      const res = await this.axios.get(`http://localhost:1323/tournaments/${this.tournament.id}/matches`)
      this.matches = res.data
    },
    generateMatches: async function () {
      await this.axios.post(`http://localhost:1323/tournaments/${this.tournament.id}/generate-matches`)
      await this.load()
    },
    clearMatches: async function () {
      await this.axios.post(`http://localhost:1323/tournaments/${this.tournament.id}/clear-matches`)
      await this.load()
    }
  }
}
</script>
