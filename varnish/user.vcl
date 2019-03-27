vcl 4.1;

backend default {
    .host = "web-client";
    .port = "80";
}

backend productService {
    .host = "product-service";
    .port = "8080";
}
acl purge {
    "localhost";
    "127.0.0.1";
}
sub vcl_recv {
    if (req.url == "/product" || req.url == "/ws") {
        set req.backend_hint = productService;
    }
    if (req.http.upgrade ~ "(?i)websocket") {
        return (pipe);
    }
    if (req.method == "PURGE") {
        if (!client.ip ~ purge) {
            return(synth(405,"Not allowed."));
        }
        return (purge);
    }
}

sub vcl_pipe {
    if (req.http.upgrade) {
        set bereq.http.upgrade = req.http.upgrade;
        set bereq.http.connection = req.http.connection;
    }
}
