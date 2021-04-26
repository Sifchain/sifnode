<script lang="ts">
import { defineComponent } from "vue";
import Layout from "@/components/layout/Layout.vue";
import PoolStatsList from "@/components/poolStats/PoolStatsList.vue";
import PoolStatsListHeader from "@/components/poolStats/PoolStatsListHeader.vue";

export default defineComponent({
  components: {
    Layout,
    PoolStatsList,
    PoolStatsListHeader,
  },
  data() {
    return {
      poolData: {},
    };
  },
  async mounted() {
    const data = await fetch(
      "https://vtdbgplqd6.execute-api.us-west-2.amazonaws.com/default/tokenstatstest",
    );
    const json = await data.json();
    this.poolData = json.body;
  },
});
</script>

<template>
  <div class="layout">
    <PoolStatsList :poolData="poolData" />
  </div>
</template>

<style scoped lang="scss">
.layout {
  background: url("../assets/World_Background_opt.jpg");
  background-size: cover;
  background-position: bottom center;
  box-sizing: border-box;
  padding-top: $header_height;
  padding-right: 32px;
  padding-left: 32px;
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
}
</style>
