package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func flagUsage() {
	usageText := `Encode and Decode url.Values to json.

Usage:
    ujsonize [SUBCOMMAND] [OPTIONS]

SUBCOMMAND:
    encode   encode url.Values to json
    decode   decode json to url.Values

See ujsonize <command> -help for more information on a specific command.`
	fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
}

func encode(in []byte, out io.Writer) error {
	parsed, err := url.ParseQuery(string(in))
	if err != nil {
		return fmt.Errorf("failed to parse query: %w", err)
	}

	if err := json.NewEncoder(out).Encode(parsed); err != nil {
		return fmt.Errorf("failed to encode url.Values: %w", err)
	}
	return nil
}

func decode(in []byte, out io.Writer) error {
	tmp_uv := make(map[string]interface{})
	if err := json.Unmarshal(in, &tmp_uv); err != nil {
		return fmt.Errorf("failed to decode json: %w", err)
	}

	keys := make([]string, 0)
	uv := make(map[string][]string)
	for k, v := range tmp_uv {
		keys = append(keys, k)
		switch v.(type) {
		case []interface{}:
			uv[k] = make([]string, 0, len(v.([]interface{})))
			for _, vv := range v.([]interface{}) {
				uv[k] = append(uv[k], fmt.Sprintf("%v", vv))
			}
		case string:
			uv[k] = []string{v.(string)}
		default:
			tmp, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("failed to encode json: %w", err)
			}
			uv[k] = []string{string(tmp)}
		}
	}
	sort.Strings(keys)

	ret := make([]string, 0, len(uv))
	for _, k := range keys {
		params := uv[k]
		for _, param := range params {
			ret = append(ret, fmt.Sprintf("%s=%s", k, param))
		}
	}
	fmt.Fprintf(out, "%s\n", strings.Join(ret, "&"))
	return nil
}

type Config struct {
	ByLine bool
	NoTrim bool
}

func bindFlag(set *flag.FlagSet, config *Config) {
	set.BoolVar(&config.NoTrim, "no-trim", false, "don't trim whitespace")
	set.BoolVar(&config.ByLine, "by-line", false, "process by line")
}

func main() {
	flag.Usage = flagUsage
	encodeCmd := flag.NewFlagSet("encode", flag.ExitOnError)
	decodeCmd := flag.NewFlagSet("decode", flag.ExitOnError)
	flag.Parse()

	var (
		config Config
		fp     func([]byte, io.Writer) error
	)
	action := flag.Arg(0)
	switch action {
	case "encode":
		bindFlag(encodeCmd, &config)
		encodeCmd.Parse(flag.Args()[1:])
		fp = encode
	case "decode":
		bindFlag(decodeCmd, &config)
		decodeCmd.Parse(flag.Args()[1:])
		fp = decode
	default:
		flag.Usage()
		os.Exit(1)
	}

	if config.ByLine {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			in := scanner.Bytes()
			if !config.NoTrim {
				in = bytes.TrimSpace(in)
			}
			if err := fp(in, os.Stdout); err != nil {
				log.Printf("%s", err)
				continue
			}
		}
	} else {
		in, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		if !config.NoTrim {
			in = bytes.TrimSpace(in)
		}
		if err := fp(in, os.Stdout); err != nil {
			log.Fatal(err)
		}
	}
}
