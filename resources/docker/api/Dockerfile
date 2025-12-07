FROM docker.io/golang:1.25.5-bookworm as builder

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

WORKDIR /app

FROM builder as deployer

COPY ./code/go .

ARG GIT_COMMIT
ENV GIT_COMMIT ${GIT_COMMIT}
ARG BUILD_DATE
ENV BUILD_DATE ${BUILD_DATE}
ARG GIT_BRANCH
ENV GIT_BRANCH ${GIT_BRANCH}
ARG BUILD_ID
ENV BUILD_ID ${BUILD_ID}
ARG GIT_TAG
ENV GIT_TAG ${GIT_TAG}

RUN make build/api

FROM gcr.io/distroless/base:nonroot AS deployable
USER 65532

ARG GIT_COMMIT
ARG BUILD_DATE
ARG GIT_BRANCH
ARG BUILD_ID
ARG GIT_TAG

LABEL \
  build-branch=${GIT_BRANCH} \
  build-commit=${GIT_COMMIT} \
  build-date=${BUILD_DATE} \
  build-id=${BUILD_ID} \
  build-version=${GIT_TAG}

COPY --from=deployer /app/bin/toll-api /toll-api

ENTRYPOINT ["/toll-api"]
