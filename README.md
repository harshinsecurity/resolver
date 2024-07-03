# DNS Resolver

A high-performance, concurrent DNS resolver tool written in Go. This tool efficiently resolves domain names to IP addresses, handling various URL formats and providing flexible output options.

## Features

- Concurrent DNS resolution for high performance
- Handles various URL formats (with or without http/https, www prefix, etc.)
- Two output formats: IP-only (default) and domain-IP mapping
- Removes duplicate IPs in the default output mode
- Customizable concurrency level
- Simple command-line interface

## Installation

1. Ensure you have Go installed on your system. If not, download and install it from [golang.org](https://golang.org/).

2. Clone this repository:
   ```
   git clone https://github.com/yourusername/concurrent-dns-resolver.git
   cd concurrent-dns-resolver
   ```

3. Build the tool:
   ```
   go build dns_resolver.go
   ```

## Usage

Basic usage:
```
./dns_resolver
```

This will use default settings: reading from `urls.txt`, writing to `resolved_ips.txt`, using 100 concurrent workers, and outputting unique IPs only.

### Command-line Options

- `-input`: Specify the input file (default: "urls.txt")
- `-output`: Specify the output file (default: "resolved_ips.txt")
- `-concurrency`: Set the number of concurrent workers (default: 100)
- `-format`: Choose the output format: "ip" (default) or "domain-ip"
- `-help`: Display help information

### Examples

1. Use custom input and output files:
   ```
   ./dns_resolver -input=my_domains.txt -output=results.txt
   ```

2. Increase concurrency to 200 workers:
   ```
   ./dns_resolver -concurrency=200
   ```

3. Output domain-IP mappings instead of just IPs:
   ```
   ./dns_resolver -format=domain-ip
   ```

4. Combine multiple options:
   ```
   ./dns_resolver -input=sites.txt -output=ips.txt -concurrency=150 -format=domain-ip
   ```

## Input File Format

The input file should contain one URL or domain per line. The tool can handle various formats:

```
example.com
http://example.com
https://www.example.com
https://subdomain.example.com/path
```

## Output Formats

1. IP-only (default):
   - Outputs only unique, successfully resolved IP addresses
   - One IP per line
   - Unresolved domains are skipped

2. Domain-IP mapping:
   - Outputs "domain,ip" for resolved domains
   - Outputs "domain,Could not resolve" for unresolved domains

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
