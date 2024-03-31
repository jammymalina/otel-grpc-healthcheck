import SemanticVersion from "./semanticVersion";

export default interface RepositoryVersionTag {
  version: SemanticVersion;
  date: Date;
}
