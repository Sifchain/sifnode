<template>
  <div class="layout">
    <Panel dark class="swap-panel">
      <template v-slot:header>
        <PanelNav />
      </template>
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

<style lang="scss" scoped>
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
  background: url("../../assets/World_Background_opt.jpg");
  background-size: cover;
  background-position: bottom center;
  box-sizing: border-box;
  padding-top: $header_height;
  width: 100%;
  height: 100vh; /* TODO: header height */
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
</style>