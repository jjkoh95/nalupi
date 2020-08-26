# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
FROM alpine:3
RUN apk add --no-cache ca-certificates

# Copy built binary to production image and its necessary dependencies
COPY ./server ./
COPY ./credentials.json ./

CMD ["./server"]
