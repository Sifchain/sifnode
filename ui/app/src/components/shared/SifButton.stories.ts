import SifButton from "./SifButton.vue";

export default {
  title: "SifButton",
  component: SifButton,
};

const Template = (args: any) => ({
  props: [],
  components: { SifButton },
  setup() {
    return { args };
  },
  template: '<sif-button :icon="icon" v-bind="args">Click Me</sif-button>',
});

export const Primary = Template.bind({});
(Primary as any).args = {};
