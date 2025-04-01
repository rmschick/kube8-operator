FROM gcr.io/distroless/static:nonroot

WORKDIR /kube8-operator

ENV PATH /kube8-operator/bin:$PATH

# Copy the binary
COPY kube8-operator /bin/

ENTRYPOINT ["kube8-operator"]
