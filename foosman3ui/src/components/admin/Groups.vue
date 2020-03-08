<template>
  <div>
    <div style="margin-bottom:10px;">
      Group count: {{ groups.length }}
      <button class="pure-button" @click="decGroupCount()">-</button>
      <button class="pure-button" @click="incGroupCount()">+</button>
      |
      <input type="range" min="50" max="1000" v-model="groupBoxWidth" class="slider">
      |
      <button class="pure-button" @click="distribute()">Auto distribute</button>
      |
      <button class="pure-button pure-button-primary" @click="save()">Save</button>
    </div>
    <div v-drag-and-drop:options="options" class="drag-wrapper">
      <div>
        <div>
          <b>Ungrouped</b><br>
          Count: {{availableTeams.length}}
        </div>
        <ul class="pink" :style="groupBoxStyle">
          <li
            v-for="t in availableTeams"
            :key="t.id"
            :data-id="t.id">
            <label v-text="t.name"></label>
          </li>
        </ul>
      </div>
      <div
        v-for="(g, idx) in groups"
        :key="g.id"
      >
      <div>
        <b>{{g.groupName}}</b><br>
        Count: {{g.teams.length}}
      </div>
      <ul
        :style="groupBoxStyle"
        :data-idx="idx"
        :data-group="true"
      >
        <li v-for="t in g.teams" :key="t.id" :data-id="t.id">
          <label v-text="t.name"></label>
        </li>
      </ul>
      </div>
    </div>
  </div>
</template>

<script>

const numberToLetter = function (n) {
  return String.fromCharCode(n + 64)
}

export default {
  props: {
    tournament: {
      type: Object,
      required: true
    }
  },
  data: function () {
    return {
      availableTeams: [],
      groups: [],
      groupBoxWidth: 200,
      options: {
        // dropzoneSelector: 'ul',
        // draggableSelector: 'li',
        // handlerSelector: null,
        // reactivityEnabled: true,
        // multipleDropzonesItemsDraggingEnabled: true,
        // showDropzoneAreas: true,
        // onDrop: function (event) {},
        // onDragstart: function (event) {},
        // onDragenter: function (event) {},
        // onDragover: function (event) {},
        onDragend: this.onDragend
      }
    }
  },
  watch: {
    tournament: {
      immediate: true,
      handler (value) {
        this.availableTeams = JSON.parse(JSON.stringify(value.teams))
        value.groups.forEach(g => {
          var group = {
            groupName: g.name,
            teams: this._.remove(this.availableTeams, t => t.groupId === g.id)
          }
          this.groups.push(group)
        })
      }
    }
  },
  computed: {
    groupBoxStyle: function () {
      return `width:${this.groupBoxWidth}px;`
    }
  },
  methods: {
    incGroupCount: function () {
      const groupNumber = this.groups.length + 1
      this.groups.push({
        // groupNumber: groupNumber,
        groupName: 'Group ' + numberToLetter(groupNumber), // TODO: Potensial issue!!
        teams: []
      })
    },
    decGroupCount: function () {
      if (this.groups.length === 0) {
        return
      }
      const g = this.groups.pop()
      this.availableTeams.push(...g.teams)
    },
    onGroupsChange: function (e) {
    },
    onDragend: function (e) {
      if (!e.droptarget) return

      e.items.forEach(element => {
        let team = null
        const elementID = parseInt(element.dataset.id)

        if (e.owner.dataset.group) {
          const teamIdx = this.groups[e.owner.dataset.idx].teams.findIndex(t => t.id === elementID)
          team = this.groups[e.owner.dataset.idx].teams.splice(teamIdx, 1)[0]
        } else {
          const teamIdx = this.availableTeams.findIndex(t => t.id === elementID)
          team = this.availableTeams.splice(teamIdx, 1)[0]
        }

        if (e.droptarget.dataset.group) {
          this.groups[e.droptarget.dataset.idx].teams.push(team)
        } else {
          this.availableTeams.push(team)
        }
      })
    },
    distribute: function () {
      if (this.groups.length === 0) return

      let gIdx = 0
      while (this.availableTeams.length > 0) {
        let team = this.availableTeams.shift()
        this.groups[gIdx].teams.push(team)
        gIdx++
        if (gIdx >= this.groups.length) {
          gIdx = 0
        }
      }
    },
    save: async function () {
      const groupDtos = this.groups.map(g => ({
        name: g.groupName,
        teams: g.teams.map(t => t.id)
      }))
      await this.axios.post(`tournaments/${this.tournament.id}/groups`, groupDtos)
      this.$emit('save')
    }
  }
}
</script>

<style scoped>

.pink {
  background: pink;
}

.drag-wrapper {
  display: flex;
}

ul {
  display: flex;
  flex-direction: column;
  padding: 3px !important;
  min-height: 50vh;
  width: 210px;
  float:left;
  list-style-type:none;
  overflow-y:auto;
  border:1px solid #888;
  border-radius:0.0em;
  background:#8adccc;
  color:#555;
  margin-right: 5px;
}

/* drop target state */
ul[aria-dropeffect="move"] {
  border-color:#68b;
  background:#fff;
}

/* drop target focus and dragover state */
ul[aria-dropeffect="move"]:focus,
ul[aria-dropeffect="move"].dragover
{
  outline:none;
  box-shadow:0 0 0 1px #fff, 0 0 0 3px #68b;
}

/* draggable items */
li {
  display:block;
  list-style-type:none;
  margin:0 0 2px 0;
  padding:0.2em 0.4em;
  border-radius:0.0em;
  line-height:1.3;
}

li:hover {
  box-shadow:0 0 0 2px #68b, inset 0 0 0 1px #ddd;
}

/* items focus state */
li:focus
{
  outline:none;
  box-shadow:0 0 0 2px #68b, inset 0 0 0 1px #ddd;
}

/* items grabbed state */
li[aria-grabbed="true"]
{
  background:#5cc1a6;
  color:#fff;
}

@keyframes nodeInserted {
    from { opacity: 0.2; }
    to { opacity: 0.8; }
}

.item-dropzone-area {
    height: 2rem;
    background: #888;
    opacity: 0.8;
    animation-duration: 0.5s;
    animation-name: nodeInserted;
}

</style>
