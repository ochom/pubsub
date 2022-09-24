FROM rabbitmq:3.9-management

# Copy the plugin
COPY ./plugins/rabbitmq_delayed_message_exchange-3.9.0.ez /plugins/

RUN rabbitmq-plugins enable rabbitmq_delayed_message_exchange
