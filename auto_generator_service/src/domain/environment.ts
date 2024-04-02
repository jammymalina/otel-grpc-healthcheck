export default class Environment {
  public readonly projectName: string;

  constructor() {
    this.projectName = Environment.required("PROJECT_NAME");
  }

  private static required(name: string): string {
    const value = process.env[name];
    if (value === undefined) {
      throw Error(`Missing the environment variable ${name}`);
    }
    return value;
  }
}
