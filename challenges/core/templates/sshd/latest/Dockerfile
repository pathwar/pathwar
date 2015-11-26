# pathwar/sshd:latest
FROM ubuntu:wily
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


# Cleanup fancy ubuntu stuff
RUN rm -f /etc/skel/* \
 && rm -f /etc/legal \
 && rm -f /etc/update-motd.d/*

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
ENV SU_USER bobby
WORKDIR /home/bobby

# Create 'super-bobby' user
RUN groupadd --gid 8484 super-bobby              \
 && useradd                                      \
      --gid super-bobby                          \
      --groups super-bobby                       \
      --uid 8484                                 \
      --shell /bin/bash                          \
      super-bobby

# Copy the content of ./overlay/ to /
COPY overlay /
