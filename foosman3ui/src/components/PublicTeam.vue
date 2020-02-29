<template>
  <div v-if="teamdata">
    <div style="margin-bottom:14px"><a :href="'#/tournament/' + teamdata.Team.tournamentId">&laquo; Tournament overview</a></div>
    <div><h1>{{teamdata.Team.name}}</h1></div>
    <div>{{teamdata.Team.Group.name}}</div>
    <hr>
    <div v-show="teamdata.Team.player1 && teamdata.Team.player1.length > 0">{{teamdata.Team.player1}}</div>
    <div v-show="teamdata.Team.player2 && teamdata.Team.player2.length > 0">{{teamdata.Team.player2}}</div>
    <div v-show="teamdata.Team.player3 && teamdata.Team.player3.length > 0">{{teamdata.Team.player3}}</div>
    <hr>
    <table style="width:99%">
      <tr>
        <td><b>Total score</b></td>
        <td><b>{{stats.points}}</b></td>
      </tr>
      <tr>
        <td>Games Won</td>
        <td>{{stats.wins}}</td>
      </tr>
      <tr>
        <td>Games Drawn</td>
        <td>{{stats.draws}}</td>
      </tr>
      <tr>
        <td>Games Lost</td>
        <td>{{stats.losses}}</td>
      </tr>
    </table>
    <hr>
    <div><b>Matches:</b></div>
    <table class="pure-table pure-table-horizontal" style="width:99%;margin-bottom:15px;">
      <thead>
        <tr>
          <th style="width:12px;">&nbsp;</th>
          <th>
            Opponent
          </th>
          <th style="text-align:right">
            Pt
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="m in teamdata.Matches" :key="m.id">
          <td>
            <emoticon-outline-icon v-if="m.MatchResults.Win === 1" fillColor="#00FF00" />
            <bomb-icon v-else-if="m.MatchResults.Loss === 1" fillColor="#FF0000" />
            <scale-balance-icon v-else-if="m.MatchResults.Draw === 1" fillColor="#0000FF" />
            <help-icon v-else />
          </td>
          <td>
            <a :href="'#/team/' + m.Opponent.id">{{m.Opponent.name}}</a>
          </td>
          <td style="text-align:right">
            {{ m | matchPoints }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import EmoticonOutlineIcon from 'vue-material-design-icons/EmoticonOutline.vue'
import BombIcon from 'vue-material-design-icons/Bomb.vue'
import HelpIcon from 'vue-material-design-icons/Help.vue'
import ScaleBalanceIcon from 'vue-material-design-icons/ScaleBalance.vue'

export default {
  components: {
    EmoticonOutlineIcon,
    BombIcon,
    HelpIcon,
    ScaleBalanceIcon
  },
  props: {
    id: {
      type: [Number, String],
      required: true
    }
  },
  data: function () {
    return {
      teamdata: null,
      stats: {
        points: 0,
        wins: 0,
        draws: 0,
        losses: 0
      }
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
      const self = this
      const res = await this.axios.get(`http://localhost:1323/teams/${this.id}`)
      this.teamdata = res.data

      this.stats = {
        points: 0,
        wins: 0,
        draws: 0,
        losses: 0
      }

      this.teamdata.Matches.forEach(m => {
        const res0 = m.MatchResults[0]
        const res1 = m.MatchResults[1]
        if (res0.TeamID === self.teamdata.Team.id) {
          m.MatchResults = res0
          m.Opponent = m.Team2
        } else {
          m.MatchResults = res1
          m.Opponent = m.Team1
        }
        self.stats.points += m.MatchResults.Points
        self.stats.wins += m.MatchResults.Win
        self.stats.draws += m.MatchResults.Draw
        self.stats.losses += m.MatchResults.Loss
      })

      this.teamdata.Matches = this._.sortBy(this.teamdata.Matches, ['Sequence'])
    }
  },
  filters: {
    matchPoints: function (m) {
      if (m.state === 3) {
        return m.MatchResults.Points
      }
      return ''
    }
  }
}
</script>

<style scoped>
h1 {
  font-size:18px;
}
div {
  font-size:12px;
}
th {
  font-size:12px;
}
td {
  font-size:12px;
}
</style>
