import { SSMClient, GetParameterCommand } from "@aws-sdk/client-ssm";

export default class SecretService {
  private readonly client: SSMClient;

  constructor() {
    this.client = new SSMClient({});
  }

  async getSecret(key: string): Promise<string> {
    const input = {
      Name: key,
      WithDecryption: true,
    };
    const command = new GetParameterCommand(input);
    const response = await this.client.send(command);
    const value = response.Parameter?.Value;
    if (typeof value === "undefined") {
      throw new Error(`Unable to find secret: ${key}`);
    }
    return value;
  }
}
