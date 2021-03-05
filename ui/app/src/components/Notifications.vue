<script lang="ts">
import { computed, defineComponent } from "vue";
import { reactive, ref, Ref } from "@vue/reactivity"; /* eslint-disable-line */
import { useCore } from "@/hooks/useCore";
import { AppEvent } from "ui-core/src/api/NotificationService";

// Message?
type Notification = {
  id?: string; // id would be used to remove timeout, may only need to be local type
  type: "error" | "success" | "info";
  message: string;
  detail?: {
    type: "etherscan" | "websocket" | "info";
    message: string;
  };
  loader?: boolean;
};

// Visual Notifications are a view level system here we work out which ones are displayed to the user
function parseEventToNotifications(event: AppEvent): Notification | null {
  if (event.type === "NoLiquidityPoolsFoundEvent") {
    return {
      type: "error",
      message: "No Liquidity Pools Found",
      detail: {
        type: "info",
        message: "Create liquidity pool to swap.",
      },
    };
  }

  if (event.type === "TransactionErrorEvent") {
    return {
      type: "error",
      message: event.payload.message,
    };
  }

  if (event.type === "PegTransactionPendingEvent") {
    return {
      type: "info",
      message: "Pegged Transaction Pending",
      detail: {
        type: "etherscan",
        message: event.payload.hash,
      },
      loader: true,
    };
  }

  if (event.type === "PegTransactionErrorEvent") {
    return {
      type: "error",
      message: event.payload.message,
    };
  }

  if (event.type === "PegTransactionCompletedEvent") {
    return {
      type: "success",
      message: `Transfer ${event.payload.hash} has succeded.`,
    };
  }

  if (event.type === "WalletConnectedEvent") {
    const message = {
      sif: "Sif Account Connected",
      eth: "Connected to Metamask",
    }[event.payload.walletType];

    return {
      type: "success",
      message,
      detail: {
        type: "info",
        message: event.payload.address,
      },
    };
  }
  console.error("Have not captured event", JSON.stringify(event));
  return null;
}

export default defineComponent({
  name: "Notifications",
  components: {},
  setup() {
    const { api } = useCore();
    const notifications = reactive<Notification[]>([]);

    api.NotificationService.onAny((event) => {
      const notification = parseEventToNotifications(event);
      if (notification !== null) notifications.unshift(notification);
    });

    return {
      notifications,
      removeItem(index: any) {
        notifications.splice(index, 1);
      },
    };
  },
});
</script>

<template>
  <div class="notifications-container">
    <transition-group name="list">
      <div
        v-for="(item, index) in notifications"
        v-bind:key="item.message"
        class="notification"
        v-bind:class="item.type"
        v-on:click="removeItem(index)"
      >
        <!-- <div class="x">x</div> -->
        <div class="inner">
          <div class="message">
            <div class="circle" v-if="item.type !== 'info'"></div>
            <div>{{ item.message }}</div>
          </div>
          <div class="detail" v-show="item.detail">
            <div v-if="item.detail?.type === 'etherscan'">
              Check on
              <a
                target="_blank"
                :href="`https://etherscan.io/tx/${item.detail?.message}`"
                @click.stop
                >Block Explorer</a
              >
            </div>
            <div v-else-if="item.detail?.type === 'info'">
              {{ item.detail.message }}
            </div>
          </div>
        </div>
      </div>
    </transition-group>
  </div>
</template>

<style lang="scss" scoped>
.notifications-container {
  position: fixed;
  bottom: 0px;
  right: 16px;
  height: auto;
  .list-enter-active,
  .list-leave-active {
    transition: all 0.5s ease;
  }
  .list-enter-from,
  .list-leave-to {
    opacity: 0;
    transform: translateX(200px);
  }
  .notification {
    background: white;
    padding: 3px;
    margin-bottom: 16px;
    text-align: left;
    border-radius: 8px;
    position: relative;
    width: 250px;
    cursor: pointer;
    .x {
      display: none;
    }
    &:hover .x {
      display: block;
      position: absolute;
    }
    .inner {
      border-radius: 6px;
      padding: 4px 8px;
      display: flex;
      flex-direction: column;
      .message {
        display: flex;
        flex-direction: row;
        align-items: center;
      }
      .circle {
        width: 8px;
        height: 8px;
        border-radius: 8px;
        margin-right: 8px;
        flex: none;
      }
      .detail {
        font-size: 12px;
        color: #9f9f9f;
        line-height: 1.1;
      }
    }
    &.error {
      .inner {
        border: 1px solid #b51a1a;
        color: #b51a1a;
        .circle {
          background: #b51a1a;
        }
      }
    }
    &.success {
      .inner {
        border: 1px solid #699829;
        color: #699829;
        .circle {
          background: #699829;
        }
      }
    }
    &.info {
      .inner {
        border: 1px solid #9f9f9f;
        color: #9f9f9f;
      }
    }
  }
}
</style>