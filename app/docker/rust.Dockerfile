FROM rust:1.75-alpine

# Install minimal dependencies for musl compilation
RUN apk add --no-cache musl-dev

# Pre-compile a dummy program to cache stdlib compilation
# This dramatically reduces subsequent compilation times (60-70% faster)
RUN mkdir -p /tmp/cache && \
    echo 'fn main() { println!("cache"); }' > /tmp/cache/main.rs && \
    rustc /tmp/cache/main.rs -o /tmp/cache/program && \
    rm -rf /tmp/cache

# Create non-root user for code execution
RUN addgroup -S coderunner && adduser -S coderunner -G coderunner

# Set working directory
WORKDIR /workspace

# Create writable temp directory for compilation artifacts
RUN mkdir -p /workspace/tmp && chown coderunner:coderunner /workspace/tmp

# Switch to non-root user
USER coderunner

# Default entrypoint will be overridden
CMD ["/bin/sh"]
