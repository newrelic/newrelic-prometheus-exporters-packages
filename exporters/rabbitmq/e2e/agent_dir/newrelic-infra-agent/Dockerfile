FROM ubuntu:18.04
USER root

RUN ln -s /bin/true /usr/local/bin/systemctl

RUN apt update && apt -y install wget curl gnupg && \
    echo "deb [arch=amd64] http://nr-downloads-main.s3-website-us-east-1.amazonaws.com/infrastructure_agent/linux/apt bionic main" > /etc/apt/sources.list.d/newrelic-infra.list  && \
    wget -nv -O- http://nr-downloads-main.s3-website-us-east-1.amazonaws.com/infrastructure_agent/gpg/newrelic-infra.gpg | apt-key add -  && \
    apt update && apt install -y newrelic-infra

ADD integrations.d/nri-config.yml /etc/newrelic-infra/integrations.d
ADD bin/* /var/db/newrelic-infra/newrelic-integrations/bin/

RUN chmod +x /var/db/newrelic-infra/newrelic-integrations/bin/* || true

ENTRYPOINT ["/usr/bin/newrelic-infra-service"]