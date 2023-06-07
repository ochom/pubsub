FROM curlimages/curl:8.1.2 as Downloader

WORKDIR /downloaded

# Get the rabbitmq_delayed_message_exchange plugin from github
RUN curl -L -o rabbitmq_delayed_message_exchange-3.12.0.ez \
  https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases/download/v3.12.0/rabbitmq_delayed_message_exchange-3.12.0.ez


FROM rabbitmq:3.12.0-alpine as Runner

# Copy the downloaded plugin to the rabbitmq plugins directory
COPY --from=Downloader /downloaded/rabbitmq_delayed_message_exchange-3.12.0.ez /opt/rabbitmq/plugins/

# Enable delayed message exchange plugin
RUN rabbitmq-plugins enable rabbitmq_delayed_message_exchange

# Enable rabbitmq_management plugin
RUN rabbitmq-plugins enable rabbitmq_management

EXPOSE 5672

EXPOSE 15672

CMD ["rabbitmq-server"]
