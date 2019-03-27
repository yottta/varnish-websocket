# varnish-websocket

This repo is intended to give a startup point for a Varnish configuration to enable WebSocket communication.


## How to run this?

Clone the repo and run the following command:
    
    docker-compose up --build

Access http://localhost and you should see a list of two products.
Add a new product using curl:

    curl -XPOST http://127.0.0.1/product -d 'product3'

Invalidate cache:

    curl -XPURGE http://127.0.0.1/product -v
