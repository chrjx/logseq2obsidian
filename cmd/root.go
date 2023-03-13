package cmd

import (
	"fmt"
	"github.com/chr11x/logseq2obisdian/scanner"
	"github.com/spf13/cobra"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var (
	obsidianPath string
)

var rootCmd = &cobra.Command{
	Use:   "logseq2obsidian",
	Short: "A converter for Logseq => Obsidian",
	Args:  cobra.ExactArgs(1),
	Run:   Logseq2Obsidian,
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringVarP(&obsidianPath, "output", "o", "out/", "output directory")
	rootCmd.MarkPersistentFlagRequired("output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}
}

func Logseq2Obsidian(cmd *cobra.Command, args []string) {
	logseqPageFolder := args[0]
	// add '/' if paths have no '/' at the end
	pages := make([]*scanner.Page, 0)
	// parse pages
	err := filepath.Walk(logseqPageFolder, func(pagePath string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			log.Printf("this is a directory: %s\n", pagePath)
			return nil
		}
		p := scanner.ParsePage(pagePath)
		pages = append(pages, p)
		log.Printf("Parsed page: %s\n", info.Name())
		return err
	})

	if err != nil {
		log.Println(err)
	}

	// write obsidian pages
	err = os.RemoveAll(obsidianPath)
	err = os.MkdirAll(obsidianPath, 0777)
	for _, p := range pages {
		if err = p.WriteInObsidian(obsidianPath); err != nil {
			log.Println(err)
			continue
		}
		log.Printf("Convert to %s", p.GetTitle())
	}
}
