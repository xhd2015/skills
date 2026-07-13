'use strict';

const fs = require('fs');
const os = require('os');
const path = require('path');
const { createRequire } = require('module');
const { chromium } = require('playwright');

function envFlag(name, fallback) {
  const v = process.env[name];
  if (v == null || v === '') return fallback;
  return v;
}

function parseExtensionPaths() {
  const raw = envFlag('PLAYWRIGHT_DEBUG_EXTENSION_PATHS', '');
  if (!raw) return [];
  const sep = path.delimiter; // ':' on unix, ';' on windows
  return raw
    .split(sep)
    .map((p) => p.trim())
    .filter(Boolean);
}

/**
 * Resolve profile dir for launchPersistentContext.
 * - If PLAYWRIGHT_DEBUG_USER_DATA_DIR is set: use it (mkdir; never delete).
 * - Else if allowTemp: create a temp dir and delete on close.
 * - Else: return empty (caller uses ephemeral chromium.launch).
 */
function resolveUserDataDir(allowTemp) {
  let userDataDir = envFlag('PLAYWRIGHT_DEBUG_USER_DATA_DIR', '');
  let createdTemp = false;
  if (userDataDir) {
    fs.mkdirSync(userDataDir, { recursive: true });
    return { userDataDir, createdTemp: false };
  }
  if (allowTemp) {
    userDataDir = fs.mkdtempSync(path.join(os.tmpdir(), 'pw-debug-profile-'));
    createdTemp = true;
  }
  return { userDataDir, createdTemp };
}

/**
 * Launch Chromium either in default headless mode, with a persistent profile
 * (--user-data-dir), and/or with unpacked extensions via launchPersistentContext.
 *
 * Env:
 *   PLAYWRIGHT_DEBUG_LAUNCH_MODE = default | extension
 *   PLAYWRIGHT_DEBUG_EXTENSION_PATHS = path-list (path.delimiter joined)
 *   PLAYWRIGHT_DEBUG_USER_DATA_DIR = optional profile dir (persists login/cookies)
 *   PLAYWRIGHT_DEBUG_HEADED = 1 | 0
 */
async function launchBrowser() {
  const mode = envFlag('PLAYWRIGHT_DEBUG_LAUNCH_MODE', 'default');
  const headed = envFlag('PLAYWRIGHT_DEBUG_HEADED', mode === 'extension' ? '1' : '0') === '1';
  const extPaths = parseExtensionPaths();
  const wantExtension = mode === 'extension' || extPaths.length > 0;
  const explicitProfile = envFlag('PLAYWRIGHT_DEBUG_USER_DATA_DIR', '') !== '';

  if (wantExtension) {
    if (extPaths.length === 0) {
      throw new Error(
        'extension launch mode requires PLAYWRIGHT_DEBUG_EXTENSION_PATHS (unpacked extension dir with manifest.json)',
      );
    }
    for (const p of extPaths) {
      if (!fs.existsSync(path.join(p, 'manifest.json'))) {
        throw new Error(`extension path missing manifest.json: ${p}`);
      }
    }

    // Temp profile only when user did not pass --user-data-dir.
    const { userDataDir, createdTemp } = resolveUserDataDir(true);
    const joined = extPaths.join(',');
    const context = await chromium.launchPersistentContext(userDataDir, {
      // Extensions are unreliable in classic headless; default headed for extension mode.
      headless: !headed,
      args: [
        `--disable-extensions-except=${joined}`,
        `--load-extension=${joined}`,
        '--no-first-run',
        '--no-default-browser-check',
      ],
    });

    const page = context.pages()[0] || (await context.newPage());
    const browser = typeof context.browser === 'function' ? context.browser() : null;

    return {
      mode: 'extension',
      browser,
      context,
      page,
      userDataDir,
      createdTemp,
      extensionPaths: extPaths,
      async close() {
        await context.close();
        if (createdTemp) {
          try {
            fs.rmSync(userDataDir, { recursive: true, force: true });
          } catch (_) {
            /* ignore cleanup errors */
          }
        }
      },
    };
  }

  // Persistent profile without extensions: keep cookies/localStorage/login
  // across reboots when the dir lives outside /tmp (user-chosen path).
  if (explicitProfile) {
    const { userDataDir } = resolveUserDataDir(false);
    const context = await chromium.launchPersistentContext(userDataDir, {
      headless: !headed,
      args: ['--no-first-run', '--no-default-browser-check'],
    });
    const page = context.pages()[0] || (await context.newPage());
    const browser = typeof context.browser === 'function' ? context.browser() : null;
    return {
      mode: 'persistent',
      browser,
      context,
      page,
      userDataDir,
      createdTemp: false,
      extensionPaths: [],
      async close() {
        // Never delete a user-specified profile.
        await context.close();
      },
    };
  }

  const browser = await chromium.launch({ headless: !headed });
  const page = await browser.newPage();
  return {
    mode: 'default',
    browser,
    context: null,
    page,
    userDataDir: '',
    createdTemp: false,
    extensionPaths: [],
    async close() {
      await browser.close();
    },
  };
}

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

  const launched = await launchBrowser();

  const AsyncFunction = Object.getPrototypeOf(async function () {}).constructor;
  // Extra bindings after the original six: context, extensionPaths
  const fn = new AsyncFunction(
    'browser',
    'page',
    'chromium',
    'require',
    '__filename',
    '__dirname',
    'context',
    'extensionPaths',
    source,
  );

  try {
    await fn(
      launched.browser,
      launched.page,
      chromium,
      userRequire,
      userScript,
      scriptDir,
      launched.context,
      launched.extensionPaths,
    );
  } finally {
    await launched.close();
  }
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
