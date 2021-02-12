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

      const symbol = computed(() =>
              props.pool?.symbol ? formatSymbol(props.pool?.symbol) : ""
      );
      const asset = useAssetItem(symbol);
      const token = asset.token;
      const image = computed(() => {
        if (!token.value) {
          return "";
        } else {
          return token.value.imageUrl;
        }
      });

      const poolUSD = formatNumberString(
              parseFloat(props.pool?.poolUSD).toFixed(2)
      );
      const poolRowan = formatNumberString(
              parseFloat(props.pool?.poolRowan).toFixed(2)
      );
      const externalUSD = formatNumberString(
              parseFloat(props.pool?.externalUSD).toFixed(8)
      );
      const externalRowan = formatNumberString(
              parseFloat(props.pool?.externalRowan).toFixed(8)
      );

      return {
        symbol,
        poolUSD,
        poolRowan,
        externalUSD,
        externalRowan,
        image,
      };
    },
  });
</script>

<template>
  <div class="pool-list-item">
    <div class="pool-asset" :class="{ inline: inline }">
      <div class="col-sm">
        <img v-if="image" width="22" height="22" :src="image" class="image" />
        <div class="placeholder" v-else></div>
        <div class="icon">
          <span>{{ symbol.replace("c", "") }}</span>
        </div>
      </div>
      <div class="col-md">
        <span>{{ externalUSD }}</span>
      </div>
      <div class="col-md">
        <span>{{ poolUSD }}</span>
      </div>
      <div class="col-lg">
        <span>{{ liqAPY }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .pool-list-item {
    padding: 12px 12px;

    &:not(:last-of-type) {
      border-bottom: $divider;
    }
  }

  .pool-asset {
    display: flex;
    justify-content: space-evenly;
    align-items: center;

    .image {
      height: 22px;
      margin-right: 10px;

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
      height: 22px;
      width: 22px;
      text-align: center;
    }

    .col-lg {
      min-width: 200px;
      width: 27%;
      font-size: $fs_md;
      color: $c_text;
    }

    .col-md {
      min-width: 160px;
      width: 27%;
      font-size: $fs_md;
      color: $c_text;
    }

    .col-sm {
      min-width: 100px;
      width: 19%;
      display: flex;
      justify-content: start;
    }
  }
</style>
