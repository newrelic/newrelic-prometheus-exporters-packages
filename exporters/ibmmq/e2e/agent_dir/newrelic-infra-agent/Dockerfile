FROM ubuntu:22.04
USER root

RUN ln -sf /bin/true /bin/systemctl

RUN apt update && apt -y install curl wget gnupg && \
    echo "deb [arch=amd64] http://nr-downloads-main.s3-website-us-east-1.amazonaws.com/infrastructure_agent/linux/apt jammy main" > /etc/apt/sources.list.d/newrelic-infra.list  && \
    wget -nv -O- http://nr-downloads-main.s3-website-us-east-1.amazonaws.com/infrastructure_agent/gpg/newrelic-infra.gpg | apt-key add -  && \
    apt update && apt install -y newrelic-infra

# Create location for the MQ installation
RUN mkdir -p /opt/mqm \
  && chmod a+rx /opt/mqm

# Location of the downloadable MQ client package \
ENV RDURL="https://public.dhe.ibm.com/ibmdl/export/pub/software/websphere/messaging/mqdev/redist" \
    RDTAR="IBM-MQC-Redist-LinuxX64.tar.gz" \
    VRMF=9.2.5.0

# Install the MQ client from the Redistributable package. This also contains the
# header files we need to compile against. Setup the subset of the package
# we are going to keep - the genmqpkg.sh script removes unneeded parts
ENV genmqpkg_incnls=1 \
    genmqpkg_incsdk=1 \
    genmqpkg_inctls=1

RUN cd /opt/mqm \
 && curl -LO "$RDURL/$VRMF-$RDTAR" \
 && tar -zxf ./*.tar.gz \
 && rm -f ./*.tar.gz \
 && bin/genmqpkg.sh -b /opt/mqm \
 && mkdir -p /IBM/MQ/data/errors \
   && mkdir -p /.mqm \
   && chmod -R 777 /IBM \
   && chmod -R 777 /.mqm

ENV LD_LIBRARY_PATH="/opt/mqm/lib64:/usr/lib64" \
    MQ_CONNECT_TYPE=CLIENT

# IBMMQ Client needs permissions to write in a dir passed specified in HOME
# https://www.ibm.com/support/pages/apar/IT32860 
ENV HOME=/tmp
ENV NRIA_PASSTHROUGH_ENVIRONMENT=HOME

ENTRYPOINT ["/usr/bin/newrelic-infra-service"]
