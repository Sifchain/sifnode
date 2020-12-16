<script lang="ts">
import { computed, defineComponent } from "vue";
import { ref, Ref } from "@vue/reactivity"; /* eslint-disable-line */
import { useCore } from "@/hooks/useCore";

export default defineComponent({
  name: "Notifications",
  components: {},
  setup() {

    const { store, actions } = useCore()
    const notifications = computed(() => store.notifications)

    return {
      notifications
    };
  },
});
</script>

<template>
  <div class="notifications">
    {{JSON.stringify(notifications)}}
  </div>
</template>

<style scoped>
.notifications {
  position: fixed;
  bottom: 16px;
  right: 16px;
  background: white;
  width: 200px;
  height: 20px
}
</style>