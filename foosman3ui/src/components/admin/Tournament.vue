<template>
  <div>
    <div>Tournament admin</div>
    <div>ID: {{ id }}</div>

    <table class="pure-table pure-table-horizontal">
      <thead>
        <tr>
          <th>Team name</th>
          <th>Player 1</th>
          <th>Player 2</th>
          <th>Player 3</th>
          <th>&nbsp;</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="t in teams" :key="t.id">
          <td>{{ t.name }}</td>
          <td>{{ t.player1 }}</td>
          <td>{{ t.player2 }}</td>
          <td>{{ t.player3 }}</td>
          <td>
          </td>
        </tr>
      </tbody>
    </table>
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
      teams: []
    }
  },
  watch: {
    id: async function () {
      await this.load()
    }
  },
  mounted: async function () {
    await this.load()
  },
  methods: {
    load: async function () {
      const res = await this.axios.get(`http://localhost:1323/tournaments/${this.id}/teams`)
      this.teams = res.data
    }
  }
}
</script>
