<template>
  <div class="dash-content">

    <div class="pure-g" v-if="flash">
      <div class="pure-u-1-1 flash">
        <span v-html="flash.Raw"></span>
      </div>
    </div>

    <div class="pure-g">
      <div class="pure-u-1-2" v-show="inProgress.length > 0">

        <table class="pure-table pure-table-horizontal" style="width:99%;margin-bottom:15px;">
          <thead>
            <tr>
              <th colspan="3" style="text-align:center">
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
              <td>
                {{ m.team2Name }}
              </td>
            </tr>
          </tbody>
        </table>

      </div>
      <div class="pure-u-1-2" v-show="scheduled.length > 0">

        <table class="pure-table pure-table-horizontal" style="width:99%;margin-bottom:15px;">
          <thead>
            <tr>
              <th colspan="2" style="text-align:center">
                Upcoming matches
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="m in scheduled" :key="m.id">
              <td>
                {{ m.team1Name }}
              </td>
              <td>
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
      flash: null,
      inProgress: [],
      scheduled: [],
      intervalId: null
    }
  },
  mounted: function () {
    this.load()
    this.intervalId = window.setInterval(this.load, 8000)
  },
  beforeDestroy: function () {
    window.clearInterval(this.intervalId)
  },
  methods: {
    load: function () {
      let self = this
      const dashboardRequest = this.axios.get(`tournaments/${this.id}/dashboard`)
      const inProgressRequest = this.axios.get(`tournaments/${this.id}/matches/in-progress`)
      const scheduledRequest = this.axios.get(`tournaments/${this.id}/matches/scheduled`)
      this.axios.all([dashboardRequest, inProgressRequest, scheduledRequest]).then(this.axios.spread(function (dashboardRes, inProgressRes, scheduledRes) {
        self.flash = dashboardRes.data
        self.inProgress = inProgressRes.data
        self.scheduled = scheduledRes.data
      }))
    }
  }
}
</script>

<style scoped>
.flash {
  padding-top: 50px;
  padding-bottom: 50px;
  padding-left: 30px;
  padding-right: 30px;
  font-weight: bold;
  font-size: 24px;
  font-family: 'Baloo Bhaina', cursive;
  text-align: center;
  line-height: 32px;
}
table {
  border: none;
}
th {
  font-family: 'Baloo Bhaina', cursive;
  font-size: 16px !important;
  background-color: white;
  border-bottom: solid 1px black;
  border-top:none;
}
td {
  font-size: 14px !important;
  font-weight: bold;
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
