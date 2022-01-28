
docker-build:
ifneq (,$(wildcard /etc/fedora-release))
	CI_COMMIT_BRANCH=$CI_COMMIT_BRANCH podman build -t registry.gitlab.com/r1chjames/photobox:$(CI_COMMIT_BRANCH) .
else
	CI_COMMIT_BRANCH=$CI_COMMIT_BRANCH docker build -t registry.gitlab.com/r1chjames/photobox:$(CI_COMMIT_BRANCH) .
endif

docker-push: docker-build
ifneq (,$(wildcard /etc/fedora-release))
	CI_BUILD_TOKEN=$CI_BUILD_TOKEN podman login -u gitlab-ci-token -p ${CI_BUILD_TOKEN} registry.gitlab.com
	CI_COMMIT_BRANCH=$CI_COMMIT_BRANCH podman push registry.gitlab.com/r1chjames/photobox:$CI_COMMIT_BRANCH:$(CI_COMMIT_BRANCH)
else
	CI_BUILD_TOKEN=$CI_BUILD_TOKEN docker login -u gitlab-ci-token -p ${CI_BUILD_TOKEN} registry.gitlab.com
	CI_COMMIT_BRANCH=$CI_COMMIT_BRANCH docker push registry.gitlab.com/r1chjames/photobox:$(CI_COMMIT_BRANCH)
endif
