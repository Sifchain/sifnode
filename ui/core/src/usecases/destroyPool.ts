import { Context } from ".";

// No async means this cannot use the store or remote apis.
function renderDestroyPool(
  isAdmin: boolean
): {
  destroyPoolButtonAvailable: boolean;
} {
  // ...
  return {} as any;
}

export default ({ api, store }: Context) => ({
  // Render helpers that are business logic
  renderDestroyPool,

  // Command and effect usecases
  async destroyPool() {
    //
  },
});
