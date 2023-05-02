FROM gcr.io/distroless/static:nonroot

WORKDIR /locomotive-collector-template

ENV PATH /locomotive-collector-template/bin:$PATH

# Copy the binary
COPY locomotive-collector-template /bin/

ENTRYPOINT ["locomotive-collector-template"]
