<template>
  <div v-if="tournament.state === 2">
    <p>All group matches are done, start elimination matches.</p>
    <form class="pure-form">
      <fieldset>
        <label for="elimTeamCount">Play-off size:</label>
        <select if="elimTeamCount" v-model="includeTeamCount">
          <option value="2" v-show="groupCount <= 2">Final only (2 teams, 1 tier, 1 match)</option>
          <option value="4" v-show="teamCount > 4 && groupCount <= 4">Semifinals (4 teams, 2 tiers, 3 matches)</option>
          <option value="8" v-show="teamCount > 8 && groupCount <= 8">Quarterfinals (8 teams, 3 tiers, 7 matches)</option>
          <option value="16" v-show="teamCount > 16 && groupCount <= 16">Last sixteen (16 teams, 4 tiers, 15 matches)</option>
          <option value="32" v-show="teamCount > 32 && groupCount <= 32">Last thirty two (32 teams, 5 tiers, 31 matches)</option>
        </select>
        <button
          type="button"
          class="pure-button pure-button-primary"
          @click="start"
          style="margin-left:10px;">
          Generate matches
        </button>

      </fieldset>
    </form>
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
      includeTeamCount: undefined
    }
  },
  computed: {
    teamCount: function () {
      return this.tournament.teams.length
    },
    groupCount: function () {
      return this.tournament.groups.length
    }
  },
  methods: {
    start: async function () {
      if (this.includeTeamCount === undefined) {
        alert('Please specify number of matches')
        return
      }
      const res = await this.axios.post(`tournaments/${this.tournament.id}/start-elimination`, {
        teamCount: parseInt(this.includeTeamCount)
      })
      if (res !== undefined) {
        alert('Elimination started')
        this.$router.push({ path: `/tournament/${this.tournament.id}` })
      }
    }
  }
}
</script>

<style scoped>
label {
  padding-left:6px;
  padding-right: 6px;

}
</style>
