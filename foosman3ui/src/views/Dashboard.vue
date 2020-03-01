<template>
  <div class="dash-content">

    <div class="pure-g">
      <div class="pure-u-1-2">

        <table class="pure-table pure-table-horizontal" style="width:99%;margin-bottom:15px;">
          <thead>
            <tr>
              <th colspan="4" style="text-align:center">
                Matches in progress
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="m in inProgress" :key="m.id">
              <td>
                {{ m.table }}
              </td>
              <td>
                {{ m.team1Name }}
              </td>
              <td style="text-align:center;">
                vs
              </td>
              <td style="text-align:right">
                {{ m.team2Name }}
              </td>
            </tr>
          </tbody>
        </table>

      </div>
      <div class="pure-u-1-2">

        <table class="pure-table pure-table-horizontal" style="width:99%;margin-bottom:15px;">
          <thead>
            <tr>
              <th colspan="3" style="text-align:center">
                Upcoming matches
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="m in scheduled" :key="m.id">
              <td>
                {{ m.team1Name }}
              </td>
              <td style="text-align:center;">
                vs
              </td>
              <td style="text-align:right">
                {{ m.team2Name }}
              </td>
            </tr>
          </tbody>
        </table>

      </div>
    </div>

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
      inProgress: [],
      scheduled: []
    }
  },
  mounted: function () {
    this.load()
    window.setInterval(this.load, 8000)
  },
  methods: {
    load: function () {
      let self = this
      const inProgressRequest = this.axios.get(`tournaments/${this.id}/matches/in-progress`)
      const scheduledRequest = this.axios.get(`tournaments/${this.id}/matches/scheduled`)
      this.axios.all([inProgressRequest, scheduledRequest]).then(this.axios.spread(function (inProgressRes, scheduledRes) {
        self.inProgress = inProgressRes.data
        self.scheduled = scheduledRes.data
      }))
    }
  }
}
</script>

<style scoped>
body {
  background-color: #1d2024;
}
tr {
  background-color: black;
}
tr {
  border-color: red;
}
th {
  font-family: 'Press Start 2P';
  font-size: 16px !important;
  color: greenyellow;
}
td {
  font-size: 14px !important;
  color: white;
}
.dash-content {
  margin: 8px;
}
.dash-info {
  border:solid 1px black;
  padding:5px;
  height: auto;
}
</style>
