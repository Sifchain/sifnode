<template>
  <div class="layout">
    <Panel dark>
      <template v-slot:header v-if="header">
        <PanelNav />
      </template>
      <div class="header" v-if="backLink || title">
        <div v-if="backLink">
          <router-link class="back-link" :to="backLink"
            ><Icon icon="back"
          /></router-link>
        </div>
        <div v-if="emitBack">
          <span @click="$emit('back')" class="back-link"
            ><Icon icon="back"
          /></span>
        </div>
        <div class="title">
          <SubHeading>{{ title }}</SubHeading>
        </div>
      </div>
      <slot></slot>
    </Panel>
    <Panel v-if="!!$slots.after" class="after">
      <slot name="after"></slot>
    </Panel>
  </div>
  <Footer />
  <div class="layout-bg" />
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Panel from "@/components/shared/Panel.vue";
import Footer from "@/components/shared/Footer/Footer.vue";
import PanelNav from "@/components/shared/PanelNav/PanelNav.vue";
import Icon from "@/components/shared/Icon.vue";
import { SubHeading } from "@/components/shared/Text";

export default defineComponent({
  components: { Panel, PanelNav, Icon, SubHeading, Footer },
  props: {
    backLink: String,
    header: { type: Boolean, default: true },
    emitBack: {
      type: Boolean,
      default: false,
    },
    title: String,
  },
});
</script>

<style lang="scss" scoped>
.layout {
  box-sizing: border-box;
  padding-top: $header_height;
  width: 100%;
  height: 100vh; /* TODO: header height */
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
.layout-bg {
  background: url("../../assets/World_Background_opt.jpg");
  background-size: cover;
  background-position: bottom center;
  // filter: blur(10px);

  z-index: -1;
  width: 100%;
  height: 100vh; /* TODO: header height */
  position: absolute;
  top: 0;
  left: 0;
}

.after {
  margin-top: 15px;
  padding: 25px;
  background: linear-gradient(180deg, $c_gray_50 0%, $c_gray_200 100%);
}

.header {
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
}
.back-link {
  text-align: left;
  display: block;
  text-decoration: none;
  position: relative;
  top: 2px;
  cursor: pointer;
}
.title {
  display: flex;
  justify-content: center;
  width: 100%;
}
</style>
