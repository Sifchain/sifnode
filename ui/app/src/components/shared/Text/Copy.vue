<script>
/*
TODO - Experimental Text component
Need to figure out best way to use it
e.g. import Text from 'core/Text'
<Text.Body>Something</Text.Body>
or <Text type="heading"> etc
*/
</script>

<template>
  <div class="container">
    <slot></slot>
  </div>
</template>

<style lang="scss" scoped>
.container {
  font-size: 12px;
  line-height: 18px;
  text-align: left;
}
</style>
