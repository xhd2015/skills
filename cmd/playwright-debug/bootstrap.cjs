'use strict';

const fs = require('fs');
const path = require('path');
const { createRequire } = require('module');
const { chromium } = require('playwright');

async function main() {
  const userScript = path.resolve(process.argv[2]);
  if (!fs.existsSync(userScript)) {
    console.error(`Script file not found: ${userScript}`);
    process.exit(1);
  }

  const scriptDir = path.dirname(userScript);
  process.chdir(scriptDir);

  const source = fs.readFileSync(userScript, 'utf8');
  const userRequire = createRequire(userScript);

  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();

  const AsyncFunction = Object.getPrototypeOf(async function () {}).constructor;
  const fn = new AsyncFunction(
    'browser',
    'page',
    'chromium',
    'require',
    '__filename',
    '__dirname',
    source,
  );

  try {
    await fn(browser, page, chromium, userRequire, userScript, scriptDir);
  } finally {
    await browser.close();
  }
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});