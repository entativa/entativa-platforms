# Multi-stage Dockerfile for Scala services
# Stage 1: Build
FROM hseeberger/scala-sbt:11.0.16.1_1.8.2_2.13.10 AS builder

WORKDIR /app

# Copy build files
COPY build.sbt .
COPY project ./project

# Copy source code
COPY src ./src

# Build the application
RUN sbt clean assembly

# Stage 2: Runtime
FROM openjdk:11-jre-slim

WORKDIR /root/

# Copy JAR from builder
COPY --from=builder /app/target/scala-*/story-service-assembly-*.jar ./app.jar

# Expose ports
EXPOSE 8090 50005

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT:-8090}/health || exit 1

# Run the application
CMD ["java", "-jar", "./app.jar"]
