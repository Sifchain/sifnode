<script lang="ts">
import { defineComponent } from "vue";
import jsonp from "jsonp";
import queryString from "query-string";
import Icon from "@/components/shared/Icon.vue";
import SifButton from "@/components/shared/SifButton.vue";

const MAILCHIMP_URL = "https://finance.us2.list-manage.com/subscribe/post-json";

export default defineComponent({
  components: { Icon, SifButton },
  props: {},
  data() {
    return {
      active: true,
      email: "",
      submitted: false,
    };
  },
  setup() {
    return {};
  },
  methods: {
    toggleActive() {
      this.active = !this.active;
    },
    async submitEmail() {
      this.submitted = true;
      const query = `u=400787e0a5e23ec37b7b51f74&id=c1ee83387b&EMAIL=${this.email}`;
      const url = `${MAILCHIMP_URL}?${query}`;
      jsonp(url, { param: "c" });
      // TODO - Set this as a callback to jsonp
      // , (error, data) => {
      //   // TODO - add a loading spinner
      //   if (error) {
      //     // say, try again
      //   } else {
      //     // say, thank you
      //   }
      // });
    },
  },
});
</script>

<template>
  <div class="container">
    <div
      class="footer"
      :class="{
        active: active,
      }"
    >
      <div class="backdrop"></div>
      <div class="items">
        <div class="left">
          <div class="toggle-button" @click="toggleActive">
            <Icon icon="close-full-circle" />Close
          </div>
        </div>
        <div class="right">
          <div class="cta">
            <div>Sign up for Sifchain updates</div>
            <div v-if="!submitted">
              <form @submit.prevent="submitEmail">
                <input
                  v-model="email"
                  type="email"
                  name="email"
                  class="email-input"
                />
                <SifButton primary type="submit">Stay Informed</SifButton>
              </form>
            </div>
            <div v-else class="thankyou">Thank you!</div>
          </div>
          <div class="links">
            <a
              target="_blank"
              href="https://sifchain.finance/wp-content/uploads/2020/12/Sifchain-Website-Privacy-Policy.pdf"
              >Privacy Policy</a
            >
            <a target="_blank" href="https://sifchain.finance/#roadmap"
              >Roadmap</a
            >
            <a target="_blank" href="https://sifchain.finance/legal-disclamer/"
              >Legal Disclaimer</a
            >
          </div>
        </div>
      </div>
    </div>
    <div
      class="footer"
      :class="{
        active: !active,
      }"
    >
      <div class="items">
        <div class="toggle-button" @click="toggleActive">
          <Icon icon="info-full-circle" />Stay in the Loop
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.container {
  position: fixed;
  width: 100vw;
  height: 51px;
  bottom: 0;
  left: 0;
}
.footer {
  font-size: 12px; /* make $variable */
  line-height: 14px; /* make $variable */
  text-align: left;
  font-weight: 400; /* make $variable */
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 51px;
}

.footer {
  opacity: 0;
  z-index: 100;
}
.footer.active {
  opacity: 1;
  transition: opacity 0.5s ease-out;
  z-index: 200;
}
.items {
  display: flex;
  position: absolute;
  width: 100%;
  align-items: center;
  height: 51px;
}

.left {
  display: flex;
  justify-content: center;
  height: 51px;
  align-items: center;
}
.right {
  flex: 1 1 auto;
  height: 51px;
  font-size: 13px;
  justify-content: flex-end;
  display: flex;
}
.email-input {
  border-bottom-left-radius: 6px;
  border-top-left-radius: 6px;
  padding: 8px;
  outline: none;
  border: none;
  font: inherit; /* TODO - This is used around the app, makes things hard */
  height: 30px;
  margin-left: 10px;
}
.cta {
  display: flex;
  align-items: center;
  padding-left: 15px;
  height: 51px;
  color: #fff;
}
.thankyou {
  margin-left: 10px;
  font-weight: bold;
}
.toggle-button {
  margin-left: 20px;
  padding: 0px 10px 0px 2px;
  height: 24px;
  border: 1px solid #ffffff;
  box-sizing: border-box;
  border-radius: 20px;
  color: #fff;
  cursor: pointer;
  background: rgba(0, 0, 0, 0.4);
  font-weight: bold;
  border-radius: 20px;
  display: flex;
  align-items: center;
  &:hover {
    background: rgba(0, 0, 0, 0.6);
  }
  span {
    margin-right: 6px;
    margin-left: 2px;
    margin-top: 3px;
  }
}
.links {
  display: flex;
  justify-content: flex-end;
  padding-right: 15px;
  height: 51px;
  align-items: center;
  a {
    margin-left: 35px;
    color: #fff;
  }
}
.backdrop {
  background: black;
  position: absolute;
  width: 100%;
  opacity: 0.5;
  height: 51px;
}
</style>
<style lang="scss">
.cta {
  .btn {
    border-bottom-left-radius: 0px;
    border-top-left-radius: 0px;
    display: inline;
    font-size: 13px;
    font: unset;
  }
}
</style>
