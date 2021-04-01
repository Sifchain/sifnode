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
  font-size: 16px;
  line-height: 22px;
  text-align: left;
  font-weight: 700;
  color: #343434;
}
</style>
