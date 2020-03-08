<template>
  <div id="app">
    <div id="nav">
      <router-link to="/"><home-icon /> Home</router-link> |
      <router-link to="/admin"><cogs-icon /> Admin</router-link>
    </div>
    <router-view/>
    <hr class="hr-footer">
    <div class="footer">
      <p>
        <b>Foosman v3</b>
        Copyright 2020 Torbjørn Marø
      </p>
      <p>
        Powered by Go, Vue, Docker, SQLite, Linode, and lots of love for programming and foosball.
      </p>
    </div>
  </div>
</template>

<script>
import HomeIcon from 'vue-material-design-icons/Home.vue'
import CogsIcon from 'vue-material-design-icons/Cogs.vue'

export default {
  components: {
    HomeIcon,
    CogsIcon
  }
}
</script>

<style>
body {
  background-color: #FFFFFF;
}

#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #000000;
}

.hr-footer {
  border: solid 1px #E9E9E9;
  margin-left:10px;
  margin-right:10px;
  margin-top: 50px;
  margin-bottom:10px;
}

.footer {
  margin-left:10px;
  margin-right:10px;
  color: #969696;
  font-size: 12px;
}

th {
  background-color: #E9E9E9;
}

tr {
  background-color: #F9F9F9;
}

#nav {
  padding: 20px;
  background-color: #E9E9E9;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #47A8DE;
}

/* unvisited link */
a:link {
  color: #47A8DE;
  text-decoration: none;
}

/* visited link */
a:visited {
  color: #47a8de;
}

/* mouse over link */
a:hover {
  color: hotpink;
}

/* selected link */
a:active {
  color: #47A8DE;
}
</style>
