# Use Wolfi for the build stage
FROM cgr.dev/chainguard/go:latest-dev AS build-provider

# Copy source code
COPY . /app
WORKDIR /app

# Build the provider binary
RUN go build -v -o terraform-provider-glesys


# Final image with Terraform
FROM hashicorp/terraform:1.6.6

# Set plugin version and path variables
ENV PROVIDER_NAME=glesys \
    PROVIDER_NAMESPACE=github.com/norrland \
    PROVIDER_VERSION=0.3.2 \
    PROVIDER_OS=linux \
    PROVIDER_ARCH=amd64

# Construct plugin directory path
ENV PLUGIN_PATH="/root/.terraform.d/plugins/${PROVIDER_NAMESPACE}/${PROVIDER_NAME}/${PROVIDER_VERSION}/${PROVIDER_OS}_${PROVIDER_ARCH}"

# Create plugin directory
RUN mkdir -p ${PLUGIN_PATH}

# Copy provider binary from build stage
COPY --from=build-provider /app/terraform-provider-glesys ${PLUGIN_PATH}/terraform-provider-glesys

# Default working dir and entry
WORKDIR /home
ENTRYPOINT ["/bin/sh"]
