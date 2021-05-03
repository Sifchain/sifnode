import Pill from "./Pill.vue";
import { withDesign } from "storybook-addon-designs";
export default {
  title: "Pill",
  component: Pill,

  decorators: [withDesign],
};

const Template = (args: any) => ({
  props: ["icon"],
  components: { Pill },
  setup() {
    return { args };
  },
  template: `
  <div>
      <Pill v-bind="args">TVL: $22,919,930</Pill>
  </div>`,
});

export const Default = Template.bind({});

Default.args = {
  color: {
    control: {
      type: "select",
      options: ["primary", "secondary"],
    },
  },
};

Default.parameters = {
  design: {
    type: "figma",
    url:
      "https://www.figma.com/file/gcSOKvZrSNKmvFDFMqrbTt/Sifchain?node-id=59%3A61",
  },
};
