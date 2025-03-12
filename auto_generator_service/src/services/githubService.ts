import { Octokit } from "@octokit/rest";

import logger from "../utils/logger";
import RepositoryVersionTag from "../domain/repositoryVersionTag";
import SemanticVersion from "../domain/semanticVersion";

interface Tag {
  name: string;
    commit: {
        sha: string;
        url: string;
    };
    zipball_url: string;
    tarball_url: string;
    node_id: string;
}

export default class GithubService {
  private readonly octokit: Octokit;

  constructor() {
    this.octokit = new Octokit();
  }

  async getRepositoryTags(owner: string, repo: string, prefix?: string): Promise<RepositoryVersionTag[]> {
    const listedTags = []
    let pageIndex = 1;

    const prefixFilter = typeof prefix === "undefined" ? () => true : (tag: Tag) => tag.name.startsWith(prefix);
    const extractVersion = typeof prefix === "undefined" ? (tagName: string) => tagName : (tagName: string) => tagName.substring(prefix.length, tagName.length);

    while (listedTags.length < 5 && pageIndex < 100) {
      const response = await this.octokit.rest.repos.listTags({ owner, repo, per_page: 100, page: pageIndex });
      logger.info("Received response from list tags octo API", { details: { response } });

      const allowedStatusCodes = [200, 204];

      if (!allowedStatusCodes.includes(response.status)) {
        throw new Error(`Unable to get repository tags for ${owner}/${repo}`);
      }

      const items = Array.isArray(response.data) ? response.data : [];
      listedTags.push(...items.filter(prefixFilter))

      pageIndex += 1;
    }

    const tags = listedTags
      .map((tagItem) => {
        const version = SemanticVersion.initFromVersionString(extractVersion(tagItem.name));
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
