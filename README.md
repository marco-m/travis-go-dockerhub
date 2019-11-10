# Travis Go DockerHub

[![Travis Build Status](https://travis-ci.org/marco-m/travis-go-dockerhub.svg?branch=master)](https://travis-ci.org/marco-m/travis-go-dockerhub)

How to use TravisCI to build and test Go code for a Docker image and publish the image to DockerHub. Supports issuing releases based on git tags. The Docker images are tagged following best practices (moving stable tags).

## Status

Work in progress.

### Currently missing features

- Moving stable tags, for example `1.2.3` is fixed, but `1.2` should give the latest release in the `1.2.x` series and `1` should give the latest release in the `1.x` series.
- **BUG**: If you support two series of releases, say `1.2.x` and `1.3.x`, then tag `latest` will flip between the two series! Currently tag `latest` works fine as long as you support only one series of releases.

## Secure setup for secrets

We need to give credentials to Travis to publish the Docker image to our DockerHub (or other Docker registry) account. We want to do it in such a way to preserve the secrecy of the credentials and to control what happens when somebody not belonging to the project issues a PR.

### DockerHub token setup

Do now use your DockerHub password, instead create a dedicated access token, see documentation at [dockerhub access tokens](https://docs.docker.com/docker-hub/access-tokens/). This allows to:

1. Reduce exposure (principle of least privilege), since a token has less capabilities than an account password.
2. Enable auditing of token usage.
3. Enable token revocation.

Unfortunately it is not possible to limit the scope of a token to a given image repository: a token has access to all repositories of an account. Nonetheless, it still makes sense to use a separate token per image repository, since it enables better auditing.

Login to your account and go to Settings | Security. Create a token, give it a name such as `Travis Project Foo` and securely back it up in your OS key store.

From an API point of view, the token can be used with `docker login` as if it was a password.

### Travis secrets setup

Please read the reference documentation [travis encryption-keys] before continuing.

The main idea is to store the secrets in the source repository (the repository containing the `.travis.yml` file), using the encrypted environment variables feature of Travis.

Note that this feature, for security reasons, does NOT make secure environment variables available to PRs coming from a forked source repository.

The [travis encryption-keys] documentation contains also pointers to the `travis` CLI. For macOS, `brew install travis` just works.

Do not follow the documentation example (`travis encrypt SOMEVAR="secretvalue"`) because it would leave the secrets in the shell history. Instead, run the tool in interactive mode with the `-i` flag:

```
$ cd the-repo
$ travis encrypt --add -i
Detected repository as marco-m/travis-go-dockerhub, is this correct? |yes|
Reading from stdin, press Ctrl+D when done
DOCKER_TOKEN="YOUR_TOKEN"  <= this is a real secret
THE_SECRET="42"            <= this shows how to pass additional secrets; see the tests
```

The `--add` will add the entry to the `.travis.yml` file.

All the encrypted secrets are defined in the `.travis.yml` file under key `secure`.

### Risk analysis

Due to the fact that a DockerHub token cannot be scoped to a specific image repository, the leaking of such token from CI (for example: the secret environment variable is not redacted and appears in logs, or user misconfiguration, or vulnerability in the secure environment mechanism) gives to an adversary write access to all image repositories of the given DockerHub account.

This problem is exacerbated if the source repository and especially the Travis build is publicly accessible.

## Releases, git and Docker tags

With reference to the corresponding [image repository on DockerHub](https://hub.docker.com/r/marcomm/travis-go-dockerhub):

Each time a commit is made on a branch, a new Docker image will be pushed with a tag corresponding to the branch name. This is done to enable integration testing of the image before merging. Said in another way: the Docker tag with the branch name always gives you an image built off the current tip of the branch.

Making a release.

If a branch is tagged (you should tag only the default branch), a new Docker image will be built and will be pushed with two tags:

1. Same Docker tag as the git tag, without the optional `v` prefix. For example, git tag `v1.2.3` will become DOcker tag `1.2.3`.
2. The `latest` Docker tag. If git tags are made only on the default branch, then Docker tag `latest` represents the latest release of the project, not the latest commit to the default branch. As such, it is as stable as pinning a specific release.

```
$ git tag -a -m 'Release 0.0.2' v0.0.2
git push origin v0.0.2
```

## Travis Docker build

See

* The Travis documentation [Using Docker in Builds](https://docs.travis-ci.com/user/docker/).
* The file `travis.yml` in this repo.
* The file `Taskfile.yml` in this repo.

## Local build

The same Taskfile can be used for local builds, for CI builds and inside a Docker container.

For easiness of customizations, you need to setup some environment variables.

Use a **secure means** to protect the environment variables, since some of them are sensitive, such as the DockerHub token!

We suggest to use [envchain](https://github.com/sorah/envchain) or [gopass](https://github.com/gopasspw/gopass).

### Setup envchain

```
envchain --set travis-docker DOCKER_USERNAME

$ envchain --set travis-docker DOCKER_TOKEN
travis-docker.DOCKER_TOKEN: YOUR_TOKEN_HERE

$ envchain --set travis-docker THE_SECRET
travis-docker.THE_SECRET: 42
```

### Build

```
$ envchain travis-docker task test
$ envchain travis-docker task build

$ envchain travis-docker task docker-build
$ envchain travis-docker task docker-smoke
$ envchain travis-docker task docker-push
```

[travis encryption-keys]: https://docs.travis-ci.com/user/encryption-keys/
