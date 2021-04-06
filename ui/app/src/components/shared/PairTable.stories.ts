import PairTable from "./PairTable.vue";
import { withDesign } from "storybook-addon-designs";
export default {
  title: "PairTable",
  component: PairTable,

  decorators: [withDesign],
};

const Template = (args: any) => ({
  props: ["items"],
  components: { PairTable },
  setup() {
    const { items } = args;
    return {
      items,
      args,
    };
  },
  template: '<div><pair-table :items="items" /></div>',
});

export const Basic = Template.bind({});
Basic.args = {
  items: [
    { key: "Your Multiplier Date", value: "12 Aug 2020" },
    { key: "Your Current Multiplier", value: "1.2x" },
  ],
};

Basic.parameters = {
  design: {
    type: "figma",
    url: "",
  },
};
