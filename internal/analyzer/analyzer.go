package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var configPath string

var Analyzer = &analysis.Analyzer{
	Name:     "loglint",
	Doc:      "checks logs",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func init() {
	Analyzer.Flags.StringVar(&configPath, "config", ".loglint.yml", "path to loglint config file")
}

func run(pass *analysis.Pass) (interface{}, error) {
	cfg := loadConfig(configPath)

	inspectResult := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspectResult.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		if !isSupportedLogCall(call, cfg) {
			return
		}

		if len(call.Args) == 0 {
			return
		}

		message, ok := getStringExpression(call.Args[0])
		if !ok {
			return
		}

		checkMessage(pass, call.Args[0], message, cfg)
	})

	return nil, nil
}

func isSupportedLogCall(call *ast.CallExpr, cfg Config) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if !isLogMethod(selector.Sel.Name) {
		return false
	}

	ident, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}

	for _, name := range cfg.AllowedLoggerNames {
		if ident.Name == name {
			return true
		}
	}

	return false
}

func isLogMethod(method string) bool {
	switch method {
	case "Debug", "Info", "Warn", "Error":
		return true
	default:
		return false
	}
}

func getStringLiteral(expr ast.Expr) (string, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return "", false
	}

	if lit.Kind != token.STRING {
		return "", false
	}

	value, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}

	return value, true
}

func checkMessage(pass *analysis.Pass, expr ast.Expr, msg string, cfg Config) {
	if cfg.CheckLowercase && !startsWithLower(msg) {
		reportWithFix(pass, expr, msg, "log message should start with lowercase letter", fixLowercase)
	}

	if cfg.CheckEnglish && !isEnglishText(msg) {
		pass.Reportf(expr.Pos(), "log message should contain only English text")
	}

	if cfg.CheckSpecialChars && !hasOnlyAllowedChars(msg) {
		reportWithFix(pass, expr, msg, "log message should not contain special characters or emoji", removeSpecialChars)
	}

	if cfg.CheckSensitiveData && hasSensitiveWords(msg, cfg.SensitivePatterns) {
		pass.Reportf(expr.Pos(), "log message should not contain sensitive data keywords")
	}
}

func getStringExpression(expr ast.Expr) (string, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return getStringLiteral(e)

	case *ast.BinaryExpr:
		if e.Op.String() != "+" {
			return "", false
		}

		left, leftOK := getStringExpression(e.X)
		right, rightOK := getStringExpression(e.Y)

		if leftOK && rightOK {
			return left + right, true
		}

		if leftOK {
			return left, true
		}

		if rightOK {
			return right, true
		}

		return "", false

	default:
		return "", false
	}
}

func reportWithFix(
	pass *analysis.Pass,
	expr ast.Expr,
	msg string,
	reportMessage string,
	fixFunc func(string) string,
) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		pass.Reportf(expr.Pos(), "%s", reportMessage)
		return
	}

	fixed := fixFunc(msg)
	if fixed == msg {
		pass.Reportf(expr.Pos(), "%s", reportMessage)
		return
	}

	pass.Report(analysis.Diagnostic{
		Pos:     expr.Pos(),
		End:     expr.End(),
		Message: reportMessage,
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "fix log message",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     lit.Pos(),
						End:     lit.End(),
						NewText: []byte(strconv.Quote(fixed)),
					},
				},
			},
		},
	})
}

func fixLowercase(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])

	return string(runes)
}

func removeSpecialChars(s string) string {
	var result []rune

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' || r == '_' {
			result = append(result, r)
		}
	}

	return string(result)
}
