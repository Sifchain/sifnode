import ErrorBoundary from "./ErrorBoundary.vue";
import { withDesign } from "storybook-addon-designs";
export default {
  title: "ErrorBoundary",
  component: ErrorBoundary,
  decorators: [withDesign],
};

const Template = (args: any) => ({
  props: [],
  components: { ErrorBoundary },
  setup() {
    return { args };
  },
  template: `<ErrorBoundary :error="args.error">
    <template #fallback>
    <div>
      <h1>This is fallback</h1>
      </div>
    </template>
    <template #default>
    <div>
      <h1>I am content</h1>
    </div>
    </template>
    </ErrorBoundary>`,
});

export const ShowError = Object.assign(Template.bind({}), {
  args: { error: true },
});

export const HideError = Object.assign(Template.bind({}), {
  args: { error: false },
});
