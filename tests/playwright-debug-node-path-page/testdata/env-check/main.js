const nodePath = process.env.NODE_PATH || '';
console.log(nodePath.includes('playwright-debug') ? 'node-path-set' : 'node-path-missing');