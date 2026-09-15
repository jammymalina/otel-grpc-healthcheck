
variable "REPOSITORY" {}
variable "TAG" {}

target "default" {
    context = "."
    dockerfile = "Dockerfile"
    args = {
        GO_VERSION = "1.27.1"
        ALPINE_VERSION = "3.24"
    }
    platforms = ["linux/amd64"]

    tags = ["${REPOSITORY}:${TAG}"]
}
