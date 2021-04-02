import Icon from "./Icon.vue";

export default {
  title: "Icon",
  component: Icon,
};

const Template = (args: any) => ({
  props: ["icon"],
  components: { Icon },
  setup() {
    return { args };
  },
  template: '<icon :icon="icon" v-bind="args" />',
});

export const InfoBoxBlack = Template.bind({});
InfoBoxBlack.args = {
  icon: "info-box-black",
};

export const Back = Template.bind({});
Back.args = {
  icon: "back",
};
