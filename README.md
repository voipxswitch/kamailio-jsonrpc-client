# kamailio-jsonrpc-client

`kamailio-jsonrpc-client` is a lightweight wrapper around kamailio's jsonrpcs module. It exposes a REST API endpoints which executes jsonrpc requests and returns the results. 

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
    xhttp_reply("404", "Not Found", "application/javascript", "{\"$var(y)\"}\n\r");
    return;
}
```

### htable dump

```
curl http://localhost:8080/v1/htable/dump?table=mytable
```

### htable flush

```
curl -X POST 'http://localhost:8080/v1/htable/mytable?action=flush'
```

### htable delete

```
curl -X DELETE 'http://localhost:8080/v1/htable/mytable/mykey'
```


### htable get

```
curl http://localhost:8080/v1/htable/mytable?key=mykey
```

### uac add registration
```
curl -X POST -d '{"id":"test123","username": "test123", "domain": "testdomain", "auth_username": "user01", "auth_password": "pass01", "proxy": "sip:5.6.6.7;transport=tcp", "random_delay": 10}' http://localhost:8080/v1/uacreg/register
```

### uac remove registration
```
curl -X POST 'http://localhost:8080/v1/uacreg/unregister?domain=testdomain&username=test123'
```

### uac list by all
```
curl http://localhost:8080/v1/uacreg/list
```

### uac list by domain
```
curl 'http://localhost:8080/v1/uacreg/list?domain=testdomain'
```

### uac list by user
```
curl 'http://localhost:8080/v1/uacreg/list?domain=testdomain&username=test123'
```
