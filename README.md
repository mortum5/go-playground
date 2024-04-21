<p align="center">
  <img src="https://socialify.git.ci/mortum5/go-playground/image?description=1&descriptionEditable=&font=Inter&issues=1&language=1&name=1&owner=1&pattern=Signal&pulls=1&stargazers=1&theme=Light"     alt="go-playground" width="640" height="320" />
</p>

![Repository Top Language](https://img.shields.io/github/languages/top/mortum5/go-playground)
![Github Open Issues](https://img.shields.io/github/issues/mortum5/go-playground)
![GitHub contributors](https://img.shields.io/github/contributors/mortum5/go-playground)

## About

This repository contains all the snippets and small code samples I'm testing and exploring.

### Web Chat

Simple chat on web sockets with redis pubsub.

See [README](webchat/README.md) for additional information.

### RabbitMQ example

Just send and receive simple value via rabbitMQ. 

See [README](rabbitmq/README.md) for additional information.

### Simple brocker message

REST Service that handle two operations:

```
PUT /queue?v=msg     # Add message to queue
GET /queue?timeout=N # Retrive first message from queue or wait N seconds
```

See [README](message-broker/README.md) for additional information.

### Logger

Colorful slog configuration that decorates logs with pretty color.

See [README](logger/README.md) for additional information.

### Cluster

Simple cluster CRDT counter based on Hashicorp Memberlist library.

See [README](cluster/README.md) for additional information.

### Minio

File uploader on Fiber with Minio library.

See [README](minio/README.md) for additional information.