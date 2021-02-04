<template>
  <div class="text">
    <p>{{ title }}</p>
    <slot></slot>
    <br />
    <p class="sub">{{ sub }}<slot name="sub"></slot></p>
  </div>
</template>


<script>
export default {
  props: ["title", "sub"],
};
</script>

<style lang="scss" scoped>
.sub {
  font-weight: normal;
  font-size: $fs_sm;
}
.anchor {
  color: $c_black;
}
</style>