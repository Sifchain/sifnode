describe("Page loads", () => {
  it("Get rowan", () => {
    cy.visit("/");
    cy.contains("Get Rowan");
  });
});
