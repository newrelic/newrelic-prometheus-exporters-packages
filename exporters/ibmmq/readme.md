# Bump exporter runbook 
- Update commit in exporter.yml
- Update client library version
  - Update 'ibmmq_client_libs_version' in build-exporter-windows.sh. Use the same as the on [used](https://github.com/ibm-messaging/mq-metric-samples/blob/master/Dockerfile.build#L64) to build Linux 
  - Update client library version in the [e2e test](./e2e/agent_dir/newrelic-infra-agent/Dockerfile)
