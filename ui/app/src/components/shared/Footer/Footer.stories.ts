import Footer from "./Footer.vue";
import { withDesign } from "storybook-addon-designs";
export default {
  title: "Footer",
  component: Footer,

  decorators: [withDesign],
};

const Template = (args: any) => ({
  props: ["icon"],
  components: { Footer },
  setup() {
    return { args };
  },
  template: `
  <div>
    <div style="position: fixed; left: 0; top: 0; width: 100vw; height: 100%; height: 150px;">
    <div style="position: relative; width: 100%; height: 155px; background-image: url(https://dex.sifchain.finance/img/World_Background_opt.d3b2323b.jpg);">
      <Footer>This is an eaxample title</Footer>
      </div>
    </div>
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
