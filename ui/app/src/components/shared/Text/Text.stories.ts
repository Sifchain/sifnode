import Text from "./Text.vue";

// TODO - This is just a place holder to practise with the new format/amount lib
//        Though we should investigate the idea of having a typography system
//        Will push to work with the new designer to build a design language

function numberWithCommas(n) {
  var parts = n.toString().split(".");
  return (
    parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ",") +
    (parts[1] ? "." + parts[1] : "")
  );
}

export default {
  title: "Text",
  component: Text,
  argTypes: {
    amount: {
      control: "text",
      description: "Amount",
      defaultValue: "000",
    },
    decimals: {
      control: "text",
      description: "Decimals",
      defaultValue: "2",
    },
    commas: {
      control: "boolean",
      description: "Decimals",
      defaultValue: false,
    },
    padding: {
      control: "text",
      description: "Padding of 0's onto decimal",
      defaultValue: "0",
    },
  },
};

const formatAmount = (amount, { decimals, commas, padding }) => {
  // Place holder function for the final prototype
  let formattedAmount = `${amount}`;
  if (decimals > 0) {
    formattedAmount = Number(formattedAmount).toFixed(decimals);
  }
  if (commas) {
    formattedAmount = numberWithCommas(formattedAmount);
  }
  return formattedAmount;
};

const Template = (args, { argTypes }) => ({
  props: [""],
  components: { Text },
  setup: () => {
    const { amount, decimals, commas, padding } = args;
    const formattedAmount = formatAmount(amount, { decimals, commas, padding });
    return { formattedAmount };
  },
  template: "<text>{{ formattedAmount }}</text>",
});

export const PrimaryText = Template.bind({});
PrimaryText.args = {
  amount: "1000002.000123",
};
