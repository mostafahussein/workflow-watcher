IMAGE_REPO            =  ghcr.io/mostafahussein/workflow-watcher
PACKAGE_NAME          := github.com/mostafahussein/workflow-watcher
GOLANG_CROSS_VERSION  ?= v1.20
SYSROOT_DIR           ?= sysroots
SYSROOT_ARCHIVE       ?= sysroots.tar.bz2

.PHONY: sysroot-pack
sysroot-pack:
	@tar cf - $(SYSROOT_DIR) -P | pv -s $[$(du -sk $(SYSROOT_DIR) | awk '{print $1}') * 1024] | pbzip2 > $(SYSROOT_ARCHIVE)

.PHONY: sysroot-unpack
sysroot-unpack:
	@pv $(SYSROOT_ARCHIVE) | pbzip2 -cd | tar -xf -

.PHONY: dry-run
dry-run:
	@docker run \
		--rm \
		-e CGO_ENABLED=1 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--clean --skip-validate --skip-publish --snapshot

.PHONY: lint
lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v1.52.2 golangci-lint run -v

.PHONY: build-release
build-release:
	docker run \
		--rm \
		-e CGO_ENABLED=1 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --clean --skip-publish

.PHONY: push
push:
	@if [ -z "$$VERSION" ]; then \
		echo "VERSION is required"; \
		exit 1; \
	fi
	export IMAGE_TAG=$(shell echo $$VERSION | sed -e s/^v//); \
	docker push $(IMAGE_REPO):$$IMAGE_TAG-amd64; \
	docker push $(IMAGE_REPO):$$IMAGE_TAG-arm64; \
	docker manifest create $(IMAGE_REPO):$$IMAGE_TAG \
	--amend $(IMAGE_REPO):$$IMAGE_TAG-amd64 \
	--amend $(IMAGE_REPO):$$IMAGE_TAG-arm64; \
	docker manifest push $(IMAGE_REPO):$$IMAGE_TAG;