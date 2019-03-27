vcl 4.1;

backend default {
    .host = "web-client";
    .port = "80";
}

backend productService {
    .host = "product-service";
    .port = "8080";
}

sub vcl_recv {
    if (req.url == "/product" || req.url == "/ws") {
        set req.backend_hint = productService;
    }
    if (req.http.upgrade ~ "(?i)websocket") {
        return (pipe);
    }
    if (req.method == "PURGE") {
        return (purge);
    }
}

sub vcl_pipe {
    if (req.http.upgrade) {
        set bereq.http.upgrade = req.http.upgrade;
        set bereq.http.connection = req.http.connection;
    }
}
