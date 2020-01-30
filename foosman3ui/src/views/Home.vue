<template>
  <div class="home">
    <HelloWorld msg="Welcome to Foosman3"/>
  </div>
</template>

<script>
// @ is an alias to /src
import HelloWorld from '@/components/HelloWorld.vue'

export default {
  name: 'home',
  components: {
    HelloWorld
  }
}
</script>

<style scoped>
.home {
  margin: 8px;
}
</style>
