<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  components: {},
  async setup() {
    function isNumeric(s: any) {
      return s - 0 == s && ("" + s).trim().length > 0;
    }

    function formatNumberString(x: string) {
      return x.replace(/\B(?=(?=\d*\.)(\d{3})+(?!\d))/g, ",");
    }

    const data = await fetch(
      "https://vtdbgplqd6.execute-api.us-west-2.amazonaws.com/default/tokenstats"
    );
    const json = await data.json();
    const pools = json.body ? json.body.pools : "";
    if (!pools || pools.length < 1) {
      return "";
    }

    let total: number = 0.0;
    pools.map((p: any) => {
      const depth = p.poolDepth;
      if (isNumeric(depth)) {
        total += parseFloat(depth) * 2;
      }
    });

    const tlv = "TVL: $" + formatNumberString(total.toFixed(1));

    return {
      tlv: tlv.substring(0, tlv.length - 2),
    };
  },
});
</script>

<template>
  <div>
    <div v-if="tlv" class="rowan">
      <div>
        {{ tlv }}
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.rowan {
  min-height: 20px;
  background: transparent;
  color: $c_gray_700;
  border-radius: $br_md;
  letter-spacing: 0.1px;
  font-family: Arial, Helvetica, sans-serif;
  font-style: normal;
  font-weight: bold;
  font-size: 12px;
  padding: 1px 5px;
  height: auto;
  line-height: initial;
  border: 1px solid $c_gray_400;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 2px;
  margin-right: 48px;

  .image {
    height: 16px;
    margin-right: 4px;
  }
}
</style>
