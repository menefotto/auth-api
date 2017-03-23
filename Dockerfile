# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM  ubuntu

# Copy the local package files to the container's workspace.
COPY auth-api /usr/local/bin
COPY api.conf /usr/local/bin
COPY waterandboards-go-bacaa0893077.json /etc/waterandboards-go-bacaa0893077.json
ENV GOOGLE_APPLICATION_CREDENTIALS=/etc/waterandboards-go-bacaa0893077.json
ENV AUTH_API_CONF=/usr/local/bin

RUN apt-get update -y
RUN apt-get install vim openssl -y
COPY ssl/ /etc/ssl
# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN cd /usr/local/bin
CMD ["auth-api"]
# Run the outyet command by default when the container starts.
EXPOSE 8080
# Document that the service listens on port 8080.
