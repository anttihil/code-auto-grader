FROM python:3.11-alpine

# Create non-root user for code execution
RUN addgroup -S coderunner && adduser -S coderunner -G coderunner

# Set working directory
WORKDIR /workspace

# Switch to non-root user
USER coderunner

# Default command will be overridden with the code file
CMD ["python3", "-u"]
