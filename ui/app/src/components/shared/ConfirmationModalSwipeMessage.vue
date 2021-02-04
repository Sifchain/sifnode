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