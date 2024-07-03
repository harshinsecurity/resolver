package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
)

type Result struct {
	InputURL string
	Domain   string
	IP       string
}

func extractDomain(input string) string {
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		input = "http://" + input
	}
	u, err := url.Parse(input)
	if err != nil {
		return input
	}
	return strings.TrimPrefix(u.Hostname(), "www.")
}

func resolveDomain(domain string) string {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return ""
	}
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}
	return ""
}

func worker(jobs <-chan string, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for inputURL := range jobs {
		domain := extractDomain(inputURL)
		ip := resolveDomain(domain)
		results <- Result{InputURL: inputURL, Domain: domain, IP: ip}
	}
}

func main() {
	inputFile := flag.String("input", "urls.txt", "Input file containing URLs or domains (one per line)")
	outputFile := flag.String("output", "resolved_ips.txt", "Output file to write results")
	concurrency := flag.Int("concurrency", 100, "Number of concurrent workers")
	outputFormat := flag.String("format", "ip", "Output format: 'ip' (default) or 'domain-ip'")
	helpFlag := flag.Bool("help", false, "Display help information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Concurrent Resolver Tool\n\n")
		fmt.Fprintf(os.Stderr, "This tool concurrently resolves domain names to IP addresses. It can handle various URL formats and extracts the domain for resolution.\n\n")
		fmt.Fprintf(os.Stderr, "Usage: resolver [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nOutput Formats:\n")
		fmt.Fprintf(os.Stderr, "  ip: Only outputs unique, successfully resolved IP addresses (default)\n")
		fmt.Fprintf(os.Stderr, "  domain-ip: Outputs 'domain,ip' for resolved domains and 'domain,Could not resolve' for failures\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  resolver\n")
		fmt.Fprintf(os.Stderr, "  resolver -input=my_domains.txt -output=results.txt -concurrency=200 -format=domain-ip\n")
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	if *outputFormat != "domain-ip" && *outputFormat != "ip" {
		fmt.Println("Invalid output format. Use 'ip' or 'domain-ip'.")
		return
	}

	input, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Error opening input file %s: %v\n", *inputFile, err)
		return
	}
	defer input.Close()

	output, err := os.Create(*outputFile)
	if err != nil {
		fmt.Printf("Error creating output file %s: %v\n", *outputFile, err)
		return
	}
	defer output.Close()

	jobs := make(chan string, *concurrency)
	results := make(chan Result, *concurrency)

	var wg sync.WaitGroup
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	go func() {
		scanner := bufio.NewScanner(input)
		for scanner.Scan() {
			jobs <- strings.TrimSpace(scanner.Text())
		}
		close(jobs)
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading input file: %v\n", err)
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	writer := bufio.NewWriter(output)
	defer writer.Flush()

	uniqueIPs := make(map[string]bool)

	for result := range results {
		var line string
		if *outputFormat == "domain-ip" {
			if result.IP != "" {
				line = fmt.Sprintf("%s,%s\n", result.Domain, result.IP)
			} else {
				line = fmt.Sprintf("%s,Could not resolve\n", result.Domain)
			}
			fmt.Print(line)
			writer.WriteString(line)
		} else { // "ip" format
			if result.IP != "" && !uniqueIPs[result.IP] {
				uniqueIPs[result.IP] = true
				line = fmt.Sprintf("%s\n", result.IP)
				fmt.Print(line)
				writer.WriteString(line)
			}
		}
	}
}
