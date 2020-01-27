<template>
  <div>
    <p>
      Group count: {{ groups.length }}
      <button class="pure-button" @click="decGroupCount()">-</button>
      <button class="pure-button" @click="incGroupCount()">+</button>
      |
      <button class="pure-button" @click="distribute()">Auto distribute</button>
      |
      <button class="pure-button pure-button-primary" @click="save()">Save</button>
    </p>
    <div v-drag-and-drop:options="options" class="drag-wrapper">
      <div>
        <div>
          <b>Ungrouped</b><br>
          Count: {{availableTeams.length}}
        </div>
        <ul style="background: pink">
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
        <b>Group {{g.groupLetter}}</b><br>
        Count: {{g.teams.length}}
      </div>
      <ul
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
        console.log('prop changed!')
        console.dir(value)
        this.availableTeams = JSON.parse(JSON.stringify(value.teams))
        // TODO: Move grouped teams to groups...
      }
    }
  },
  methods: {
    incGroupCount: function () {
      const groupNumber = this.groups.length + 1
      this.groups.push({
        groupNumber: groupNumber,
        groupLetter: numberToLetter(groupNumber),
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
      console.dir(e)
    },
    onDragend: function (e) {
      console.dir(e)

      if (!e.droptarget) return

      e.items.forEach(element => {
        let team = null
        const elementID = parseInt(element.dataset.id)
        console.log(`Move object id ${elementID}`)

        if (e.owner.dataset.group) {
          const teamIdx = this.groups[e.owner.dataset.idx].teams.findIndex(t => t.id === elementID)
          team = this.groups[e.owner.dataset.idx].teams.splice(teamIdx, 1)[0]
        } else {
          const teamIdx = this.availableTeams.findIndex(t => t.id === elementID)
          console.log(`teamIdx=${teamIdx}, ${this.availableTeams[teamIdx]}`)
          team = this.availableTeams.splice(teamIdx, 1)[0]
        }

        if (e.droptarget.dataset.group) {
          this.groups[e.droptarget.dataset.idx].teams.push(team)
        } else {
          this.availableTeams.push(team)
        }

        console.dir(element.dataset.id)
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
    save: function () {
      const groupDtos = this.groups.map(g => ({
        name: 'Group ' + g.groupLetter,
        teams: g.teams.map(t => t.id) 
      }))
      this.$emit('save')
    }
  }
}
</script>

<style scoped>

.drag-wrapper {
  display: flex;
  justify-content: center;
}

ul {
  display: flex;
  flex-direction: column;
  padding: 3px !important;
  min-height: 70vh;
  width: 110px;
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
