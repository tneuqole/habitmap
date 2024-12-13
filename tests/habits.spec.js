import { test, expect } from "@playwright/test";

test("test /habits", async ({ page }) => {
  await page.goto("http://localhost:4000/habits");
  await expect(page).toHaveTitle(/Habits/);
  await expect(page.getByRole("navigation")).toContainText("Habits");
  await expect(page.getByRole("link", { name: "+" })).toBeVisible();
  await expect(page.getByRole("link", { name: "hello, world!" })).toBeVisible();
  await expect(
    page
      .locator("div")
      .filter({ hasText: /^hello, world!TODO$/ })
      .locator("div"),
  ).toBeVisible();
  await page.getByRole("link", { name: "+" }).click();
  await page.getByRole("main").click();
  await expect(page.getByRole("main")).toMatchAriaSnapshot(`
    - main:
      - heading "Create Habit" [level=1]
      - text: Name
      - textbox "Name"
      - button "Submit"
    `);
  await page.goto("http://localhost:4000/habits");
  await page.getByRole("link", { name: "hello, world!" }).click();
  await expect(page.getByRole("main")).toMatchAriaSnapshot(`
    - main:
      - heading "hello, world!" [level=1]:
        - link "hello, world!"
      - text: TODO
      - button [expanded]
    `);
});
