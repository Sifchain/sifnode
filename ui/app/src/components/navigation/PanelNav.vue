<template>
  <div class="panel-nav">
    <router-link class="link" to="/swap">Swap</router-link>
    <router-link class="link" to="/pool">Pool</router-link>
  </div>
</template>
<style scoped>
.panel-nav {
  margin-bottom: 1rem;
  display: flex;
}
.link {
  color: rgb(10, 10, 10);
  height: 3rem;
  width: 3rem;
  border-radius: 1rem;
  background: #aaa;
  margin-right: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.router-link-active {
  color: rgb(95, 95, 95);
}
</style>