FROM rabbitmq:3.9-management

# Copy the plugin
COPY plugins /plugins

# Enable delayed message exchange plugin
RUN rabbitmq-plugins enable rabbitmq_delayed_message_exchange

EXPOSE 5672 15672