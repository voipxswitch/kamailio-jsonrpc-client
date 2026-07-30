# kamailio-jsonrpc-client

`kamailio-jsonrpc-client` is a lightweight wrapper around kamailio's jsonrpcs module. It exposes a REST API endpoints which executes jsonrpc requests and returns the results.

## kamailio config

```kamailio
#!substdef "!HTTP_PORT!8081!g"

tcp_accept_no_cl=yes
listen=tcp:127.0.0.1:HTTP_PORT

loadmodule "xhttp.so"
loadmodule "jsonrpcs.so"
loadmodule "jansson.so"

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
    if ($var(x) == "/htable/delete-prefix") {
        if ($rm != "POST") {
            xhttp_reply("405", "Method Not Allowed", "application/json", "");
            return;
        }
        $var(htable) = $null;
        $var(key_prefix) = $null;
        if (!jansson_get("htable", "$rb", "$var(htable)")
                || !jansson_get("key_starts_with", "$rb", "$var(key_prefix)")
                || $var(htable) == $null
                || $var(key_prefix) == $null
                || $var(htable) == ""
                || $var(key_prefix) == "") {
            xhttp_reply("400", "Bad Request", "application/json", "");
            return;
        }
        if (!sht_has_name("$var(htable)", "sw", "$var(key_prefix)")) {
            xhttp_reply("404", "Not Found", "application/json", "");
            return;
        }
        if (!sht_rm_name("$var(htable)", "sw", "$var(key_prefix)")) {
            xhttp_reply("500", "Internal Server Error", "application/json", "");
            return;
        }
        xhttp_reply("204", "No Content", "", "");
        return;
    }
    xhttp_reply("404", "Not Found", "application/javascript", "{\"$var(y)\"}\n\r");
    return;
}
```

### htable dump

```bash
curl http://localhost:8080/v1/htable/dump?table=mytable
```

### htable flush

```bash
curl -X POST 'http://localhost:8080/v1/htable/mytable?action=flush'
```

### htable delete

```bash
curl -X DELETE 'http://localhost:8080/v1/htable/mytable/mykey'
```

### htable delete by key prefix

```bash
curl -X DELETE 'http://localhost:8080/v1/htable/mytable?key_starts_with=mykey'
```

### htable get

```bash
curl http://localhost:8080/v1/htable/mytable?key=mykey
```

### uac add registration

```bash
curl -X POST -d '{"id":"test123","username": "test123", "domain": "testdomain", "auth_username": "user01", "auth_password": "pass01", "proxy": "sip:5.6.6.7;transport=tcp", "random_delay": 10}' http://localhost:8080/v1/uacreg/register
```

### uac remove registration

```bash
curl -X POST 'http://localhost:8080/v1/uacreg/unregister?domain=testdomain&username=test123'
```

### uac list by all

```bash
curl http://localhost:8080/v1/uacreg/list
```

### uac list by domain

```bash
curl 'http://localhost:8080/v1/uacreg/list?domain=testdomain'
```

### uac list by user

```bash
curl 'http://localhost:8080/v1/uacreg/list?domain=testdomain&username=test123'
```
