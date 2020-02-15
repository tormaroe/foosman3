<template>
  <div v-if="team">
    <div><b>{{team.name}}</b></div>
    <div><a :href="'#/group/' + team.groupId">{{team.groupName}}</a></div>
    <hr>
    <div v-show="team.player1 && team.player1.length > 0">{{team.player1}}</div>
    <div v-show="team.player2 && team.player2.length > 0">{{team.player2}}</div>
    <div v-show="team.player3 && team.player3.length > 0">{{team.player3}}</div>
    <hr>
    <table style="width:99%">
      <tr>
        <td>Total score</td>
        <td>{{team.totalScore}}</td>
      </tr>
      <tr>
        <td>Games Won</td>
        <td>{{team.gamesWon}}</td>
      </tr>
      <tr>
        <td>Games Drawn</td>
        <td>{{team.gamesDraw}}</td>
      </tr>
      <tr>
        <td>Games Lost</td>
        <td>{{team.gamesLost}}</td>
      </tr>
    </table>
    <hr>
    <div><b>Matches:</b></div>
    <table class="pure-table pure-table-horizontal" style="width:99%;margin-bottom:15px;">
      <thead>
        <tr>
          <th>
            Vs.
          </th>
          <th>
            Pts.
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="m in team.matches" :key="m.id">
          <td>
            ...
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
      team: null
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
      const res = await this.axios.get(`http://localhost:1323/teams/${this.id}`)
      this.team = res.data
    }
  }
}
</script>
