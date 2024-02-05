#!/usr/bin/env node

const fs = require("fs/promises");
const ejs = require("ejs");
const yargs = require("yargs/yargs");
const { hideBin } = require("yargs/helpers");

const argv = yargs(hideBin(process.argv))
  .options({
    otelVersion: {
      alias: "o",
      describe: "provide an otel version",
      demandOption: true,
    },
    goVersion: {
      alias: "g",
      describe: "provide a go version",
      demandOption: true,
    },
  })
  .help()
  .parse();

(async (argv) => {
  const goModTemplateData = (await fs.readFile("./go.mod.ejs")).toString("utf8");
  const goModTemplate = ejs.compile(goModTemplateData);

  const goModData = await goModTemplate({
    goVersion: argv.goVersion,
    otelVersion: argv.otelVersion,
  });

  await fs.writeFile("./go.mod", goModData);
})(argv);
