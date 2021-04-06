import SubHeading from "./SubHeading.vue";
import { withDesign } from "storybook-addon-designs";
export default {
  title: "SubHeading",
  component: SubHeading,

  decorators: [withDesign],
};

const Template = (args: any) => ({
  props: ["icon"],
  components: { SubHeading },
  setup() {
    return { args };
  },
  template: "<div><sub-heading>This is an example title</sub-heading></div>",
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
