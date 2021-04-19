import Copy from "./Copy.vue";
import { withDesign } from "storybook-addon-designs";
export default {
  title: "Copy",
  component: Copy,

  decorators: [withDesign],
};

const Template = (args: any) => ({
  props: ["icon"],
  components: { Copy },
  setup() {
    return { args };
  },
  template: "<div><copy>This is an example title</copy></div>",
});

export const Body = Template.bind({});
Body.args = {};

Body.parameters = {
  design: {
    type: "figma",
    url:
      "https://www.figma.com/file/gcSOKvZrSNKmvFDFMqrbTt/Sifchain?node-id=59%3A61",
  },
};
