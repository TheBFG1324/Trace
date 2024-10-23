package main

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenType represents the type of lexical tokens.
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF
	WS
	COMMENT

	// Symbols and literals
	IDENT   // identifier
	STRING  // string literal
	NUMBER  // number literal
	LBRACE  // {
	RBRACE  // }
	LPAREN  // (
	RPAREN  // )
	COMMA   // ,
	SEMICOL // ;
	EQUAL   // =
)

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Literal string
}

// Lexer represents a lexer for AICL.
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position (after current char)
	ch           byte // current char under examination
}

// NewLexer initializes a new Lexer with the input.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar advances the lexer to the next character.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for NUL, signifies EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// NextToken lexes the next token from the input.
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '{':
		tok = newToken(LBRACE, l.ch)
	case '}':
		tok = newToken(RBRACE, l.ch)
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case ';':
		tok = newToken(SEMICOL, l.ch)
	case '=':
		tok = newToken(EQUAL, l.ch)
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			tok.Type = COMMENT
			tok.Literal = l.readLineComment()
			l.readChar() // Advance to the next character after the comment
			return tok
		} else if l.peekChar() == '*' {
			l.readChar()
			tok.Type = COMMENT
			tok.Literal = l.readBlockComment()
			l.readChar() // Advance to the next character after the comment
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = IDENT
			return tok
		} else if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// Helper functions for the lexer
func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) skipWhitespace() {
	for l.ch != 0 && unicode.IsSpace(rune(l.ch)) {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	var result strings.Builder
	l.readChar() // skip opening quote

	for {
		if l.ch == '"' || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			l.readChar()
			if l.ch == '"' {
				result.WriteByte('"')
			} else {
				result.WriteByte('\\')
				result.WriteByte(l.ch)
			}
		} else {
			result.WriteByte(l.ch)
		}
		l.readChar()
		if l.ch == 0 {
			break
		}
	}

	l.readChar() // skip closing quote
	return result.String()
}

func (l *Lexer) readLineComment() string {
	position := l.position + 1
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readBlockComment() string {
	position := l.position + 1
	for {
		if l.ch == 0 {
			break
		}
		if l.ch == '*' && l.peekChar() == '/' {
			// Consume the '*' and '/'
			l.readChar()
			l.readChar()
			break
		}
		l.readChar()
	}
	return l.input[position : l.position-2]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// isLetter checks if a character is a letter.
func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}

// isDigit checks if a character is a digit.
func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}

// Data represents a global data declaration.
type Data struct {
	DataName     string
	DataType     string
	InitialValue string
}

// Permission represents permissions assigned to an agent.
type Permission struct {
	AgentName   string
	DataNames   []string
	Permissions []string
}

// Task represents a task to be executed.
type Task struct {
	TaskName   string
	AgentName  string
	Parameters map[string]string
}

// WaitStatement represents a WAIT statement.
type WaitStatement struct {
	TaskNames []string
}

// ParentRequest represents the root of the parsed script.
type ParentRequest struct {
	Statements  []interface{}          // Slice of tasks and blocks (RUNSEQ, RUNCON)
	GlobalData  map[string]*Data       // Mapping of data name to Data
	Permissions map[string]*Permission // Mapping of agent name to Permission
}

// RunSeqBlock represents a RUNSEQ block.
type RunSeqBlock struct {
	Statements []interface{} // Ordered slice of tasks and blocks
}

// RunConBlock represents a RUNCON block.
type RunConBlock struct {
	Statements map[string]interface{} // Mapping of task/block names to tasks and blocks
}

// Parser represents the parser for AICL.
type Parser struct {
	lexer     *Lexer
	curToken  Token
	peekToken Token
	errors    []string

	parentRequest *ParentRequest
}

// NewParser creates a new Parser.
func NewParser(l *Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
		parentRequest: &ParentRequest{
			Statements:  []interface{}{},
			GlobalData:  make(map[string]*Data),
			Permissions: make(map[string]*Permission),
		},
	}
	// Read two tokens to initialize curToken and peekToken
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken advances the tokens.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

// Helper functions to check for keywords
func (p *Parser) curTokenIsKeyword(keyword string) bool {
	return p.curToken.Type == IDENT && strings.ToUpper(p.curToken.Literal) == strings.ToUpper(keyword)
}

func (p *Parser) peekTokenIsKeyword(keyword string) bool {
	return p.peekToken.Type == IDENT && strings.ToUpper(p.peekToken.Literal) == strings.ToUpper(keyword)
}

func (p *Parser) expectCurTokenKeyword(keyword string) bool {
	if p.curTokenIsKeyword(keyword) {
		return true
	} else {
		p.errors = append(p.errors, fmt.Sprintf("Expected current token to be '%s', got '%s' instead", keyword, p.curToken.Literal))
		return false
	}
}

func (p *Parser) expectPeekKeyword(keyword string) bool {
	if p.peekTokenIsKeyword(keyword) {
		p.nextToken()
		return true
	} else {
		p.peekErrorKeyword(keyword)
		return false
	}
}

func (p *Parser) peekErrorKeyword(keyword string) {
	msg := fmt.Sprintf("Expected next token to be '%s', got '%s' instead", keyword, p.peekToken.Literal)
	p.errors = append(p.errors, msg)
}

// ParseProgram parses the entire program.
func (p *Parser) ParseProgram() *ParentRequest {
	for p.curToken.Type != EOF {
		if p.curToken.Type == COMMENT {
			// Skip comments
			p.nextToken()
			continue
		}
		if p.curTokenIsKeyword("START") {
			p.nextToken()
		} else if p.curTokenIsKeyword("END") {
			p.nextToken()
		} else if p.curTokenIsKeyword("DATA") {
			data := p.parseData()
			if data != nil {
				p.parentRequest.GlobalData[data.DataName] = data
			}
		} else if p.curTokenIsKeyword("PERM") {
			p.parsePermission()
		} else if p.curTokenIsKeyword("RUNSEQ") {
			seqBlock := p.parseRunSeqBlock()
			if seqBlock != nil {
				p.parentRequest.Statements = append(p.parentRequest.Statements, seqBlock)
			}
		} else if p.curTokenIsKeyword("RUNCON") {
			conBlock := p.parseRunConBlock()
			if conBlock != nil {
				p.parentRequest.Statements = append(p.parentRequest.Statements, conBlock)
			}
		} else if p.curTokenIsKeyword("TASK") {
			task := p.parseTask()
			if task != nil {
				p.parentRequest.Statements = append(p.parentRequest.Statements, task)
			}
		} else {
			p.nextToken()
		}
	}
	return p.parentRequest
}

// Parsing functions

func (p *Parser) parseData() *Data {
	data := &Data{}

	// Expect DATA_NAME
	if !p.expectPeek(IDENT) {
		return nil
	}
	data.DataName = p.curToken.Literal

	// Expect TYPE keyword
	if !p.expectPeekKeyword("TYPE") {
		return nil
	}

	// Expect DATA_TYPE
	if !p.expectPeek(IDENT) {
		return nil
	}
	data.DataType = p.curToken.Literal

	// Optionally expect VALUE keyword
	if p.peekTokenIsKeyword("VALUE") {
		p.nextToken() // move to 'VALUE' keyword
		p.nextToken() // move to value
		if p.curToken.Type != STRING && p.curToken.Type != NUMBER && p.curToken.Type != IDENT {
			p.errors = append(p.errors, "Expected value after VALUE")
			return nil
		}
		data.InitialValue = p.curToken.Literal
	}

	// Expect ';'
	if !p.expectPeek(SEMICOL) {
		return nil
	}

	p.nextToken()
	return data
}

func (p *Parser) parsePermission() {
	// Expect AGENT keyword
	if !p.expectPeekKeyword("AGENT") {
		return
	}

	// Expect AGENT_NAME
	if !p.expectPeek(IDENT) {
		return
	}
	agentName := p.curToken.Literal

	// Retrieve or create Permission object
	perm, exists := p.parentRequest.Permissions[agentName]
	if !exists {
		perm = &Permission{
			AgentName:   agentName,
			DataNames:   []string{},
			Permissions: []string{},
		}
	}

	// Expect DATA keyword
	if !p.expectPeekKeyword("DATA") {
		return
	}

	p.nextToken() // Move to first data name
	dataNames := p.parseIdentifierList()
	perm.DataNames = append(perm.DataNames, dataNames...)

	// Expect ACCESS keyword
	if !p.expectPeekKeyword("ACCESS") {
		return
	}

	p.nextToken() // Move to first permission
	permissions := p.parseIdentifierList()
	perm.Permissions = append(perm.Permissions, permissions...)

	// Expect ';'
	if !p.expectPeek(SEMICOL) {
		return
	}

	p.nextToken()
	p.parentRequest.Permissions[agentName] = perm
}

func (p *Parser) parseTask() *Task {
	task := &Task{}

	// Expect TASK_NAME
	if !p.expectPeek(IDENT) {
		return nil
	}
	task.TaskName = p.curToken.Literal

	// Expect AGENT keyword
	if !p.expectPeekKeyword("AGENT") {
		return nil
	}

	// Expect AGENT_NAME
	if !p.expectPeek(IDENT) {
		return nil
	}
	task.AgentName = p.curToken.Literal

	// Expect PARAMETERS keyword
	if !p.expectPeekKeyword("PARAMETERS") {
		return nil
	}

	// Expect '('
	if !p.expectPeek(LPAREN) {
		return nil
	}

	task.Parameters = p.parseParameters()

	// Expect ';'
	if !p.curTokenIs(SEMICOL) {
		p.errors = append(p.errors, "Expected ';' after TASK")
		return nil
	}

	p.nextToken()
	return task
}

func (p *Parser) parseRunSeqBlock() *RunSeqBlock {
	seqBlock := &RunSeqBlock{
		Statements: []interface{}{},
	}

	// Expect '{'
	if !p.expectPeek(LBRACE) {
		return nil
	}
	p.nextToken()

	for p.curToken.Type != RBRACE && p.curToken.Type != EOF {
		if p.curToken.Type == COMMENT {
			p.nextToken()
			continue
		}
		if p.curTokenIsKeyword("TASK") {
			task := p.parseTask()
			if task != nil {
				seqBlock.Statements = append(seqBlock.Statements, task)
			}
		} else if p.curTokenIsKeyword("RUNSEQ") {
			nestedSeq := p.parseRunSeqBlock()
			if nestedSeq != nil {
				seqBlock.Statements = append(seqBlock.Statements, nestedSeq)
			}
		} else if p.curTokenIsKeyword("RUNCON") {
			nestedCon := p.parseRunConBlock()
			if nestedCon != nil {
				seqBlock.Statements = append(seqBlock.Statements, nestedCon)
			}
		} else if p.curTokenIsKeyword("WAIT") {
			wait := p.parseWait()
			if wait != nil {
				seqBlock.Statements = append(seqBlock.Statements, wait)
			}
		} else {
			p.nextToken()
		}
	}
	if p.curToken.Type != RBRACE {
		p.errors = append(p.errors, "Expected '}' at the end of RUNSEQ block")
		return nil
	}
	p.nextToken()
	return seqBlock
}

func (p *Parser) parseRunConBlock() *RunConBlock {
	conBlock := &RunConBlock{
		Statements: make(map[string]interface{}),
	}

	// Expect '{'
	if !p.expectPeek(LBRACE) {
		return nil
	}
	p.nextToken()

	count := 0

	for p.curToken.Type != RBRACE && p.curToken.Type != EOF {
		if p.curToken.Type == COMMENT {
			p.nextToken()
			continue
		}
		if p.curTokenIsKeyword("TASK") {
			task := p.parseTask()
			if task != nil {
				conBlock.Statements[task.TaskName] = task
			}
		} else if p.curTokenIsKeyword("RUNSEQ") {
			nestedSeq := p.parseRunSeqBlock()
			if nestedSeq != nil {
				key := fmt.Sprintf("RUNSEQ_%d", count)
				conBlock.Statements[key] = nestedSeq
				count++
			}
		} else if p.curTokenIsKeyword("RUNCON") {
			nestedCon := p.parseRunConBlock()
			if nestedCon != nil {
				key := fmt.Sprintf("RUNCON_%d", count)
				conBlock.Statements[key] = nestedCon
				count++
			}
		} else if p.curTokenIsKeyword("WAIT") {
			wait := p.parseWait()
			if wait != nil {
				key := fmt.Sprintf("WAIT_%d", count)
				conBlock.Statements[key] = wait
				count++
			}
		} else {
			p.nextToken()
		}
	}
	if p.curToken.Type != RBRACE {
		p.errors = append(p.errors, "Expected '}' at the end of RUNCON block")
		return nil
	}
	p.nextToken()
	return conBlock
}

func (p *Parser) parseWait() *WaitStatement {
	wait := &WaitStatement{}

	p.nextToken() // Move to first task name

	wait.TaskNames = p.parseIdentifierList()

	// Expect ';'
	if !p.expectPeek(SEMICOL) {
		return nil
	}

	p.nextToken()
	return wait
}

// Helper functions

func (p *Parser) parseParameters() map[string]string {
	params := make(map[string]string)

	p.nextToken() // move to first parameter or ')'

	for p.curToken.Type != RPAREN && p.curToken.Type != EOF {
		if p.curToken.Type != IDENT {
			p.errors = append(p.errors, fmt.Sprintf("Expected parameter name, got '%s'", p.curToken.Literal))
			return nil
		}
		key := p.curToken.Literal

		if !p.expectPeek(EQUAL) {
			return nil
		}

		p.nextToken() // move to value
		if p.curToken.Type != STRING && p.curToken.Type != NUMBER && p.curToken.Type != IDENT {
			p.errors = append(p.errors, fmt.Sprintf("Expected parameter value, got '%s'", p.curToken.Literal))
			return nil
		}
		value := p.curToken.Literal
		params[key] = value

		if p.peekTokenIs(COMMA) {
			p.nextToken() // consume ','
			p.nextToken() // move to next parameter
		} else if p.peekTokenIs(RPAREN) {
			p.nextToken() // consume ')'
			break
		} else {
			p.errors = append(p.errors, "Expected ',' or ')' in parameters")
			return nil
		}
	}

	if p.curToken.Type != RPAREN {
		p.errors = append(p.errors, "Expected ')' at the end of parameters")
		return nil
	}

	p.nextToken() // Move past ')'
	return params
}

func (p *Parser) parseIdentifierList() []string {
	identifiers := []string{}

	if p.curToken.Type != IDENT {
		p.errors = append(p.errors, "Expected identifier")
		return nil
	}

	identifiers = append(identifiers, p.curToken.Literal)

	for p.peekTokenIs(COMMA) {
		p.nextToken() // consume ','
		p.nextToken() // move to next identifier
		if p.curToken.Type != IDENT {
			p.errors = append(p.errors, "Expected identifier")
			return nil
		}
		identifiers = append(identifiers, p.curToken.Literal)
	}

	return identifiers
}

// Utility functions

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekErrorToken(t)
		return false
	}
}

func (p *Parser) peekErrorToken(t TokenType) {
	msg := fmt.Sprintf("Expected next token to be %v, got %v instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func main() {
	input := `
START
    DATA data1 TYPE String VALUE "Initial Data" ;
    DATA data2 TYPE String ;
    DATA globalData TYPE string ;

    PERM AGENT Agent1 DATA data1 ACCESS READ, WRITE ;
    PERM AGENT Agent2 DATA data2 ACCESS READ ;
    PERM AGENT Agent2 DATA globalDataAdd ACCESS ADD, WRITE ;

    RUNSEQ {
        TASK FetchData AGENT Agent1 PARAMETERS (source="DB", output=data1) ;
        RUNCON {
            TASK ProcessData AGENT Agent2 PARAMETERS (input=data1, output=data2) ;
            TASK LogData AGENT Agent3 PARAMETERS (input=data1) ;
        }
        WAIT ProcessData ;
        TASK SaveData AGENT Agent4 PARAMETERS (input=data2) ;
    }
END
`

	l := NewLexer(input)
	p := NewParser(l)
	parentRequest := p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Println("Parser errors:")
		for _, e := range p.Errors() {
			fmt.Println(e)
		}
		return
	}

	// Print the parsed data structures
	fmt.Println("Global Data:")
	for name, data := range parentRequest.GlobalData {
		fmt.Printf("Name: %s, Type: %s, InitialValue: %s\n", name, data.DataType, data.InitialValue)
	}
	fmt.Println()

	fmt.Println("Permissions:")
	for agent, perm := range parentRequest.Permissions {
		fmt.Printf("Agent: %s, Data: %v, Permissions: %v\n", agent, perm.DataNames, perm.Permissions)
	}
	fmt.Println()

	fmt.Println("Parent Request Statements:")
	printStatements(parentRequest.Statements, 1)
}

// Helper function to print statements.
func printStatements(statements []interface{}, indent int) {
	prefix := strings.Repeat("    ", indent)
	for _, stmt := range statements {
		switch s := stmt.(type) {
		case *Task:
			fmt.Printf("%sTask: %s, Agent: %s, Parameters: %v\n", prefix, s.TaskName, s.AgentName, s.Parameters)
		case *RunSeqBlock:
			fmt.Printf("%sRunSeqBlock:\n", prefix)
			printStatements(s.Statements, indent+1)
		case *RunConBlock:
			fmt.Printf("%sRunConBlock:\n", prefix)
			for key, conStmt := range s.Statements {
				fmt.Printf("%s    Key: %s\n", prefix, key)
				printStatements([]interface{}{conStmt}, indent+2)
			}
		case *WaitStatement:
			fmt.Printf("%sWait for tasks: %v\n", prefix, s.TaskNames)
		default:
			fmt.Printf("%sUnknown statement type\n", prefix)
		}
	}
}
