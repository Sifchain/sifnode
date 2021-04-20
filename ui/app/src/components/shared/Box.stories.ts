import Box from "./Box.vue";
import { withDesign } from "storybook-addon-designs";
export default {
  title: "Box",
  component: Box,

  decorators: [withDesign],
};

const Template = (args: any) => ({
  props: ["icon"],
  components: { Box },
  setup() {
    return { args };
  },
  template:
    "<div style='background: gray; padding: 20px; width: 400px;'><box><span>This is an example title</span></box></div>",
});

export const Basic = Template.bind({});
(Basic as any).args = {};

(Basic as any).parameters = {
  design: {
    type: "figma",
    url:
      "https://www.figma.com/file/gcSOKvZrSNKmvFDFMqrbTt/Sifchain?node-id=59%3A61",
  },
};
