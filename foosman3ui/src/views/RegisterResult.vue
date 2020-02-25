<template>
  <div class="regres-content">
    <template v-if="selectedResult">
      <p>
        Please confirm result:
      </p>
      <p>
        {{ selectedResult.description }}
      </p>
      <div class="result">
        <b>YES</b>
      </div>
      <div class="result">
        <b>Cancel</b>
      </div>
    </template>
    <template v-else-if="selectedMatch">
      <p>
        <b>Choose result..</b>
      </p>
      <div class="result">
        {{ selectedMatch.team1Name }} won!
      </div>
      <div class="result">
        {{ selectedMatch.team2Name }} won!
      </div>
      <div class="result">
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
      const res = await this.axios.get(`http://localhost:1323/tournaments/${this.id}/matches/in-progress`)
      this.inProgress = res.data
    },
    selectMatch: function (m) {
      this.selectedMatch = m
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
