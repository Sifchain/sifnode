<script>
export default {
  props: {
    color: "primary" | "info" | "danger",
  },
  setup(props) {
    const { color } = props;
    return { classes: { [color]: true, pill: true } };
  },
};
</script>
<template>
  <div class="container">
    <div :class="classes">
      <slot></slot>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.container {
}
.pill {
  border: 2px solid #d4b553;
  color: #d4b553;
  font-weight: 600;
  border-radius: 20px;
  padding: 0px 8px 0px 4px;
  display: inline-block;
  background: #000;
  display: flex;
  cursor: pointer;
}
.primary {
  border: 2px solid #d4b553;
  color: #d4b553;
}
.info {
  border: 2px solid #3b7fba;
  color: #3b7fba;
}
.danger {
  border: 2px solid #f05089;
  color: #f05089;
}
.success {
  border: 2px solid #97d937;
  color: #97d937;
}
</style>
