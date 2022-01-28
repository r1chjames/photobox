
docker-build:
ifneq (,$(wildcard /etc/fedora-release))
	podman build -t registry.gitlab.com/r1chjames/photobox .
else
	docker build -t registry.gitlab.com/r1chjames/photobox .
endif

docker-push: docker-build
ifneq (,$(wildcard /etc/fedora-release))
	CI_BUILD_TOKEN=$CI_BUILD_TOKEN podman login -u gitlab-ci-token -p ${CI_BUILD_TOKEN} registry.gitlab.com
	podman push registry.gitlab.com/r1chjames/photobox:$CI_COMMIT_BRANCH
else
	CI_BUILD_TOKEN=$CI_BUILD_TOKEN docker login -u gitlab-ci-token -p ${CI_BUILD_TOKEN} registry.gitlab.com
	docker push registry.gitlab.com/r1chjames/photobox:$CI_COMMIT_BRANCH
endif
