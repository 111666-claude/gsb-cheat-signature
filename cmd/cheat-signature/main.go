// Command cheat-signature 跑作弊特征扫描样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/signatures"
)

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("cheat-signature", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "stale", "stale / dup / partial / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "stale":
		scanner := signatures.New()
		scanner.Scan("p", "s1", []string{"aim"})
		scanner.AddSignature("aim")
		fmt.Fprintf(stdout, "hits=%d\n", len(scanner.Scan("p", "s2", []string{"aim"})))
	case "dup":
		scanner := signatures.New()
		scanner.AddSignature("aim")
		scanner.Scan("p", "s1", []string{"aim"})
		scanner.Scan("p", "s1", []string{"aim"})
		fmt.Fprintf(stdout, "hits=%d\n", len(scanner.Hits("p")))
	case "partial":
		scanner := signatures.New()
		scanner.AddSignature("aimbot")
		fmt.Fprintf(stdout, "hits=%d\n", len(scanner.Scan("p", "s1", []string{"aim"})))
	case "scan":
		scanner := signatures.New()
		for index := 0; index < 3000; index++ {
			scanner.AddSignature(fmt.Sprintf("sig-%d", index))
		}
		for index := 0; index < 3000; index++ {
			scanner.Scan("p-0", fmt.Sprintf("s-%d", index), []string{"sig-0"})
		}
		fmt.Fprintf(stdout, "scanned=%d\n", scanner.Scanned())
	default:
		fmt.Fprintln(stderr, "需要 --sample stale|dup|partial|scan")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
