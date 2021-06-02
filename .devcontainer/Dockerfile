# Install Go version 1.16
ARG VARIANT="1.16"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# Install Node
ARG INSTALL_NODE="true"
ARG NODE_VERSION="lts/*"
RUN if [ "${INSTALL_NODE}" = "true" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

# Set up system dependencies
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends netcat tmux

# Install Go packages
# RUN go get -x <your-dependency-or-tool>

# Install Node packages
RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install" 2>&1
RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g ganache-cli truffle" 2>&1
