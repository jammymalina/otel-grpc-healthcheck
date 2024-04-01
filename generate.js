#!/usr/bin/env node

import * as fs from "node:fs/promises";
import { compile as ejsCompile } from "ejs";
import yargs from "yargs/yargs";
import { hideBin } from "yargs/helpers";

const argv = yargs(hideBin(process.argv))
  .options({
    otelVersion: {
      alias: "o",
      describe: "provide an otel version",
      demandOption: true,
      type: "string",
    },
    goVersion: {
      alias: "g",
      describe: "provide a go version",
      demandOption: true,
      type: "string",
    },
  })
  .help()
  .parse();

(async (argv) => {
  const goModTemplateData = (await fs.readFile("./go.mod.ejs")).toString("utf8");
  const goModTemplate = ejsCompile(goModTemplateData);

  const goModData = goModTemplate({
    goVersion: argv.goVersion,
    otelVersion: argv.otelVersion,
  });

  await fs.writeFile("./go.mod", goModData);
})(argv);
