package util

import (
	"math/rand"
	"regexp"
	"strings"
	"unicode"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
)

var (
	ignoreComments = regexp.MustCompile(`^--.*\s*`)
)

type FmtSqlConfig struct {
	PrintWidth int // default: 60
	UseSpaces  bool
	TabWidth   int    // default: 4
	CaseMode   string //be one of: upper, lower, title, spongebob, default: upper
	Align      bool
	NoSimplify bool //don't simplify the output
}

func DefaultFmtSqlConfig() FmtSqlConfig {
	return FmtSqlConfig{
		PrintWidth: 60,
		UseSpaces:  false,
		TabWidth:   4,
		CaseMode:   "upper",
		Align:      false,
	}
}

func FmtSQL(cfg FmtSqlConfig, stmts []string) (string, error) {
	cfg0 := tree.DefaultPrettyCfg()
	cfg0.UseTabs = cfg.UseSpaces
	if cfg.PrintWidth <= 0 {
		cfg0.LineWidth = 60
	} else {
		cfg0.LineWidth = cfg.PrintWidth
	}
	if cfg.TabWidth <= 0 {
		cfg0.TabWidth = 4
	} else {
		cfg0.TabWidth = cfg.TabWidth
	}
	cfg0.Simplify = !cfg.NoSimplify
	cfg0.Align = tree.PrettyNoAlign
	if cm, ok := caseModes[cfg.CaseMode]; ok {
		cfg0.Case = cm
	} else {
		cfg0.Case = caseModes["upper"]
	}
	cfg0.JSONFmt = true
	if cfg.Align {
		cfg0.Align = tree.PrettyAlignAndDeindent
	}

	var prettied strings.Builder
	for _, stmt := range stmts {
		for len(stmt) > 0 {
			stmt = strings.TrimSpace(stmt)
			hasContent := false
			// Trim comments, preserving whitespace after them.
			for {
				found := ignoreComments.FindString(stmt)
				if found == "" {
					break
				}
				// Remove trailing whitespace but keep up to 2 newlines.
				prettied.WriteString(strings.TrimRightFunc(found, unicode.IsSpace))
				newlines := strings.Count(found, "\n")
				if newlines > 2 {
					newlines = 2
				}
				prettied.WriteString(strings.Repeat("\n", newlines))
				stmt = stmt[len(found):]
				hasContent = true
			}
			// Split by semicolons
			next := stmt
			if pos, _ := parser.SplitFirstStatement(stmt); pos > 0 {
				next = stmt[:pos]
				stmt = stmt[pos:]
			} else {
				stmt = ""
			}
			// This should only return 0 or 1 responses.
			allParsed, err := parser.Parse(next)
			if err != nil {
				return "", err
			}
			for _, parsed := range allParsed {
				pretty, err := cfg0.Pretty(parsed.AST)
				if err != nil {
					return "", err
				}
				prettied.WriteString(pretty)
				prettied.WriteString(";\n")
				hasContent = true
			}
			if hasContent {
				prettied.WriteString("\n")
			}
		}
	}

	return strings.TrimRightFunc(prettied.String(), unicode.IsSpace), nil
}

var caseModes = map[string]func(string) string{
	"upper":     strings.ToUpper,
	"lower":     strings.ToLower,
	"title":     titleCase,
	"spongebob": spongeBobCase,
}

func titleCase(s string) string {
	return strings.Title(strings.ToLower(s))
}

func spongeBobCase(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, c := range s {
		b.WriteRune(unicode.To(rand.Intn(2), c))
	}
	return b.String()
}
