package filter

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenIdent
	TokenString
	TokenNumber
	TokenEQ
	TokenNE
	TokenGT
	TokenGE
	TokenLT
	TokenLE
	TokenAND
	TokenOR
	TokenLParen
	TokenRParen
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input string
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return rune(l.input[l.pos])
}

func (l *Lexer) advance() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	ch := rune(l.input[l.pos])
	l.pos++
	return ch
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && unicode.IsSpace(rune(l.input[l.pos])) {
		l.pos++
	}
}

func (l *Lexer) readString() string {
	var sb strings.Builder
	for {
		ch := l.peek()
		if ch == 0 || ch == '\'' || ch == '"' {
			break
		}
		sb.WriteRune(l.advance())
	}
	return sb.String()
}

func (l *Lexer) readNumber() string {
	var sb strings.Builder
	for {
		ch := l.peek()
		if ch == 0 || (!unicode.IsDigit(ch) && ch != '.' && ch != '-') {
			break
		}
		sb.WriteRune(l.advance())
	}
	return sb.String()
}

func (l *Lexer) readIdent() string {
	var sb strings.Builder
	for {
		ch := l.peek()
		if ch == 0 || (!unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_') {
			break
		}
		sb.WriteRune(l.advance())
	}
	return sb.String()
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	ch := l.peek()
	if ch == 0 {
		return Token{Type: TokenEOF}
	}

	switch ch {
	case '(':
		l.advance()
		return Token{Type: TokenLParen, Value: "("}
	case ')':
		l.advance()
		return Token{Type: TokenRParen, Value: ")"}
	case '=':
		l.advance()
		return Token{Type: TokenEQ, Value: "="}
	case '!':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			return Token{Type: TokenNE, Value: "!="}
		}
		return Token{Type: TokenNE, Value: "!"}
	case '>':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			return Token{Type: TokenGE, Value: ">="}
		}
		return Token{Type: TokenGT, Value: ">"}
	case '<':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			return Token{Type: TokenLE, Value: "<="}
		}
		return Token{Type: TokenLT, Value: "<"}
	case '\'', '"':
		l.advance()
		val := l.readString()
		if l.peek() == '\'' || l.peek() == '"' {
			l.advance()
		}
		return Token{Type: TokenString, Value: val}
	}

	if unicode.IsDigit(ch) || (ch == '-' && len(l.input) > l.pos+1 && unicode.IsDigit(rune(l.input[l.pos+1]))) {
		return Token{Type: TokenNumber, Value: l.readNumber()}
	}

	if unicode.IsLetter(ch) || ch == '_' {
		ident := l.readIdent()
		switch strings.ToUpper(ident) {
		case "AND":
			return Token{Type: TokenAND, Value: ident}
		case "OR":
			return Token{Type: TokenOR, Value: ident}
		}
		return Token{Type: TokenIdent, Value: ident}
	}

	l.advance()
	return Token{Type: TokenEOF}
}

func (l *Lexer) Tokenize() []Token {
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}
	return tokens
}

type Expr interface {
	exprNode()
}

type ComparisonExpr struct {
	Field    string
	Operator TokenType
	Value    interface{}
}

func (e *ComparisonExpr) exprNode() {}

type LogicalExpr struct {
	Left     Expr
	Operator TokenType
	Right    Expr
}

func (e *LogicalExpr) exprNode() {}

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) current() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() Token {
	tok := p.current()
	p.pos++
	return tok
}

func (p *Parser) Parse() (Expr, error) {
	if p.current().Type == TokenEOF {
		return nil, nil
	}
	return p.parseOr()
}

func (p *Parser) parseOr() (Expr, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	for p.current().Type == TokenOR {
		p.advance()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &LogicalExpr{Left: left, Operator: TokenOR, Right: right}
	}

	return left, nil
}

func (p *Parser) parseAnd() (Expr, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for p.current().Type == TokenAND {
		p.advance()
		right, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		left = &LogicalExpr{Left: left, Operator: TokenAND, Right: right}
	}

	return left, nil
}

func (p *Parser) parsePrimary() (Expr, error) {
	if p.current().Type == TokenLParen {
		p.advance()
		expr, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		if p.current().Type != TokenRParen {
			return nil, fmt.Errorf("expected closing parenthesis")
		}
		p.advance()
		return expr, nil
	}

	return p.parseComparison()
}

func (p *Parser) parseComparison() (Expr, error) {
	if p.current().Type != TokenIdent {
		return nil, fmt.Errorf("expected field name, got %v", p.current())
	}

	field := p.advance().Value

	opTypes := []TokenType{TokenEQ, TokenNE, TokenGT, TokenGE, TokenLT, TokenLE}
	isOp := false
	for _, t := range opTypes {
		if p.current().Type == t {
			isOp = true
			break
		}
	}

	if !isOp {
		return nil, fmt.Errorf("expected comparison operator, got %v", p.current())
	}

	op := p.advance()

	var value interface{}
	switch p.current().Type {
	case TokenString:
		value = p.advance().Value
	case TokenNumber:
		numStr := p.advance().Value
		if strings.Contains(numStr, ".") {
			var f float64
			_, err := fmt.Sscanf(numStr, "%f", &f)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", numStr)
			}
			value = f
		} else {
			var n int
			_, err := fmt.Sscanf(numStr, "%d", &n)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", numStr)
			}
			value = n
		}
	case TokenIdent:
		value = p.advance().Value
	default:
		return nil, fmt.Errorf("expected value, got %v", p.current())
	}

	return &ComparisonExpr{Field: field, Operator: op.Type, Value: value}, nil
}

func Parse(input string) (Expr, error) {
	lexer := NewLexer(input)
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)
	return parser.Parse()
}
