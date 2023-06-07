FROM debian:bookworm-slim as Downloader

WORKDIR /downloaded

# Install curl
RUN apt-get update && apt-get install -y curl

# Get the rabbitmq_delayed_message_exchange plugin from github
RUN curl -L -o ./rabbitmq_delayed_message_exchange-3.12.0.ez \
  https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases/download/v3.12.0/rabbitmq_delayed_message_exchange-3.12.0.ez


FROM rabbitmq:3.12.0 as Runner

# Copy the downloaded plugin to the rabbitmq plugins directory
COPY --from=Downloader /downloaded/rabbitmq_delayed_message_exchange-3.12.0.ez /opt/rabbitmq/plugins/

# Enable delayed message exchange plugin
RUN rabbitmq-plugins enable rabbitmq_delayed_message_exchange

# Enable rabbitmq_management plugin
RUN set eux; \
  rabbitmq-plugins enable --offline rabbitmq_management; \
  # make sure the metrics collector is re-enabled (disabled in the base image for Prometheus-style metrics by default)
  rm -f /etc/rabbitmq/conf.d/20-management_agent.disable_metrics_collector.conf; \
  # grab "rabbitmqadmin" from inside the "rabbitmq_management-X.Y.Z" plugin folder
  # see https://github.com/docker-library/rabbitmq/issues/207
  cp /plugins/rabbitmq_management-*/priv/www/cli/rabbitmqadmin /usr/local/bin/rabbitmqadmin; \
  [ -s /usr/local/bin/rabbitmqadmin ]; \
  chmod +x /usr/local/bin/rabbitmqadmin; \
  apt-get update; \
  apt-get install -y --no-install-recommends python3; \
  rm -rf /var/lib/apt/lists/*; \
  rabbitmqadmin --version

EXPOSE 5672 15672 15671

CMD ["rabbitmq-server"]
