<template>
  <div class="panel-nav">
    <button @click="$router.push('swap')">swap</button>
    <button @click="$router.push('pool')">pool</button>
  </div>
</template>
<style scoped>
.panel-nav {
  margin-bottom: 1rem;
  display: flex;
}
</style>