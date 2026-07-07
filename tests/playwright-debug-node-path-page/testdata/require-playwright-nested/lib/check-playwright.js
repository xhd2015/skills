const { chromium } = require('playwright');
module.exports = {
  check() {
    console.log(typeof chromium.launch === 'function' ? 'playwright-ok' : 'playwright-fail');
  },
};