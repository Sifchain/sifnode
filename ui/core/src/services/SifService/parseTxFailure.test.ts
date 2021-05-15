import { parseTxFailure } from "./parseTxFailure";

// possibly slightly redundant test for coverage
test("parseTxFailure", () => {
  expect(
    parseTxFailure({
      transactionHash: "123",
      rawLog: "",
    }),
  ).toMatchObject({
    hash: "123",
    memo: "There was an unknown failure",
    state: "failed",
  });

  expect(
    parseTxFailure({
      transactionHash: "123",
      rawLog: "something was below expected",
    }),
  ).toMatchObject({
    hash: "123",
    memo: "Your transaction has failed - Received amount is below expected",
    state: "failed",
  });

  expect(
    parseTxFailure({
      transactionHash: "123",
      rawLog: "yegads swap_failed!",
    }),
  ).toMatchObject({
    hash: "123",
    memo: "Your transaction has failed",
    state: "failed",
  });

  expect(
    parseTxFailure({
      transactionHash: "123",
      rawLog: "your Request rejected!",
    }),
  ).toMatchObject({
    hash: "123",
    memo: "You have rejected the transaction",
    state: "rejected",
  });
});
