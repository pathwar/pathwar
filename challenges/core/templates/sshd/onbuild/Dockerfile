# pathwar/sshd:latest
FROM ubuntu:vivid
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Install dependencies
RUN apt-get update \
 && apt-get -y -q upgrade \
 && apt-get -y -q install \
      curl \
      openssh-server \
      supervisor \
      gcc \
      mg \
      vim \
      emacs-nox \
 && apt-get clean

# Add Pathwar tools
COPY skeleton /pathwar

# Override pathwar scripts
RUN echo /usr/bin/supervisord > /pathwar/scripts/run


# Cleanup fancy /home
RUN rm -f /etc/skel/*

# Prepare environment for SSH
RUN mkdir -p /var/run/sshd

# Create 'bobby' user
RUN groupadd --gid 4242 bobby                    \
 && useradd                                      \
      --home /home/bobby                         \
      --gid bobby                                \
      --groups bobby                             \
      --uid 4242                                 \
      --shell /bin/bash                          \
      --password $(openssl passwd -crypt bobby)  \
      bobby                                      \
 && mkdir /home/bobby                            \
 && chown bobby:bobby /home/bobby
WORKDIR /home/bobby

# Create 'super-bobby' user
RUN groupadd --gid 8484 super-bobby              \
 && useradd                                      \
      --gid super-bobby                          \
      --groups super-bobby                       \
      --uid 8484                                 \
      --shell /bin/bash                          \
      super-bobby

COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf


# Onbuild
ONBUILD WORKDIR /home/bobby
ONBUILD CMD /pathwar/scripts/init
ONBUILD COPY scripts /pathwar/scripts/
ONBUILD COPY level.yml /pathwar/
ONBUILD EXPOSE 22

# ONBUILD COPY ...
