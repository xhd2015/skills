const args = process.argv.slice(3);
if (args.includes('--help')) {
  console.log('SCRIPT_HELP_OK');
  process.exit(0);
}