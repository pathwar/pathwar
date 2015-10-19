# pathwar/unix:onbuild
FROM ubuntu:vivid
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)


# Install dependencies
RUN apt-get update \
 && apt-get -y -q upgrade \
 && apt-get -y -q install \
      curl \
      emacs-nox \
      gcc \
      mg \
      sudo \
      vim \
 && apt-get clean


# Add Pathwar tools
COPY skeleton /pathwar


# Override pathwar scripts
RUN echo 'echo "level created, now you can enter with `level-enter`"; exit 0' > /pathwar/scripts/run


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
 && chown bobby:bobby /home/bobby  \
 && chmod 755 /home/bobby


WORKDIR /home/bobby


# Create 'super-bobby' user
RUN groupadd --gid 8484 super-bobby                  \
 && useradd                                          \
      --home /home/super-bobby                       \
      --gid super-bobby                              \
      --groups super-bobby                           \
      --uid 8484                                     \
      --shell /bin/bash                              \
      --no-create-home                               \
      super-bobby                                    \
 && mkdir -p /home/super-bobby                       \
 && chown super-bobby:super-bobby /home/super-bobby  \
 && chmod 755 /home/super-bobby


# Copy the content of ./overlay/ to /
COPY overlay /


# Onbuild
ONBUILD WORKDIR /home/bobby
ONBUILD CMD /pathwar/scripts/init
ONBUILD COPY scripts /pathwar/scripts/
ONBUILD COPY level.yml /pathwar/
