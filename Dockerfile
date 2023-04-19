FROM gcr.io/distroless/static:nonroot
COPY workflow-watcher /usr/local/bin/workflow-watcher
CMD ["/usr/local/bin/workflow-watcher"]