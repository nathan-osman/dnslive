## dnslive

[![Go Reference](https://pkg.go.dev/badge/github.com/nathan-osman/dnslive.svg)](https://pkg.go.dev/github.com/nathan-osman/dnslive)
[![MIT License](https://img.shields.io/badge/license-MIT-9370d8.svg?style=flat)](https://opensource.org/licenses/MIT)

This tool provides both a simple DNS server for dynamic IP addresses and a client application for providing the updates.

### Features

- No external dependencies
- Runs on any platform supported by Go
- Secured using HTTPS client authentication
- [Available as a Docker image](https://hub.docker.com/repository/docker/nathanosman/dnslive/general)

### Overview

Typical setup for dnslive looks like this:

1. A certificate authority and client / certificates are generated ([easy-rsa](https://github.com/OpenVPN/easy-rsa) is good for this)
1. Application is installed on a server and run with `dnslive server`
1. Application is installed on each client behind a dynamic IP address and run with `dnslive client` (each client provides a domain name that it wants to use)
1. The server responds to DNS requests for the domain names selected by the clients
1. As the clients' IP addresses change, the server receives the updates and adjusts the records accordingly

### Server Setup

Running the server is as simple as pointing it to a configuration file:

    dnslive server --config config.yaml

The configuration file is in YAML format and consists of the following keys:

| Key | Description | Default |
|---|---|---|
| `ca_cert_filename` | certificate authority used to sign the connecting clients ||
| `cert_filename` | certificate presented to connecting clients ||
| `key_filename` | private key for certificate ||
| `http_server_addr` | listen here for client connections | `0.0.0.0:443` |
| `dns_server_addr` | listen here for DNS requests | `0.0.0.0:53` |
| `persistent_file` | store address assignments in this file | `entries.json` |

### Client Setup

Running the client, likewise, is as simple as pointing it to a configuration file:

    dnslive client --config config.yaml

The configuration consists of the following keys:

| Key | Description | Default |
|---|---|---|
| `ca_cert_filename` | certificate authority used to verify the server ||
| `cert_filename` | certificate to present to the server ||
| `key_filename` | private key for certificate ||
| `server_addr` | address to use for the server ||
| `interval` | update interval for checking IP address(es) | `1h` |
| `name` | domain name to reserve on the server ||
