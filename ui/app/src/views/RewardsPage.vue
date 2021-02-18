<script lang="ts">
import { defineComponent } from "vue";
import Layout from "@/components/layout/Layout.vue";

export default defineComponent({
  components: {
    Layout,
  },
  data() {
    return {
      poolData: {},
    };
  },
  async mounted() {
    const data = await fetch(
      "https://vtdbgplqd6.execute-api.us-west-2.amazonaws.com/default/liqvalrewards"
    );
    const json = await data.json();
    console.log(json)
    this.poolData = json.body;
  },
});
</script>

<template>
  <Layout>
    
      <template slot="header">
          hi
      </template>
  </Layout>
</template>

<style scoped lang="scss">
// .layout {
//   background: url("../assets/World_Background_opt.jpg");
//   background-size: cover;
//   background-position: bottom center;
//   box-sizing: border-box;
//   padding-top: $header_height;
//   padding-right: 32px;
//   padding-left: 32px;
//   width: 100%;
//   height: 100vh;
//   display: flex;
//   flex-direction: column;
//   align-items: center;
// }
</style>