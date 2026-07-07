async function greet(page) {
  await page.goto('about:blank');
  console.log('explicit-page-ok');
}
module.exports = { greet };