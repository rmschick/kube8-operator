FROM gcr.io/distroless/static:nonroot

WORKDIR /terminal-poc-deployment

ENV PATH /terminal-poc-deployment/bin:$PATH

# Copy the binary
COPY terminal-poc-deployment /bin/

ENTRYPOINT ["terminal-poc-deployment"]
