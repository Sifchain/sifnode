import Menu from "./Menu.vue";
import { withDesign } from "storybook-addon-designs";
export default {
  title: "Menu",
  component: Menu,

  decorators: [withDesign],
};

const Template = (args: any) => ({
  props: ["icon"],
  components: { Menu },
  setup() {
    return { args };
  },
  template: `
  <div>
      <Menu>This is an eaxample title</Menu>
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
