FROM ubuntu:22.04 AS builder

RUN apt update && \
    apt install -yq wget curl git build-essential && \
    rm -rf /var/lib/apt/lists

ARG TARGETARCH
RUN curl -L https://go.dev/dl/go1.20.14.linux-${TARGETARCH}.tar.gz | tar -C /usr/local -xz

ENV PATH=/usr/local/go/bin/:${PATH}
ENV GOPATH /go
ENV GOCACHE=/root/.cache/go-build

WORKDIR /go/delivery/zeta-node

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .
RUN --mount=type=cache,target="/root/.cache/go-build" make install


# ============
#  Cosmovisor
# ============

FROM golang:1.20 AS cosmovisor-builder
RUN apt update && apt install -y bash clang tar wget musl-dev git make gcc bc ca-certificates

ARG GIT_REF=cosmovisor/v1.5.0
ARG REPO_URL=https://github.com/cosmos/cosmos-sdk
RUN git clone -n "${REPO_URL}" cosmos-sdk \
    && cd cosmos-sdk \
    && git fetch origin "${GIT_REF}" \
    && git reset --hard FETCH_HEAD

# WORKDIR /go/cosmos-sdk/cosmovisor/
WORKDIR /go/cosmos-sdk/

RUN go mod download
RUN make cosmovisor


FROM ubuntu:22.04

COPY contrib/docker-scripts/start.sh /scripts/start.sh

RUN chmod +x /scripts/start.sh

RUN apt update && \
    apt install -yq wget curl jq && \
    rm -rf /var/lib/apt/lists

COPY --from=builder /go/bin/zetaclientd /usr/local/bin/zetaclientd
COPY --from=builder /go/bin/zetacored /usr/local/bin/zetacored
COPY --from=cosmovisor-builder /go/cosmos-sdk/tools/cosmovisor /usr/local/bin/cosmovisor

EXPOSE 26656
EXPOSE 1317
EXPOSE 8545
EXPOSE 8546
EXPOSE 9090
EXPOSE 26657
EXPOSE 9091
