<template>
  <div id="app">
    <div id="nav">
      <router-link to="/"><home-icon /> Home</router-link> |
      <router-link to="/admin"><cogs-icon /> Admin</router-link>
    </div>
    <router-view/>
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
  background-color: #f5fff6;
}

#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #243325;
}

th {
  background-color: #97e69f;
}

#nav {
  padding: 30px;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #42b983;
}

/* unvisited link */
a:link {
  color: #3f9148;
  text-decoration: none;
}

/* visited link */
a:visited {
  color: #3f9148;
}

/* mouse over link */
a:hover {
  color: hotpink;
}

/* selected link */
a:active {
  color: #3f9148;
}
</style>
