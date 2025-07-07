FROM ubuntu:22.04

# Gerekli paketleri kur
RUN apt-get update && \
    apt-get install -y bash curl wget git iputils-ping

# Go'yu indir ve kur
RUN wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz && \
    rm go1.22.3.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}" 