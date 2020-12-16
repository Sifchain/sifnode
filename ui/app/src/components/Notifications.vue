<script lang="ts">
import { computed, defineComponent } from "vue";
import { ref, Ref } from "@vue/reactivity"; /* eslint-disable-line */
import { useCore } from "@/hooks/useCore";

export default defineComponent({
  name: "Notifications",
  components: {},
  setup() {

    const { store, actions } = useCore()
    // this will have to generate a new 
    // @ts-ignore ??
    const notifications = computed(() => store.notifications)

    // function 

    return {
      notifications
    };
  },
});


</script>

<template>
  <div v-if="notifications.length > 0" 
    class="notifications-container">
    <div 
      v-for="item in notifications" 
      v-bind:key="item.message"
      class="notification"
      v-bind:class="item.type"
      >
      {{item.message}}
    </div>
  </div>
</template>

<style lang="scss" scoped>
.notifications-container {
  position: fixed;
  bottom: 64px;
  right: 16px;
  width: 200px;
  height: 20px
}
.notification {
  background: white;
  padding: 16px;
  &.error {
    border: 1px solid red
  }
}
</style>