import Header from "./Header.vue";
import { withDesign } from "storybook-addon-designs";
export default {
  title: "Header",
  component: Header,

  decorators: [withDesign],
};

const Template = (args: any) => ({
  props: ["icon"],
  components: { Header },
  setup() {
    return { args };
  },
  template: `
  <div>
      <Header>This is an eaxample title</Header>
  </div>`,
});

export const Default = Template.bind({});

Default.args = {};

Default.parameters = {
  design: {
    type: "figma",
    url:
      "https://www.figma.com/file/gcSOKvZrSNKmvFDFMqrbTt/Sifchain?node-id=59%3A61",
  },
};
