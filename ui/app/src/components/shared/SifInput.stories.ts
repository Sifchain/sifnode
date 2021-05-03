import SifInput from "./SifInput.vue";

export default {
  title: "SifInput",
  component: SifInput,
};

const Template = (args: any) => ({
  props: [],
  components: { SifInput },
  setup() {
    return { args };
  },
  template: '<SifInput gold placeholder="Search name or paste address" />',
});

export const Primary = Template.bind({});
Primary.args = {};
