<template>
  <div class="layout">
    <Panel class="swap-panel">
      <PanelNav />
      <div v-if="backLink">
        <router-link class="back-link" :to="backLink">◀</router-link>
      </div>
      <slot></slot>
    </Panel>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Panel from "@/components/shared/Panel.vue";
import PanelNav from "@/components/navigation/PanelNav.vue";
export default defineComponent({
  components: { Panel, PanelNav },
  props: {
    backLink: String,
  },
});
</script>

<style scoped>
.back-link {
  display: block;
  width: 1rem;
  height: 1rem;
  text-decoration: none;
  color: #333;
  margin-bottom: 1rem;
}
.back-link:hover {
  background: #aaa;
  border-radius: 1rem;
}
.layout {
  background: #f9f9f9;
  box-sizing: border-box;
  padding-top: 3rem;
  width: 100%;
  height: calc(100vh - 4rem); /* TODO: header height */
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: flex-start;
}
</style>