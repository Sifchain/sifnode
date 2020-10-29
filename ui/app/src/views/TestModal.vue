<template>
  <Modal
    ><template v-slot:activator="{ openModal }"
      ><button @click="openModal">Go</button></template
    >
    <h1>{{ foo }}</h1>
    <button @click="setBar">bar</button>
  </Modal>
</template>

<script lang="ts">
import { defineComponent, onMounted, ref } from "vue";
import Modal from "@/components/shared/Modal.vue";
export default defineComponent({
  components: { Modal },
  setup() {
    const foo = ref("Hello");
    onMounted(() => {
      foo.value = "hello";
    });
    function setBar() {
      foo.value = "bar";
    }
    return { foo, setBar };
  },
});
</script>