export default class SemanticVersion {
  constructor(
    public readonly major: number,
    public readonly minor: number,
    public readonly patch: number
  ) {}

  static initFromVersionString(versionString: string): SemanticVersion | null {
    const regex = /v?(\d+)\.(\d+)\.(\d+)/;
    const match = regex.exec(versionString);

    if (!match) {
      return null;
    }

    return new SemanticVersion(parseInt(match[1], 10), parseInt(match[2], 10), parseInt(match[3], 10));
  }

  compare(otherVersion: SemanticVersion): number {
    if (this.major !== otherVersion.major) {
      return this.major - otherVersion.major;
    }
    if (this.minor !== otherVersion.minor) {
      return this.minor - otherVersion.minor;
    }
    return this.patch - otherVersion.patch;
  }
}
