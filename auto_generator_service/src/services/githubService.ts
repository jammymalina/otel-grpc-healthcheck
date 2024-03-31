import { Octokit } from "@octokit/rest";

import logger from "../utils/logger";
import RepositoryVersionTag from "../domain/repositoryVersionTag";

export default class GithubService {
  private readonly octokit: Octokit;

  constructor() {
    this.octokit = new Octokit();
  }

  async getRepositoryTags(owner: string, repo: string): Promise<RepositoryVersionTag[]> {
    const response = await this.octokit.rest.repos.listTags({ owner, repo, per_page: 100 });
    logger.info("Received response from list tags octo API", response);
    return [];
  }
}
