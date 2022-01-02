# retry

## Name

Plugin *Retry* is able to selectively fanout the query to another upstream server, depending the error result provided by the initial resolver

## Description

The *retry* plugin allows set of upstreams be specified which will be used
if the plugin chain returns specific error messages. The *retry* plugin utilizes the *fanout* plugin (<https://github.com/networkservicemesh/fanout>) to sent fanout query the specified upstreams.



## Syntax

```
{
    retry [original] RCODE_1[,RCODE_2,RCODE_3...] . DNS_RESOLVERS
}
```

* **original** is optional flag. If it is set then retry uses original request instead of potentially changed by other plugins
* **RCODE** is the string representation of the error response code. The complete list of valid rcode strings are defined as `RcodeToString` in <https://github.com/miekg/dns/blob/master/msg.go>, examples of which are `SERVFAIL`, `NXDOMAIN` and `REFUSED`. At least one rcode is required, but multiple rcodes may be specified, delimited by commas.
* **DNS_RESOLVERS** accepts dns resolvers list.

## Building CoreDNS with Retry

When building CoreDNS with this plugin, _retry_ should be positioned **before** _forward_ in `/plugin.cfg`.

## Examples

### Retry to local DNS server

The following specifies that all requests are forwarded to 8.8.8.8. If the response is `NXDOMAIN`, *retry* will fanout the request to 192.168.1.1:53, and reply to client accordingly.

```
. {
	forward . 8.8.8.8
	retry NXDOMAIN . 192.168.1.1:53
	log
}

```
### Retry with original request used

The following specify that `original` query will be forwarded to 192.168.1.1:53 if 8.8.8.8 response is `NXDOMAIN`. `original` means no changes from next plugins on request. With no `original` flag retry will forward request with EDNS0 option (set by rewrite).

```
. {
	forward . 8.8.8.8
	rewrite edns0 local set 0xffee 0x61626364
	retry original NXDOMAIN . 192.168.1.1:53
	log
}

```

### Multiple retries

Multiple retries can be specified, as long as they serve unique error responses.

```
. {
    forward . 8.8.8.8
    retry NXDOMAIN . 192.168.1.1:53
    retry original SERVFAIL,REFUSED . 192.168.100.1:53
    log
}

```
