<script lang="ts">
import { computed } from "@vue/reactivity";
import { defineComponent } from "vue";
import { formatSymbol, useAssetItem } from "@/components/shared/utils";

export default defineComponent({
  props: {
    pool: {
      type: Object,
    },
    liqAPY: {
      type: Object,
    },
    inline: Boolean,
  },

  components: {},

  setup(props) {
    function formatNumberString(x: string) {
      return x.replace(/\B(?=(?=\d*\.)(\d{3})+(?!\d))/g, ",");
    }

    const symbol = computed(() => props.pool?.symbol ?? "");
    const asset = useAssetItem(symbol);
    const token = asset.token;
    const image = computed(() => {
      if (!token.value) {
        return "";
      } else {
        return token.value.imageUrl;
      }
    });

    const priceToken = formatNumberString(
      parseFloat(props.pool?.priceToken).toFixed(2)
    );
    const poolDepth = formatNumberString(
      parseFloat(props.pool?.poolDepth).toFixed(2)
    );
    const volume = formatNumberString(
      parseFloat(props.pool?.volume).toFixed(1)
    );
    const poolAPY = formatNumberString(
      (
        parseFloat(props.pool?.volume) / parseFloat(props.pool?.poolDepth)
      ).toFixed(1)
    );

    return {
      symbol,
      image,
      priceToken,
      poolDepth,
      volume,
      poolAPY,
      formatNumberString,
    };
  },
});
</script>

<template>
  <div class="pool-list-item">
    <div class="pool-asset" :class="{ inline: inline }">
      <div class="col-sm-s">
        <img v-if="image" width="22" height="22" :src="image" class="image" />
        <div class="placeholder" v-else></div>
        <div class="icon">
          <span>{{ "c" + symbol.toString().slice(1).toUpperCase() }}</span>
        </div>
      </div>
      <div class="col-sm">
        <span>${{ priceToken }}</span>
      </div>
      <div class="col-sm">
        <span>${{ poolDepth }}</span>
      </div>
      <div class="col-sm">
        <span>${{ volume }}</span>
      </div>
      <div class="col-sm">
        <span>{{ poolAPY }}%</span>
      </div>
      <div class="col-lg">
        <span>{{ formatNumberString(parseFloat(liqAPY).toFixed(1)) }}%</span>
      </div>
      <div class="col-sm">
        <span
          >{{
            formatNumberString(
              (parseFloat(poolAPY) + parseFloat(liqAPY)).toFixed(1)
            )
          }}%</span
        >
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.pool-list-item {
  padding: 12px 1em;

  &:not(:last-of-type) {
    border-bottom: $divider;
  }
}

.pool-asset {
  display: flex;
  justify-content: space-evenly;
  align-items: center;

  .image {
    height: 20px;
    margin-right: 8px;

    & > * {
      border-radius: 16px;

      &:nth-child(2) {
        position: relative;
        left: -6px;
      }
    }
  }

  &.inline {
    display: inline-flex;

    & > span {
      margin-right: 0;
    }
  }

  .placeholder {
    display: inline-block;
    background: #aaa;
    box-sizing: border-box;
    border-radius: 16px;
    height: 20px;
    width: 20px;
    text-align: center;
  }

  .col-sm-s {
    padding-left: 12px;
    min-width: 102px;
    width: 10%;
    display: flex;
    justify-content: start;
  }

  .col-sm {
    min-width: 102px;
    display: flex;
    justify-content: center;
    font-size: $fs_md;
    color: $c_text;
  }

  .col-md {
    min-width: 110px;
    font-size: $fs_md;
    color: $c_text;
  }

  .col-lg {
    min-width: 168px;
    font-size: $fs_md;
    color: $c_text;
  }
}
</style>
