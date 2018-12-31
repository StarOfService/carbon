## Building a Carbon package without an application code
Sometimes you may need to build a Carbon package without a real application code. For example, it may be necessary, when you want to use community Docker image (say, RabbitMQ) but with custom configuration. It's still achievable with Carbon.

Docker image is a foundation for a Carbon package, but we want to avoid senseless size bloating, we recommend to use `scratch` as a parent image. It's a special base image, which doesn't have any data. Hence, you'll get the most lightweight Docker image possible, without any useless data. Taking into account all mentioned above, your full Dockerfile will be:
```
FROM scratch
```
It doesn't weight a single byte, but can be used for distribution of your Kubernetes manifests and allows you to use other amazing features of Carbon.
