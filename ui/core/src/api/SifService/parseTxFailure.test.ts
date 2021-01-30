import { parseTxFailure } from "./parseTxFailure";

// possibly slightly redundant test for coverage
test("parseTxFailure", () => {
  expect(
    parseTxFailure({
      transactionHash: "123",
      code: 123,
      height: 123,
      rawLog: "",
    })
  ).toEqual({
    hash: "123",
    memo: "Unknown failure",
    state: "failed",
  });

  expect(
    parseTxFailure({
      transactionHash: "123",
      code: 123,
      height: 123,
      rawLog: "something was below expected",
    })
  ).toEqual({
    hash: "123",
    memo: "Swap failed - Received amount is below expected",
    state: "failed",
  });

  expect(
    parseTxFailure({
      transactionHash: "123",
      code: 123,
      height: 123,
      rawLog: "yegads swap_failed!",
    })
  ).toEqual({
    hash: "123",
    memo: "Swap failed",
    state: "failed",
  });

  expect(
    parseTxFailure({
      transactionHash: "123",
      code: 123,
      height: 123,
      rawLog: "your Request rejected!",
    })
  ).toEqual({
    hash: "123",
    memo: "Request Rejected",
    state: "rejected",
  });
});
