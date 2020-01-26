<template>
  <div id="admincontainer">
    <h1>Admin</h1>
    <router-view />
  </div>
</template>

<script>

export default {
  mounted: function () {}
}
</script>

<style>
#admincontainer {
  margin:8px;
}

.column {
}

legend {
  font-weight: bold;
}
</style>
