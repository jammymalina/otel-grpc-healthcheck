#!/usr/bin/env node

const fs = require("fs/promises")
const ejs = require("ejs")
const yargs = require("yargs/yargs");
const { hideBin } = require("yargs/helpers");

const argv = yargs(hideBin(process.argv))
  .options({
    version: {
      alias: "v",
      describe: "provide an otel version",
      demandOption: true,
    },
    goVersion: {
      alias: "g",
      describe: "provide a go version",
      demandOption: true,
    }
  })
  .help()
  .parse().argv;

const goModTemplateData = await fs.readFile("./go.mod.ejs")
const goModTemplate = ejs.compile(goModTemplateData)

const goModData = await goModTemplate({
  goVersion: argv.goVersion,
  version: argv.version,
})

await fs.writeFile("./go.mod", goModData)
