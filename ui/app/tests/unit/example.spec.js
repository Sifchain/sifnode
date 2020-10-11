"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const tslib_1 = require("tslib");
const test_utils_1 = require("@vue/test-utils");
const HelloWorld_vue_1 = tslib_1.__importDefault(require("@/components/HelloWorld.vue"));
describe('HelloWorld.vue', () => {
    it('renders props.msg when passed', () => {
        const msg = 'new message';
        const wrapper = test_utils_1.shallowMount(HelloWorld_vue_1.default, {
            propsData: { msg }
        });
        expect(wrapper.text()).toMatch(msg);
    });
});
//# sourceMappingURL=example.spec.js.map