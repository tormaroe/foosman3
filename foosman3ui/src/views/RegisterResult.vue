<template>
  <div class="regres-content">
    <div style="margin-bottom:14px"><a :href="'#/tournament/' + id">&laquo; Tournament overview</a></div>
    <template v-if="selectedResult">
      <p>
        Please confirm result:
      </p>
      <p>
        {{ selectedResult.description }}
      </p>
      <div class="result" @click="confirmResult">
        <b>YES</b>
      </div>
      <div class="result" @click="cancel">
        <b>Cancel</b>
      </div>
    </template>
    <template v-else-if="selectedMatch">
      <p>
        <b>Choose result..</b>
      </p>
      <div class="result" @click="setWinner(selectedMatch.team1Id, selectedMatch.team1Name)">
        {{ selectedMatch.team1Name }} won!
      </div>
      <div class="result" @click="setWinner(selectedMatch.team2Id, selectedMatch.team2Name)">
        {{ selectedMatch.team2Name }} won!
      </div>
      <div class="result" @click="setDraw()" v-show="selectedMatch.groupName != ''">
        Draw!
      </div>
    </template>
    <template v-else>
      <p>
        <b>Ongoing matches</b>
      </p>
      <div class="match" v-for="m in inProgress" :key="m.id" @click="selectMatch(m)">
        <div>{{ m.table }}</div>
        <div><b>{{ m.team1Name }}</b></div>
        <div><b>{{ m.team2Name }}</b></div>
      </div>
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
      inProgress: [],
      selectedMatch: null,
      selectedResult: null
    }
  },
  mounted: async function () {
    await this.load()
  },
  methods: {
    load: async function () {
      const res = await this.axios.get(`tournaments/${this.id}/matches/in-progress`)
      this.inProgress = res.data
    },
    selectMatch: function (m) {
      this.selectedMatch = m
    },
    cancel: function () {
      this.selectedResult = null
      this.selectedMatch = null
    },
    setWinner: function (id, name) {
      this.selectedResult = {
        description: 'WINNER: ' + name,
        matchId: this.selectedMatch.id,
        isDraw: false,
        winnerId: id
      }
    },
    setDraw: function () {
      this.selectedResult = {
        description: 'DRAW',
        matchId: this.selectedMatch.id,
        isDraw: true,
        winnerId: -1
      }
    },
    confirmResult: async function () {
      await this.axios.post(`tournaments/${this.id}/match/set-result`, this.selectedResult)
      this.cancel()
      await this.load()
    }
  }
}
</script>

<style scoped>
.regres-content {
  margin: 8px;
}
.match {
  padding:10px;
  background-color: lightblue;
  margin-bottom: 8px;
  text-align: center;
}
.result {
  padding:10px;
  padding-top: 20px;
  padding-bottom: 20px;
  background-color: lightblue;
  margin-bottom: 12px;
  text-align: center;
}
</style>
