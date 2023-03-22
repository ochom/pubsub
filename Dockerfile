FROM rabbitmq:3.9-alpine
# FROM rabbitmq:3.9-management

# Copy the plugin
COPY plugins /plugins

# Enable delayed message exchange plugin
RUN rabbitmq-plugins enable rabbitmq_delayed_message_exchange

# Enable rabbitmq_management plugin
RUN rabbitmq-plugins enable rabbitmq_management

EXPOSE 5672

EXPOSE 15672

CMD ["rabbitmq-server"]