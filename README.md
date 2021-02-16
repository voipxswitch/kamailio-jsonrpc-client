# kamailio-jsonrpc-client

### uac add registration
```
curl -X POST -d '{"username": "test123", "domain": "testdomain", "auth_username": "user01", "auth_password": "pass01", "proxy": "sip:5.6.6.7;transport=tcp", "random_delay": 10}' http://localhost:8080/v1/uac/register
```

### uac remove registration
```
curl -X POST 'http://localhost:8080/v1/uac/unregister?domain=testdomain&username=test123'
```

### uac list by all
```
curl http://localhost:8080/v1/uac/list
```

### uac list by domain
```
curl 'http://localhost:8080/v1/uac/list?domain=testdomain'
```

### uac list by user
```
curl 'http://localhost:8080/v1/uac/list?domain=testdomain&username=test123'
```

### kamailio config
```
#!substdef "!HTTP_PORT!8081!g"

tcp_accept_no_cl=yes
listen=tcp:127.0.0.1:HTTP_PORT

loadmodule "xhttp.so"
loadmodule "jsonrpcs.so"

modparam("jsonrpcs", "transport", 1)

event_route[xhttp:request] {
    if ($Rp != "HTTP_PORT") {
        xlog("L_WARN", "HTTP request received on $Rp from $si");
        xhttp_reply("403", "Forbidden", "", "");
        exit;
    }
    $var(x) = $(hu{url.path});
    xlog("L_INFO", "received request [$hu] [$var(x)]");
    if ($hu =~ "^/RPC") {
        xlog("L_INFO", "jsonrpc dispatch [$hu] [$var(x)]");
        jsonrpc_dispatch();
        return;
    }
    xhttp_reply("404", "OK", "application/javascript", "{\"$var(y)\"}\n\r");
    return;
}
```
