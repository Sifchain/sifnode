<script>
import Box from "@/components/shared/Box.vue";
import { Copy, SubHeading } from "@/components/shared/Text";
import Loader from "@/components/shared/Loader.vue";
import Tooltip from "@/components/shared/Tooltip.vue";
import Icon from "@/components/shared/Icon.vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import SifButton from "@/components/shared/SifButton.vue";
import { format } from "ui-core/src/utils/format";

const REWARD_INFO = {
  lm: {
    label: "Liquidity Mining",
    description:
      "Earn additional rewards by providing liquidity to any of Sifchain's pools.",
  },
  vs: {
    label: "Validator Subsidy",
    description:
      "Earn additional rewards by staking a node or delegating to a staked node.",
  },
};

export default {
  props: {
    type: {
      type: String,
    },
    data: {
      type: Object,
    },
    address: {
      type: String,
    },
  },
  components: {
    SifButton,
    AssetItem,
    Copy,
    SubHeading,
    Box,
    Tooltip,
    Icon,
    Loader,
  },
  emits: ["openModal"],
  methods: {
    format,
    openClaimModal() {
      this.modalOpen = true;
    },
    requestClose() {
      this.modalOpen = false;
    },
    claimRewards() {
      alert("claim logic/keplr goes here");
    },
  },
  data() {
    return {
      modalOpen: false,
      loadingLm: true,
      loadingVs: true,
      REWARD_INFO,
    };
  },
};
</script>

<template>
  <Box>
    <div class="reward-container">
      <SubHeading>{{ REWARD_INFO[type].label }}</SubHeading>
      <Copy>
        {{ REWARD_INFO[type].description }}
      </Copy>
      <div class="details-container">
        <Loader v-if="!data" black />

        <div v-else class="amount-container">
          <div class="reward-rows">
            <div class="reward-row">
              <div class="row-label">Claimable Amount</div>
              <div class="row-amount">
                {{
                  format(data.totalClaimableCommissionsAndClaimableRewards, {
                    mantissa: 4,
                  }) || "0"
                }}
              </div>
              <AssetItem symbol="Rowan" :label="false" />
            </div>
            <div v-if="type === 'vs'" class="reward-row">
              <div class="row-label">Rewards Based on Commission</div>
              <div class="row-amount">
                {{
                  format(
                    data.currentTotalCommissionsOnClaimableDelegatorRewards,
                    {
                      mantissa: 4,
                    },
                  ) || "0"
                }}
              </div>
              <AssetItem symbol="Rowan" :label="false" />
            </div>

            <div class="reward-row">
              <div class="row-label">
                Dispensed Rewards
                <Tooltip>
                  <template #message>
                    <div class="tooltip">
                      Rewards that have already been dispensed.
                    </div>
                  </template>
                  <Icon icon="info-box-black" />
                </Tooltip>
              </div>
              <div class="row-amount">
                {{ format(data.dispensed, { mantissa: 4 }) || "0" }}
              </div>
              <AssetItem symbol="Rowan" :label="false" />
            </div>

            <div class="reward-row secondary">
              <div class="row-label">
                Projected Full Amount
                <Tooltip>
                  <template #message>
                    <div class="tooltip">
                      <div v-if="data.maturityDate">
                        Projected Full Maturity Date: <br />
                        <span class="tooltip-date">{{
                          data.maturityDate
                        }}</span>
                        <span v-if="data.nextRewardProjectedAPYOnTickets">
                          Projected Fully Maturated APY: <br />
                          <span class="tooltip-date">
                            {{
                              format(
                                data.nextRewardProjectedAPYOnTickets * 100,
                                {
                                  mantissa: 2,
                                },
                              )
                            }}%</span
                          >
                        </span>
                        <br /><br />
                      </div>
                      This is your estimated projected full reward amount that
                      you can earn if you were to leave your current liquidity
                      positions in place to the above mentioned date. This
                      includes projected future rewards, and already
                      claimed/disbursed previous rewards. This number can
                      fluctuate due to other market conditions and this number
                      is a representation of the current market as it is in this
                      very moment.
                    </div>
                  </template>
                  <Icon icon="info-box-black" />
                </Tooltip>
              </div>
              <div class="row-amount">
                {{
                  format(data.totalCommissionsAndRewardsAtMaturity, {
                    mantissa: 4,
                  }) || "0"
                }}
              </div>
              <AssetItem symbol="Rowan" :label="false" />
            </div>
          </div>
          <div class="reward-buttons">
            <a
              class="more-info-button mr-8"
              target="_blank"
              :href="`https://cryptoeconomics.sifchain.finance/#${address}&type=${type}`"
              >More Info</a
            >

            <!-- :disabled="(data.claimableReward - data.claimed) === 0" -->
            <SifButton
              @click="$emit('openModal', type)"
              :primary="true"
              :disabled="true"
              >Claim</SifButton
            >
          </div>
        </div>
      </div>
    </div>
  </Box>
</template>

<style lang="scss" scoped>
.rewards-container {
  display: flex;
  flex-direction: column;
  > :first-child {
    margin-top: $margin_medium;
  }
  width: 100%;
  > :nth-child(1) {
    margin-bottom: $margin_medium;
  }
  .reward-container {
    flex-direction: column;
    > :nth-child(1),
    > :nth-child(2) {
      margin-bottom: $margin_small;
    }
  }
  .reward-rows {
    display: flex;
    flex-direction: column;
    margin-bottom: 15px;
    color: #343434;
  }
  .reward-row {
    display: flex;
    width: 100%;
    justify-content: space-between;
    font-size: $fs;
    font-weight: 400;
    &.secondary {
      color: #818181;
    }
    .row-label {
      flex: 1 1 auto;
      text-align: left;
    }
    .row-amount {
      width: 100px;
      text-align: right;
    }
    .row {
      width: 15px;
      margin-left: 2px;
    }
  }
}

.reward-buttons {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  .more-info-button {
    background: #f3f3f3;
    color: #343434;
    font-weight: 100;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  .more-info-button,
  .btn {
    width: 300px;
    border-radius: 6px;
    display: flex;
    font-size: $fs;
    height: 30px;
  }
  .reward-button {
    text-align: center;
  }
}

.tooltip-date {
  font-weight: 600;
}
</style>
