import Button from "./Button.vue";

export default {
  title: "Button",
  component: Button,
};

const Template = (args: any) => ({
  props: [],
  components: { Button },
  setup() {
    return { args };
  },
  template: '<button :icon="icon" v-bind="args">Click Me</button>',
});

export const Primary = Template.bind({});
Primary.args = {};
