import Footer from "./Footer.vue";

export default {
  title: "Footer",
  component: Footer,
};

const Template = (args: any) => ({
  props: ["footer"],
  components: { Footer },
  setup() {
    return { args };
  },
  template: '<footer :footer="footer" v-bind="args" />',
});

export const Footer = Template.bind({});
Footer.args = {
  footer: "info-box-black",
};
