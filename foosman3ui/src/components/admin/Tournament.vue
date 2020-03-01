<template>
  <div v-if="tournament">
    <button
      type="button"
      @click="deleteTournament"
      class="pure-button pure-button-danger"
      style="float: right;margin-left:5px;">DELETE</button>
    <button
      type="button"
      @click="resetTournament"
      class="pure-button pure-button-danger"
      style="float: right;margin-left:5px;">Reset</button>
    <button
      v-show="mode !== 'default'"
      type="button"
      @click="mode = 'default'"
      class="pure-button pure-button-primary"
      style="float: right;margin-left:5px;">Teams</button>
    <button
      type="button"
      @click="startTournament"
      class="pure-button pure-button-primary"
      style="float: right;margin-left:5px;">START TOURNAMENT</button>
    <button
      v-show="mode !== 'matches'"
      type="button"
      @click="mode = 'matches'"
      class="pure-button pure-button-primary"
      style="float: right;margin-left:5px;">Matches</button>
    <button
      v-show="mode !== 'groups'"
      type="button"
      @click="mode = 'groups'"
      class="pure-button pure-button-primary"
      style="float: right">Groups</button>
    <h3>Tournament: {{ tournament.name }}</h3>
    <hr>

    <form class="pure-form">
      <fieldset>
        <label for="tName">Name:</label>
        <input id="tName" type="text" v-model="tournament.name">
        <label for="tTableCount">Table#:</label>
        <input id="tTableCount" type="number" v-model="tournament.tableCount">
        <button
          type="button"
          @click="updateTournament"
          class="pure-button pure-button-primary"
          style="margin-left:10px;">
            Update
        </button>
      </fieldset>
    </form>

    <hr>
    <matches
      v-if="mode === 'matches'"
      :tournament-id="tournament.id"
      />
    <groups
      v-if="mode === 'groups'"
      :tournament="tournament"
      @save="groupsSave"
      />

    <div v-show="mode === 'default'" class="pure-g">
      <div class="pure-u-2-3 column">
        <p>{{ tournament.teams.length }} teams</p>
        <table class="pure-table pure-table-horizontal" style="width:99%">
          <thead>
            <tr>
              <th>Name</th>
              <th>Player 1</th>
              <th>Player 2</th>
              <th>Player 3</th>
              <th>&nbsp;</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="t in tournament.teams" :key="t.id">
              <td>{{ t.name }}</td>
              <td>{{ t.player1 }}</td>
              <td>{{ t.player2 }}</td>
              <td>{{ t.player3 }}</td>
              <td style="text-align:right">
                <button @click="selectTeam(t)" class="pure-button" title="Edit team">→</button>
                &nbsp;
                <button @click="deleteTeam(t)" class="pure-button" title="Delete team">x</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="pure-u-1-3 column">
        <form class="pure-form pure-form-stacked">
          <fieldset>
            <legend>{{editForm.legend}}</legend>
            <label>Team name</label>
            <input type="text" placeholder="" v-model="editTeam.name">
            <label>Player #1</label>
            <input type="text" placeholder="" v-model="editTeam.player1">
            <label>Player #2</label>
            <input type="text" placeholder="" v-model="editTeam.player2">
            <label>Player #3</label>
            <input type="text" placeholder="" v-model="editTeam.player3">
            <button
              type="button"
              @click="saveTeam"
              class="pure-button pure-button-primary">
                {{editForm.button}}
            </button>
            <button
              type="button"
              @click="cancelEdit"
              class="pure-button">
                Cancel
            </button>
          </fieldset>
        </form>
      </div>
    </div>

  </div>
</template>

<script>
import Groups from './Groups.vue'
import Matches from './Matches.vue'

export default {
  components: {
    Groups,
    Matches
  },
  props: {
    id: {
      type: [Number, String],
      required: true
    }
  },
  data: function () {
    return {
      tournament: null,
      mode: 'default',
      editForm: {
        legend: 'Add a team',
        button: 'Add team'
      },
      editTeam: {
        id: 0,
        name: '',
        player1: '',
        player2: '',
        player3: ''
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
      const res = await this.axios.get(`tournaments/${this.id}`)
      let tmp = res.data
      if (tmp.teams === null) {
        tmp.teams = []
      }
      this.tournament = tmp
    },
    resetTournament: async function () {
      if (window.confirm('This will only work if there are no match results yet. Continue?')) {
        await this.axios.post(`tournaments/${this.id}/reset`)
        await this.load()
      }
    },
    deleteTournament: async function () {
      if (window.confirm('ARE YOU SURE?')) {
        await this.axios.delete(`tournaments/${this.id}`)
        this.$router.push({ path: '/admin' })
      }
    },
    updateTournament: async function () {
      this.tournament.tableCount = parseInt(this.tournament.tableCount)
      await this.axios.patch(`tournaments`, this.tournament)
      await this.load()
    },
    selectTeam: function (t) {
      this.editTeam.id = t.id
      this.editTeam.name = t.name
      this.editTeam.player1 = t.player1
      this.editTeam.player2 = t.player2
      this.editTeam.player3 = t.player3
      this.editForm.legend = 'Edit team'
      this.editForm.button = 'Update'
    },
    cancelEdit: function () {
      this.clearEditTeam()
    },
    saveTeam: async function () {
      // TODO: Validation
      if (this.editTeam.id > 0) {
        await this.axios.patch('tournaments/teams', this.editTeam)
      } else {
        await this.axios.post(`tournaments/${this.id}/teams`, this.editTeam)
      }
      this.clearEditTeam()
      await this.load()
    },
    clearEditTeam: function () {
      this.editTeam = {
        id: 0,
        name: '',
        player1: '',
        player2: '',
        player3: ''
      }
      this.editForm.legend = 'Add a team'
      this.editForm.button = 'Add team'
    },
    deleteTeam: async function (t) {
      if (confirm(`Delete ${t.name} for sure?`)) {
        await this.axios.delete(`tournaments/teams/${t.id}`)
        await this.load()
      }
    },
    groupsSave: async function (groups) {
      await this.load()
      this.mode = 'default'
    },
    startTournament: async function () {
      if (confirm('Are all teams registered and placed in groups? Then you may start the tournament..')) {
        await this.axios.post(`tournaments/${this.id}/start`)
        alert('Tournament started')
        this.$router.push({ path: `/tournament/${this.id}` })
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
.pure-button-danger {
  background-color: red;
  color: white;
}
</style>
