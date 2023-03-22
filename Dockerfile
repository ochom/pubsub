FROM rabbitmq:3.9-alpine

# Install curl
RUN apk add --no-cache curl

# Get the rabbitmq_delayed_message_exchange plugin from github
RUN curl -L -o /plugins/rabbitmq_delayed_message_exchange-3.9.0.ez \
  https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases/download/3.9.0/rabbitmq_delayed_message_exchange-3.9.0.ez


# Enable delayed message exchange plugin
RUN rabbitmq-plugins enable rabbitmq_delayed_message_exchange

# Enable rabbitmq_management plugin
RUN rabbitmq-plugins enable rabbitmq_management

EXPOSE 5672

EXPOSE 15672

CMD ["rabbitmq-server"]