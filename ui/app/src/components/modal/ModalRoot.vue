<template>
  <Modal :isOpen="!!component" :title="title" @close="handleClose">
    <component :is="component" @close="handleClose" v-bind="props" />
  </Modal>
</template>

<script lang="ts">
// https://medium.com/js-dojo/vue-js-manage-your-modal-window-s-effortlessly-using-eventbus-518977195eed
import { onMounted, onUnmounted } from "vue";
import { ref } from "@vue/reactivity";
import { ModalBus } from "./ModalBus";
import Modal from "./Modal.vue";

export default {
  setup() {
    const component = ref<unknown>(null);
    const title = ref<string>("");
    const props = ref<unknown>(null);

    function handleClose() {
      component.value = null;
    }

    function handleKeyup(e: { keyCode: number }) {
      if (e.keyCode === 27) handleClose();
    }

    onMounted(() => {
      ModalBus.on(
        "open",
        (payload: { component: unknown; title: string; props: unknown }) => {
          component.value = payload.component;
          title.value = payload.title;
          props.value = payload.props;
        }
      );
      document.addEventListener("keyup", handleKeyup);
    });

    onUnmounted(() => {
      document.removeEventListener("keyup", handleKeyup);
    });
    return { component, handleClose, title, props };
  },

  components: { Modal },
};
</script>