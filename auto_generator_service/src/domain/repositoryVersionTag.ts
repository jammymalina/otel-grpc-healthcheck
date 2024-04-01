import SemanticVersion from "./semanticVersion";

export default interface RepositoryVersionTag {
  version: SemanticVersion;
  commitSha: string;
}
