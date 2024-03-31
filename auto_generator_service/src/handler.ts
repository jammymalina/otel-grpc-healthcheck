import GithubService from "./services/githubService";

const githubService = new GithubService();

const handler = async (): Promise<string> => {
  await githubService.getRepositoryTags("open-telemetry", "opentelemetry-collector-contrib");
  return "OK";
};

export default handler;
