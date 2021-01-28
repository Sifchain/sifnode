<template>
  <div>
    <h3 class="title mb-10">{{ title }}</h3>
    <slot name="body"></slot>
    <SifButton block primary class="confirm-btn" @click="$emit('confirmed')"
      >Confirm</SifButton
    >
  </div>
</template>

<script>
import SifButton from "@/components/shared/SifButton.vue";

export default {
  components: {
    SifButton
  },

  props: {
    title: String,
  }
}
</script>

<style>

</style>