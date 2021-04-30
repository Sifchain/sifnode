import PanelNavItem from "./PanelNavItem.vue";

export default {
  title: "PanelNavItem",
  component: PanelNavItem,
};

const Template = (args: any) => ({
  props: [],
  components: { PanelNavItem },
  setup() {
    return { args };
  },
  template: `
  <div style="position: relative;width: 60px;background: #272727;padding: 20px;">
    <PanelNavItem :icon="icon" v-bind="args">Click Me</PanelNavItem>
  </div>
  `,
});

export const Primary = Template.bind({});
Primary.args = {
  icon: "circle-arrows",
  color: "pink",
  label: "SWAP",
};
