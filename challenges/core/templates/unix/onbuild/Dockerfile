# pathwar/unix:onbuild
FROM ubuntu:vivid
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)


# Install dependencies
RUN apt-get update \
 && apt-get -y -q upgrade \
 && apt-get -y -q install \
      curl \
      gcc \
      mg \
      vim \
      emacs-nox \
 && apt-get clean


# Add Pathwar tools
COPY skeleton /pathwar


# Override pathwar scripts
RUN echo /bin/bash > /pathwar/scripts/run


# Cleanup fancy ubuntu stuff
RUN rm -f /etc/skel/* \
 && rm -f /etc/legal \
 && rm -f /etc/update-motd.d/*


# Create 'bobby' user
RUN groupadd --gid 4242 bobby      \
 && useradd                        \
      --home /home/bobby           \
      --gid bobby                  \
      --groups bobby               \
      --uid 4242                   \
      --shell /bin/bash            \
      --create-home                \
      bobby                        \
 && mkdir -p /home/bobby           \
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


# Copy the content of ./overlay/ to /
COPY overlay /


# Onbuild
ONBUILD WORKDIR /home/bobby
ONBUILD USER bobby
ONBUILD CMD /pathwar/scripts/init
ONBUILD COPY scripts /pathwar/scripts/
ONBUILD COPY level.yml /pathwar/
