<script lang="ts">
import { defineComponent } from "vue";
import jsonp from "jsonp";
import queryString from "query-string";
import Panel from "@/components/shared/Panel.vue";
import PanelNav from "@/components/shared/PanelNav/PanelNav.vue";
import Icon from "@/components/shared/Icon.vue";
import SifButton from "@/components/shared/SifButton.vue";

const MAILCHIMP_URL = "https://finance.us2.list-manage.com/subscribe/post-json";

export default defineComponent({
  components: { Panel, PanelNav, Icon, SifButton },
  props: {},
  data() {
    return {
      active: true,
      email: "thomas@gmail.com",
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
      console.log("asd", this.email);
      // const data = new URLSearchParams();
      // for (const pair of new FormData(formElement)) {
      // data.append("email", "thomasalwyndavis@gmail.com");
      // }
      //
      // const mailchimpRes = await fetch(MAILCHIMP_URL, {
      //   method: "post",
      //   body: data,
      // });
      const query = queryString.stringify({
        u: "400787e0a5e23ec37b7b51f74",
        id: "c1ee83387b",
        EMAIL: "thomasalwyndavis@gmail.com",
      });
      const url = `${MAILCHIMP_URL}?${query}`;
      jsonp(url, { param: "c" }, (error, data) => {
        if (error) {
          // say, try again
        } else {
          // say, thank you
        }
      });
    },
  },
});
</script>

<template>
  <div class="container">
    <div v-if="active" class="footer">
      <div class="backdrop"></div>
      <div class="items">
        <div class="left">
          <div class="toggle-button" @click="toggleActive">
            <Icon icon="info-box-white" />Close
          </div>
        </div>
        <div class="right">
          <div class="cta">
            <form @submit.prevent="submitEmail">
              Sign up for Sifchain updates
              <input
                v-model="email"
                type="email"
                name="email"
                class="email-input"
              />
              <SifButton primary type="submit">Stay Informed</SifButton>
            </form>
          </div>
          <div class="links">
            <a href="">Privacy Policy</a>
            <a href="">Roadmap</a>
            <a href="">Legal Disclaimer</a>
          </div>
        </div>
      </div>
    </div>
    <div v-else class="footer">
      <div class="backdrop"></div>
      <div class="items">
        <div class="toggle-button" @click="toggleActive">
          <Icon icon="info-box-white" />Stay in the Loop
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
  position: relative;
  height: 51px;
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
  height: 30px;
}
.cta {
  display: flex;
  align-items: center;
  padding-left: 15px;
  height: 51px;
  color: #fff;
}
.toggle-button {
  margin-left: 20px;
  padding: 0px 8px 0px 2px;
  height: 22px;
  border: 1px solid #ffffff;
  box-sizing: border-box;
  border-radius: 20px;
  color: #fff;
  cursor: pointer;
  background: rgba(0, 0, 0, 0.4);
  border-radius: 20px;
}
.links {
  display: flex;
  justify-content: flex-end;
  padding-right: 15px;
  height: 51px;
  align-items: center;
  a {
    margin-left: 15px;
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
  }
}
</style>
