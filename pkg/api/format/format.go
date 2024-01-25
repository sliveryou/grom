package format

import (
	"bufio"
	"errors"
	"go/format"
	"strings"

	"github.com/sliveryou/goctl/api/parser"
	"github.com/sliveryou/goctl/api/util"
	"github.com/sliveryou/goctl/util/pathx"
)

const (
	leftParenthesis  = "("
	rightParenthesis = ")"
	leftBrace        = "{"
	rightBrace       = "}"
)

// APIFormat formats the api data.
func APIFormat(data string) (string, error) {
	_, err := parser.ParseContentWithParserSkipCheckTypeDeclaration(data)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	s := bufio.NewScanner(strings.NewReader(data))
	tapCount := 0
	newLineCount := 0
	var preLine string
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 {
			if newLineCount > 0 {
				continue
			}
			newLineCount++
		} else {
			if preLine == rightBrace {
				builder.WriteString(pathx.NL)
			}
			newLineCount = 0
		}

		if tapCount == 0 {
			ft, err := formatGoTypeDef(line, s, &builder)
			if err != nil {
				return "", err
			}

			if ft {
				continue
			}
		}

		noCommentLine := util.RemoveComment(line)
		if noCommentLine == rightParenthesis || noCommentLine == rightBrace {
			tapCount--
		}
		if tapCount < 0 {
			line := strings.TrimSuffix(noCommentLine, rightBrace)
			line = strings.TrimSpace(line)
			if strings.HasSuffix(line, leftBrace) {
				tapCount++
			}
		}
		if line != "" {
			util.WriteIndent(&builder, tapCount)
		}
		builder.WriteString(line + pathx.NL)
		if strings.HasSuffix(noCommentLine, leftParenthesis) || strings.HasSuffix(noCommentLine, leftBrace) {
			tapCount++
		}
		preLine = line
	}

	return strings.TrimSpace(builder.String()), nil
}

func formatGoTypeDef(line string, scanner *bufio.Scanner, builder *strings.Builder) (bool, error) {
	noCommentLine := util.RemoveComment(line)
	tokenCount := 0
	if strings.HasPrefix(noCommentLine, "type") && (strings.HasSuffix(noCommentLine, leftParenthesis) ||
		strings.HasSuffix(noCommentLine, leftBrace)) {
		var typeBuilder strings.Builder
		typeBuilder.WriteString(mayInsertStructKeyword(line, &tokenCount) + pathx.NL)
		for scanner.Scan() {
			noCommentLine := util.RemoveComment(scanner.Text())
			typeBuilder.WriteString(mayInsertStructKeyword(scanner.Text(), &tokenCount) + pathx.NL)
			if noCommentLine == rightBrace || noCommentLine == rightParenthesis {
				tokenCount--
			}
			if tokenCount == 0 {
				ts, err := format.Source([]byte(typeBuilder.String()))
				if err != nil {
					return false, errors.New("error format \n" + typeBuilder.String())
				}

				result := strings.ReplaceAll(string(ts), " struct ", " ")
				result = strings.ReplaceAll(result, "type ()", "")
				builder.WriteString(result)
				break
			}
		}
		return true, nil
	}

	return false, nil
}

func mayInsertStructKeyword(line string, token *int) string {
	insertStruct := func() string {
		if strings.Contains(line, " struct") {
			return line
		}
		index := strings.Index(line, leftBrace)
		return line[:index] + " struct " + line[index:]
	}

	noCommentLine := util.RemoveComment(line)
	if strings.HasSuffix(noCommentLine, leftBrace) {
		*token++
		return insertStruct()
	}
	if strings.HasSuffix(noCommentLine, rightBrace) {
		noCommentLine = strings.TrimSuffix(noCommentLine, rightBrace)
		noCommentLine = util.RemoveComment(noCommentLine)
		if strings.HasSuffix(noCommentLine, leftBrace) {
			return insertStruct()
		}
	}
	if strings.HasSuffix(noCommentLine, leftParenthesis) {
		*token++
	}

	if strings.Contains(noCommentLine, "`") {
		return util.UpperFirst(strings.TrimSpace(line))
	}

	return line
}
