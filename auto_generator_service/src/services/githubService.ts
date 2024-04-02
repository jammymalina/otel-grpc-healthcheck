import { Octokit } from "@octokit/rest";

import logger from "../utils/logger";
import RepositoryVersionTag from "../domain/repositoryVersionTag";
import SemanticVersion from "../domain/semanticVersion";

export default class GithubService {
  private readonly octokit: Octokit;

  constructor() {
    this.octokit = new Octokit();
  }

  async getRepositoryTags(owner: string, repo: string): Promise<RepositoryVersionTag[]> {
    const response = await this.octokit.rest.repos.listTags({ owner, repo, per_page: 10 });
    logger.info("Received response from list tags octo API", { details: { response } });

    const allowedStatusCodes = [200, 204];

    if (!allowedStatusCodes.includes(response.status)) {
      throw new Error(`Unable to get repository tags for ${owner}/${repo}`);
    }

    const tags = (response.data || [])
      .map((tagItem) => {
        const version = SemanticVersion.initFromVersionString(tagItem.name);
        return {
          version,
          commitSha: tagItem.commit.sha,
        };
      })
      .filter((repositoryTag) => !repositoryTag.version.isZero());

    tags.sort((a: RepositoryVersionTag, b: RepositoryVersionTag) => a.version.compare(b.version));
    return tags;
  }
}
