# Vignette Notification Service Dockerfile
FROM hseeberger/scala-sbt:11.0.16.1_1.9.0_2.13.11 AS builder

WORKDIR /app

# Copy build files
COPY build.sbt ./
COPY project ./project

# Download dependencies
RUN sbt update

# Copy source code
COPY src ./src

# Build fat JAR
RUN sbt assembly

# Final stage
FROM openjdk:11-jre-slim

WORKDIR /app

# Copy JAR from builder
COPY --from=builder /app/target/scala-2.13/vignette-notification-service-assembly-1.0.0.jar ./notification-service.jar

# Expose port
EXPOSE 8089

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8089/health || exit 1

# Run
CMD ["java", "-jar", "notification-service.jar"]
