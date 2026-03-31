package annotations

type Scanner struct {
	source string
	offset int
	line   int
	column int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
		offset: 0,
		line:   1,
		column: 1,
	}
}

// Scan returns the next token from the source.
// It returns a TokenEOF token when it reaches the end of the source.
func (s *Scanner) Scan() Token {
	s.skip()

	// End of source
	if s.offset >= len(s.source) {
		return Token{Kind: TokenEOF, Line: s.line, Column: s.column}
	}

	// Annotation: '@' followed by a letter
	if s.source[s.offset] == '@' && s.offset+1 < len(s.source) && s.isLetter(s.source[s.offset+1]) {
		return s.scanName()
	}

	return s.scanText()
}

// skip advances past whitespace and comment markers (#).
func (s *Scanner) skip() {
	for s.offset < len(s.source) {
		ch := s.source[s.offset]
		switch {
		case ch == '\n':
			s.offset++
			s.line++
			s.column = 1
		case s.isSpace(ch) || ch == '#' || ch == '\r':
			s.offset++
			s.column++
		default:
			return
		}
	}
}

// scanName reads an annotation name after '@'.
func (s *Scanner) scanName() Token {
	line := s.line
	column := s.column

	// Skip '@'
	s.offset++
	s.column++

	start := s.offset

	for s.offset < len(s.source) && s.isIdent(s.source[s.offset]) {
		s.offset++
		s.column++
	}

	return Token{
		Kind:   TokenName,
		Value:  s.source[start:s.offset],
		Line:   line,
		Column: column,
	}
}

// scanText reads free-form content until end of line.
func (s *Scanner) scanText() Token {
	line := s.line
	column := s.column
	start := s.offset

	for s.offset < len(s.source) && !s.isNewline(s.source[s.offset]) {
		s.offset++
		s.column++
	}

	// Trim
	end := s.offset
	for end > start && s.isSpace(s.source[end-1]) {
		end--
	}

	return Token{
		Kind:   TokenText,
		Value:  s.source[start:end],
		Line:   line,
		Column: column,
	}
}

func (s *Scanner) isSpace(ch byte) bool {
	return ch == ' ' || ch == '\t'
}

func (s *Scanner) isNewline(ch byte) bool {
	return ch == '\n' || ch == '\r'
}

func (s *Scanner) isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func (s *Scanner) isIdent(ch byte) bool {
	return s.isLetter(ch) || (ch >= '0' && ch <= '9') || ch == '_' || ch == '-'
}

type Token struct {
	Kind   TokenKind
	Value  string
	Line   int
	Column int
}

type TokenKind int

const (
	TokenUnknown TokenKind = iota
	TokenEOF
	TokenName // annotation name (e.g. "schema" from "@schema")
	TokenText // free text content
)
