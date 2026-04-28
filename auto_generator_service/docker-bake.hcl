
variable "REPOSITORY" {}
variable "TAG" {}

target "default" {
    context = "."
    dockerfile = "Dockerfile"
    args = {
        GO_VERSION = "1.26.2"
        ALPINE_VERSION = "3.23"
    }
    platforms = ["linux/amd64"]

    tags = ["${REPOSITORY}:${TAG}"]
}
