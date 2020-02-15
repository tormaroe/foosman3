<template>
  <div v-if="tournament">
    <div><h1>{{tournament.name}}</h1></div>
    <template v-for="g in tournament.groups">
      <table class="pure-table pure-table-horizontal" style="width:99%;margin-bottom:15px;" :key="g.id">
        <thead>
          <tr>
            <th style="text-align:center">
              <a :href="'#/group/' + g.id">{{g.name}}</a>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in groupTeams(g.id)" :key="t.id">
            <td>
              <a :href="'#/team/' + t.id">{{t.name}}</a>
            </td>
          </tr>
        </tbody>
      </table>
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
      tournament: null
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
      const res = await this.axios.get(`http://localhost:1323/tournaments/${this.id}`)
      this.tournament = res.data
    },
    groupTeams: function (gID) {
      return this.tournament.teams.filter(t => t.groupId === gID)
    }
  }
}
</script>
