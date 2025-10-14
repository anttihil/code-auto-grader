# Unified multi-language container for comparison
FROM rust:1.75-alpine

# Install Node.js, Python, and build dependencies
RUN apk add --no-cache \
    nodejs \
    npm \
    python3 \
    py3-pip \
    musl-dev

# Pre-compile Rust stdlib cache
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

CMD ["/bin/sh"]
