// Go version
variable "GO_VERSION" {
  default = "1.16.3"
}

// GitHub reference as defined in GitHub Actions (eg. refs/head/master))
variable "GITHUB_REF" {
  default = ""
}

// Common args for inheritance
target "_common" {
  args = {
    GO_VERSION = GO_VERSION
    GIT_REF = GITHUB_REF
  }
}

// Special target: https://github.com/docker/metadata-action#bake-definition
target "docker-metadata-action" {
  tags = ["psampaz/go-mod-outdated:local"]
}

group "default" {
  targets = ["image-local"]
}

group "validate" {
  targets = ["lint", "vendor-validate"]
}

target "lint" {
  inherits = ["_common"]
  dockerfile = "./hack/lint.Dockerfile"
  target = "lint"
  output = ["type=cacheonly"]
}

target "vendor-validate" {
  inherits = ["_common"]
  dockerfile = "./hack/vendor.Dockerfile"
  target = "validate"
  output = ["type=cacheonly"]
}

target "vendor-update" {
  inherits = ["_common"]
  dockerfile = "./hack/vendor.Dockerfile"
  target = "update"
  output = ["."]
}

target "test" {
  inherits = ["_common"]
  dockerfile = "./hack/test.Dockerfile"
  target = "test-coverage"
  output = ["."]
}

target "artifact" {
  inherits = ["_common"]
  target = "artifacts"
  output = ["./dist"]
}

target "artifact-all" {
  inherits = ["artifact"]
  platforms = [
    "darwin/amd64",
    "darwin/arm64",
    "linux/amd64",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/arm64",
    "windows/amd64"
  ]
}

target "image" {
  inherits = ["_common", "docker-metadata-action"]
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/arm64"
  ]
}
