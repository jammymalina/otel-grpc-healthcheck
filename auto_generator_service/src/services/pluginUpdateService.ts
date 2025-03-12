import { CodeBuildClient, StartBuildCommand } from "@aws-sdk/client-codebuild";

import SemanticVersion from "../domain/semanticVersion";
import logger from "../utils/logger";

export default class PluginUpdateService {
  private readonly client: CodeBuildClient;

  constructor(private readonly projectName: string) {
    this.client = new CodeBuildClient({});
  }

  async update(version: SemanticVersion, componentVersion: SemanticVersion): Promise<void> {
    const command = new StartBuildCommand({
      projectName: this.projectName,
      environmentVariablesOverride: [
        {
          name: "VERSION",
          value: version.toString(),
          type: "PLAINTEXT",
        },
        {
          name: "COMPONENT_VERSION",
          value: componentVersion.toString(),
          type: "PLAINTEXT",
        },
      ],
    });
    const response = await this.client.send(command);
    logger.info("Started codebuild", { details: { codebuildResponse: response } });
  }
}
