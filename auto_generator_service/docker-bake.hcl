
variable "REPOSITORY" {}
variable "TAG" {}

target "default" {
    context = "."
    dockerfile = "Dockerfile"
    args = {
        GO_VERSION = "1.25.5"
    }
    platforms = ["linux/amd64"]

    tags = ["${REPOSITORY}:${TAG}"]
}
