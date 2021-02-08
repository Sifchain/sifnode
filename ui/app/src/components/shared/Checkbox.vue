<template>
  <div  @click="onClickButton()" class="checkbox-container">
    <div id="sdf" :data-checked="checked" />
    <label for="sdf">
      {{toggleLabel}}
      <span></span>
    </label>
  </div>
</template>

<script>
import { computed } from "@vue/reactivity";
export default {
  props: {
    toggleLabel: { type: String, default: null },
    checked: { type: Boolean, default: true },
  },
  emits: [
    "handleToggle",
  ],
  setup(props) {
    const toggled = computed(() => {
      return props.checked
    });
    return {
      toggled
    }
  },
	methods: {
    onClickButton (event) {
      this.$emit('clicked', '')
    },
	},


}
</script>

<style lang="scss" scoped>
  .checkbox-container {
    display: flex;
    justify-content: flex-end;
  }
  div#sdf {
    display: none;
    + label {
      display: flex;
      align-items: center;
      color: #818181;
      cursor: pointer;
      font-size: 12px;
      font-weight: 600;
      span {
        display: block;
        width: 32px;
        height: 15.84px;
        margin: 0 0 0 5.328px;
        border-radius: 1000px;
        background: rgba(129, 129, 129, 0.5);
        position: relative;
        transition: background 0.2s ease-in;

        &:after {
          display: block;
          content: "";
          width: 12px;
          height: 12px;
          background: #fff;
          border-radius: 50%;
          position: absolute;
          left: 2px;
          top: 1.8px;
          transition: left 0.2s ease-in;
        }
      }

      &:before,
      &:after {
        font-size: 20px;
        font-weight: 600;
        padding-top: 3px;
        transition: opacity 0.2s ease-in;
      }

    }

    &[data-checked=true] + label {
      span {
        background: rgb(202, 169, 58);
        &:after {
          left: 18.5px;
        }
      }

      &:before {
        opacity: 0.5;
      }

      &:after {
        opacity: 1;
      }
    }
  }

</style>
