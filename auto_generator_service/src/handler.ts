import GithubService from "./services/githubService";
import logger from "./utils/logger";

const githubService = new GithubService();

enum LambdaResponse {
  OK = "OK",
  NO_TAGS = "NO_TAGS",
}

const handler = async (): Promise<LambdaResponse> => {
  const otelTags = await githubService.getRepositoryTags("open-telemetry", "opentelemetry-collector-contrib");
  const healthTags = await githubService.getRepositoryTags("jammymalina", "otel-grpc-healthcheck");

  logger.info("Received repo tags", { otelTags, healthTags });

  if (otelTags.length === 0 || healthTags.length === 0) {
    logger.info("No repo tags for at least one of the repos, exiting");
    return LambdaResponse.NO_TAGS;
  }

  const [latestOtelTag] = otelTags.slice(-1);
  const [latestHealthTag] = healthTags.slice(-1);
  logger.info("Latest tags", { otelTag: latestOtelTag.version, healthTag: latestHealthTag.version });

  return LambdaResponse.OK;
};

export default handler;
