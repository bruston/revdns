revdns
======

Attempt to perform a reverse DNS lookup on a list of IP addresses.

## Usage

```
Usage of revdns:
  -c uint
    	number of concurrent requests to make (default 10)
  -f string
    	file containing list of ips, defaults to stdin if omitted
  -v	log errors to stdout
```

## Example

```
cat ips | revdns -c=15 | tee -a ip.hosts
```
