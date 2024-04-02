import Environment from "./domain/environment";
import GithubService from "./services/githubService";
import PluginUpdateService from "./services/pluginUpdateService";
import logger from "./utils/logger";

const environment = new Environment();
const githubService = new GithubService();
const pluginUpdateService = new PluginUpdateService(environment.projectName);

enum LambdaResponse {
  OK = "OK",
  NO_UPDATE_REQUIRED = "NO_UPDATE_REQUIRED",
  NO_TAGS = "NO_TAGS",
}

const handler = async (): Promise<LambdaResponse> => {
  const [otelTags, healthTags] = await Promise.all([
    githubService.getRepositoryTags("open-telemetry", "opentelemetry-collector-contrib"),
    githubService.getRepositoryTags("jammymalina", "otel-grpc-healthcheck"),
  ]);

  if (otelTags.length === 0 || healthTags.length === 0) {
    logger.info("No repo tags for at least one of the repos, exiting");
    return LambdaResponse.NO_TAGS;
  }

  const [latestOtelTag] = otelTags.slice(-1);
  const [latestHealthTag] = healthTags.slice(-1);
  logger.info("Latest tags", { details: { otelTag: latestOtelTag.version, healthTag: latestHealthTag.version } });

  if (latestOtelTag.version.equals(latestHealthTag.version)) {
    logger.info("The latest version is already published");
    return LambdaResponse.NO_UPDATE_REQUIRED;
  }

  logger.info(`Updating the version to ${latestOtelTag.version}`);
  await pluginUpdateService.update(latestOtelTag.version);

  return LambdaResponse.OK;
};

export default handler;
