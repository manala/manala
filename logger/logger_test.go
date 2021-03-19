package logger

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/suite"
	"manala/config"
	"testing"
)

/*********/
/* Suite */
/*********/

type LoggerTestSuite struct {
	suite.Suite
}

func TestLoggerTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(LoggerTestSuite))
}

/*********/
/* Tests */
/*********/

func (s *LoggerTestSuite) TestWithConfig() {
	s.Run("Debug true", func() {
		conf := config.New("", "")
		conf.SetDebug(true)
		out := bytes.NewBufferString("")
		log := New(WithConfig(conf), WithWriter(out))

		log.Info("info")
		log.Debug("debug")
		log.Error("error")

		s.Equal(`   • info                     
   • debug                    
   ⨯ error                    
`, out.String())
	})

	s.Run("Debug false", func() {
		conf := config.New("", "")
		conf.SetDebug(false)
		out := bytes.NewBufferString("")
		log := New(WithConfig(conf), WithWriter(out))

		log.Info("info")
		log.Debug("debug")
		log.Error("error")

		s.Equal(`   • info                     
   ⨯ error                    
`, out.String())
	})
}

func (s *LoggerTestSuite) TestWithDebug() {
	s.Run("True", func() {
		out := bytes.NewBufferString("")
		log := New(WithDebug(true), WithWriter(out))

		log.Info("info")
		log.Debug("debug")
		log.Error("error")

		s.Equal(`   • info                     
   • debug                    
   ⨯ error                    
`, out.String())
	})

	s.Run("False", func() {
		out := bytes.NewBufferString("")
		log := New(WithDebug(false), WithWriter(out))

		log.Info("info")
		log.Debug("debug")
		log.Error("error")

		s.Equal(`   • info                     
   ⨯ error                    
`, out.String())
	})
}

func (s *LoggerTestSuite) TestDebug() {
	s.Run("Message", func() {
		out := bytes.NewBufferString("")
		log := New(WithDebug(true), WithWriter(out))

		log.Debug("debug")

		s.Equal(`   • debug                    
`, out.String())
	})

	s.Run("Message with field", func() {
		out := bytes.NewBufferString("")
		log := New(WithDebug(true), WithWriter(out))

		log.Debug("debug", log.WithField("bar", "baz"))

		s.Equal(`   • debug                     bar=baz
`, out.String())
	})

	s.Run("Message with fields", func() {
		out := bytes.NewBufferString("")
		log := New(WithDebug(true), WithWriter(out))

		log.Debug("debug", log.WithField("bar", "baz"), log.WithField("qux", "quux"))

		s.Equal(`   • debug                     bar=baz qux=quux
`, out.String())
	})
}

func (s *LoggerTestSuite) TestInfo() {
	s.Run("Message", func() {
		out := bytes.NewBufferString("")
		log := New(WithWriter(out))

		log.Info("info")

		s.Equal(`   • info                     
`, out.String())
	})

	s.Run("Message with field", func() {
		out := bytes.NewBufferString("")
		log := New(WithWriter(out))

		log.Info("info", log.WithField("bar", "baz"))

		s.Equal(`   • info                      bar=baz
`, out.String())
	})

	s.Run("Message with fields", func() {
		out := bytes.NewBufferString("")
		log := New(WithWriter(out))

		log.Info("info", log.WithField("bar", "baz"), log.WithField("qux", "quux"))

		s.Equal(`   • info                      bar=baz qux=quux
`, out.String())
	})
}

func (s *LoggerTestSuite) TestError() {
	s.Run("Message", func() {
		out := bytes.NewBufferString("")
		log := New(WithWriter(out))

		log.Error("error")

		s.Equal(`   ⨯ error                    
`, out.String())
	})

	s.Run("Message with error", func() {
		out := bytes.NewBufferString("")
		log := New(WithWriter(out))

		log.Error("error", errors.New("foo"))

		s.Equal(`   ⨯ error                     error=foo
`, out.String())
	})
}
