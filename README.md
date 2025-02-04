# WebRTC gRPC Backend

This project is a Go backend that provides a gRPC API for initiating a WebRTC video session and handling chat messages. It is designed to be used with a Flutter application.

## Prerequisites

- Go (version 1.18 or later)
- Docker
- Protocol Buffers compiler (\`protoc\`) with Go plugins

## Generate gRPC Code

From the project root, run:
\`\`\`bash
protoc --go_out=. --go-grpc_out=. proto/session.proto
\`\`\`
This will generate the necessary Go files in the \`sessionpb\` package.

## Running Locally with Docker

Build the Docker image:
\`\`\`bash
docker build -t webrtc-backend .
\`\`\`

Run the container:
\`\`\`bash
docker run -p 50051:50051 webrtc-backend
\`\`\`

The gRPC server will now be listening on port 50051.

## Integration with Flutter

On your Flutter side, you can use the generated gRPC client to call \`StartSession\`, \`ConnectSession\`, and maintain a chat stream via the \`Chat\` method.
