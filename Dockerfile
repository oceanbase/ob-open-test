FROM golang:1.18 as builder
WORKDIR /workspace
COPY bin .
COPY obopentest .
RUN chmod +x obopentest
