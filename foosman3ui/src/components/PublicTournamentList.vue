<template>
  <div class="hello">
    <img class="center" style="height:200px;width:200px;" src="../assets/foosmanlogo.png" />
        <table class="pure-table pure-table-horizontal" style="width:99%">
          <thead>
            <tr>
              <th>Tournaments</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="t in tournaments" :key="t.id">
              <td><a :href="'#/tournament/' + t.id">{{ t.name }}</a></td>
            </tr>
          </tbody>
        </table>
  </div>
</template>

<script>
export default {
  name: 'PublicTournamentList',
  data: function () {
    return {
      tournaments: []
    }
  },
  mounted: function () {
    this.load()
  },
  methods: {
    load: function () {
      const self = this
      this.axios.get('tournaments').then(function (res) {
        self.tournaments = res.data
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
.center {
  display: block;
  margin-left: auto;
  margin-right: auto;
  width: 50%;
}
</style>
