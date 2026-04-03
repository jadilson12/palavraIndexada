FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
  curl \
  wget \
  git \
  build-essential \
  ca-certificates \
  gnupg \
  lsb-release \
  software-properties-common \
  && rm -rf /var/lib/apt/lists/*

RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz && \
  tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz && \
  rm go1.24.0.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin

RUN ARCH=$(dpkg --print-architecture) && \
  if [ "$ARCH" = "arm64" ]; then \
    wget https://github.com/gohugoio/hugo/releases/download/v0.159.2/hugo_extended_0.159.2_linux-arm64.deb && \
    dpkg -i hugo_extended_0.159.2_linux-arm64.deb && \
    rm hugo_extended_0.159.2_linux-arm64.deb; \
  else \
    wget https://github.com/gohugoio/hugo/releases/download/v0.159.2/hugo_extended_0.159.2_linux-amd64.deb && \
    dpkg -i hugo_extended_0.159.2_linux-amd64.deb && \
    rm hugo_extended_0.159.2_linux-amd64.deb; \
  fi

WORKDIR /app

COPY go.mod go.sum ./
COPY hugo.yaml ./

RUN go mod download

COPY . .

RUN go run scripts/fix_html_entities.go && \
    go run scripts/generate_index.go

RUN hugo --gc --minify

EXPOSE 1313

CMD ["hugo", "server", "--logLevel", "debug", "--disableFastRender", "-p", "1313", "--bind", "0.0.0.0"]
