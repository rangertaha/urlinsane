// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rangertaha/urlinsane/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var downloadFlags = []cli.Flag{}

var DownloadCmd = cli.Command{
	Name:                   "download",
	Aliases:                []string{"d"},
	Usage:                  "Download datasets",
	UsageText:              "download [opt..] [directory]",
	UseShortOptionHandling: true,
	Flags:                  downloadFlags,
	Action: func(cCtx *cli.Context) error {
		if cCtx.NArg() == 0 {
			fmt.Println(text.FgRed.Sprint("\n  a directory is needed!\n"))
			cli.ShowSubcommandHelpAndExit(cCtx, 1)
		}

		if _, err := config.Init(); err != nil {
			log.Error(err)
			fmt.Println(text.FgRed.Sprint(err))
			cli.ShowSubcommandHelpAndExit(cCtx, 1)
		}
		return Download(cCtx)
	},
}

func Download(cli *cli.Context) error {
	folder := cli.Args().First()
	configDir := filepath.Join(folder, "domains")
	// Returned rather than logged. A directory that could not be created means
	// every download below fails too, and carrying on only turns one legible
	// error into a pile of them.
	if err := os.MkdirAll(configDir, 0o750); err != nil {
		return err
	}
	return DownloadSuffix(configDir)
}

// DownloadSuffix refreshes the public suffix list from publicsuffix.org.
//
// It returns an error rather than logging one. The previous version logged with
// logrus and carried on, which does not stop the function the way log.Fatal
// would: a failed request left resp nil and the next line dereferenced it, so
// the one case the check existed for was the case that panicked.
func DownloadSuffix(dirname string) error {
	fmt.Println("Downloading public suffix...")
	url := "https://publicsuffix.org/list/public_suffix_list.dat"

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download suffix: %w", err)
	}
	defer resp.Body.Close()

	// Without this an error page is data. publicsuffix.org answering 503 wrote
	// its HTML into suffix.lst, where every unindented line that is not a
	// comment reads as a public suffix.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download suffix: %s returned %s", url, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("download suffix: %w", err)
	}

	lines := strings.Split(string(body), "\n")
	var buffer bytes.Buffer

	rules := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "// ===BEGIN PRIVATE DOMAINS") {
			break
		}
		if line != "" && !strings.HasPrefix(line, "//") {
			line = strings.Replace(line, "*.", "", 1)
			line = strings.Replace(line, "!", "", 1)
			buffer.WriteString(line)
			buffer.WriteString("\n")
			rules++
		}
	}

	// A parse that yields nothing means the format moved, not that the list is
	// empty. Writing the empty result would replace a working suffix list with
	// a file that splits every domain wrongly, and the next run would have
	// nothing to fall back to.
	if rules == 0 {
		return fmt.Errorf("download suffix: %s parsed to no rules; the format has changed", url)
	}

	return replace(filepath.Join(dirname, "suffix.lst"), buffer.Bytes())
}

// replace writes a file atomically: a temporary beside the target, then a
// rename. A plain WriteFile truncates first, so an interrupted or short write
// leaves a half a suffix list that still loads and still parses -- silently,
// and wrongly, for every domain past the truncation point.
func replace(path string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name()) // no-op once the rename below succeeds

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	// Close before renaming, and check it: a write buffered by the kernel can
	// still fail here, and a rename of a file that failed to close publishes
	// the failure.
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), 0o640); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
