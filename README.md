# loopia-dns-updater

This is a small utility for updating the Loopia DNS records through
the Loopia API. It is capable of updating IPv4 A records and IPv6
AAAA records.

I'm using the tool myself and I have tested it, but please, make a backup of your
domain before running this utility. Also run a dry run first to check
that it behaves as expected.

## Design decisions

1. Must work with lots of hosts/records behind SLAAC
2. Rely on the IP of the interface rather than use an
   external service to find the IP
3. Auto discover of records in need to be updated
   (idea from one of my brothers)
4. Should work with both IPv4 and IPv6
5. The use case is home use and small office use with a single FW/Router
6. No config or Loopia credicals on every host with IPv6 services
7. Should be a small stand-alone binary

## Supported operating systems

The build script builds for Linux, FreeBSD and OpenBSD, but it should be
possible to compile for other operating systems and CPU architectures.

## Execution

The software runs as a single execution. Place it in cron or something similar
to run periodically

The util will store the relevant IPs of the interfaces. If they are not updated,
the util will not call Loopia to check whether update is needed. To force
execution, delete the data directory.

### Execution flags

| Flag      | Description                                                                                                     |
|-----------|-----------------------------------------------------------------------------------------------------------------|
| -c {file} | Config file location. If not set, /etc/loopia-dns-updater.yaml is used                                          |
| -d        | Debug output                                                                                                    |
| -h        | Show help                                                                                                       |
| -n        | Dry run                                                                                                         |
| -q        | Quiet. Minimal logging                                                                                          |
| -s        | Data directory. A directory where the utility store it's data. If not set, /var/lib/loopia-dns-updater is used. |

## Logging

Logging is performed to stdout.

## IPv4 handling

The IPv4 addresses are always treated as with netmask /32, i.e.
the complete IP of the interface is updated in the DNS.

## IPv6 handling

The IPv6 handling supports SLAAC. This is done by reading the netmask of
the interface. The prefix (most often the upper 64 bits) is replaced
in the AAAA record. The lower part (most often the lower 64 bits) is kept
as is. This method also works for /128 addresses where the full address
is replaced as with IPv4.

## Configuration file

A sample file is provided below

```yaml
domains:
  - name: example.com
    interfaces:
      - ifName: vio1
        matchUnknown4: true
        match4:
          - ssh.example.com
        match6:
          - ssh.example.com
      - ifName: vio5
        matchUnknown6: true

loopia:
  username: username-from-loopia-api
  password: password-of-loopia-api-user
```

Configuration items:

| Item                             | Description                                                                                             |
|----------------------------------|---------------------------------------------------------------------------------------------------------|
| domains                          | List of domains to update                                                                               |
| domains.name                     | FQDN of the domain                                                                                      |
| domains.interfaces               | List of interfaces where to find IP for the domain                                                      |
| domains.interfaces.ifName        | The interface name                                                                                      |
| domains.interfaces.match4        | List of FQDN that should be assigned the IPv4 of the interface                                          |
| domains.interfaces.match6        | List of FQDN that should be assigned the IPv6 of the interface                                          |
| domains.interfaces.matchUnknown4 | Update A records that are not listed in any {this}.interfaces.match4 with the IPv4 of this interface    |
| domains.interfaces.matchUnknown6 | Update AAAA records that are not listed in any {this}.interfaces.match6 with the IPv6 of this interface |
| domains.ignore                   | List of FQDN that should not be modified by the matchUnknown{4\|6} above                                |
| loopia.username                  | Loopia API user. The user will call GetSubdomains, GetZoneRecords and UpdateZoneRecord                  |
| loopia.password                  | Loopia API password                                                                                     |
| executeOnChange                  | Command to execute if the relevant IPs are changed                                                      |

Notes:

- The domain itself (the apex) is named `@` in Loopia, so write it as
  `@.example.com` in `match4`, `match6` and `ignore`. Writing
  `example.com` matches nothing.
- `executeOnChange` is split on spaces and run without a shell, so quotes,
  pipes, redirects and `&&` don't work. Put anything more than a simple
  command in a script and point `executeOnChange` at it.
