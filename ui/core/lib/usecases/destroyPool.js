// No async means this cannot use the store or remote apis.
function renderDestroyPool(isAdmin) {
    // ...
    return {};
}
export default ({ api, store }) => ({
    // Render helpers that are business logic
    renderDestroyPool,
    // Command and effect usecases
    async destroyPool() {
        //
    },
});
//# sourceMappingURL=destroyPool.js.map