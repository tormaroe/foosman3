<template>
  <div v-if="tournament">
    <button
      v-show="mode !== 'default'"
      type="button"
      @click="mode = 'default'"
      class="pure-button"
      style="float: right">Cancel</button>
    <button
      v-show="mode !== 'groups'"
      type="button"
      @click="mode = 'groups'"
      class="pure-button pure-button-primary"
      style="float: right">Groups</button>
    <h3>Tournament: {{ tournament.name }}</h3>
    <hr>

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
                <button @click="selectTeam(t)" class="pure-button">→</button>
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

export default {
  components: {
    Groups
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
      const res = await this.axios.get(`http://localhost:1323/tournaments/${this.id}`)
      let tmp = res.data
      if (tmp.teams === null) {
        tmp.teams = []
      }
      this.tournament = tmp
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
        await this.axios.patch('http://localhost:1323/tournaments/teams', this.editTeam)
      } else {
        await this.axios.post(`http://localhost:1323/tournaments/${this.id}/teams`, this.editTeam)
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
    groupsSave: async function (groups) {
      await this.load()
      this.mode = 'default'
    }
  }
}
</script>
