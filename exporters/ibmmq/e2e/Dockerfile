FROM ibmcom/mq:latest

ENV LICENSE=accept
ENV MQ_QMGR_NAME=QM1
ENV MQ_ENABLE_METRICS=true
ENV MQ_ENABLE_EMBEDDED_WEB_SERVER=1

EXPOSE 1414 9443

USER root
RUN groupadd mqm
RUN useradd newrelic -G mqm && \
    echo newrelic:passw0rd | chpasswd
COPY config.mqsc /etc/mqm/
