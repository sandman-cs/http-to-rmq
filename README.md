# http-to-rmq

This is a simple piece of middlewere that listens on a given port for http posts, including authentication tokens for security.  If the token is good, the payload is forwarded to RabbitMQ.

Example Configuration file "conf.json"
```
{
    "Broker":"freeversion.rmq.cloudamqp.com",
    "BrokerUser":"rmq.usr", 
    "BrokerPwd":"rmqPa55w0rd",
    "BrokerExchange":"amq.topic",
    "BrokerVhost":"/",
    "ChannelCount":4,
    "SrvPort":"443",
    "KeyFile": "/etc/http-rmq/privkey.pem",
    "CrtFile": "/etc/http-rmq/fullchain.pem"
}

```