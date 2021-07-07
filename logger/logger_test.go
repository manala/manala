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
		conf := config.New(config.WithDebug(true))
		buf := bytes.NewBufferString("")
		log := New(WithConfig(conf), WithWriter(buf))

		log.Info("info")
		log.Debug("debug")
		log.Error("error")

		s.Equal(`   • info                     
   • debug                    
   ⨯ error                    
`, buf.String())
	})

	s.Run("Debug false", func() {
		conf := config.New(config.WithDebug(false))
		buf := bytes.NewBufferString("")
		log := New(WithConfig(conf), WithWriter(buf))

		log.Info("info")
		log.Debug("debug")
		log.Error("error")

		s.Equal(`   • info                     
   ⨯ error                    
`, buf.String())
	})
}

func (s *LoggerTestSuite) TestWithDebug() {
	s.Run("True", func() {
		buf := bytes.NewBufferString("")
		log := New(WithDebug(true), WithWriter(buf))

		log.Info("info")
		log.Debug("debug")
		log.Error("error")

		s.Equal(`   • info                     
   • debug                    
   ⨯ error                    
`, buf.String())
	})

	s.Run("False", func() {
		buf := bytes.NewBufferString("")
		log := New(WithDebug(false), WithWriter(buf))

		log.Info("info")
		log.Debug("debug")
		log.Error("error")

		s.Equal(`   • info                     
   ⨯ error                    
`, buf.String())
	})
}

func (s *LoggerTestSuite) TestDebug() {
	s.Run("Message", func() {
		buf := bytes.NewBufferString("")
		log := New(WithDebug(true), WithWriter(buf))

		log.Debug("debug")

		s.Equal(`   • debug                    
`, buf.String())
	})

	s.Run("Message with field", func() {
		buf := bytes.NewBufferString("")
		log := New(WithDebug(true), WithWriter(buf))

		log.Debug("debug", log.WithField("bar", "baz"))

		s.Equal(`   • debug                     bar=baz
`, buf.String())
	})

	s.Run("Message with fields", func() {
		buf := bytes.NewBufferString("")
		log := New(WithDebug(true), WithWriter(buf))

		log.Debug("debug", log.WithField("bar", "baz"), log.WithField("qux", "quux"))

		s.Equal(`   • debug                     bar=baz qux=quux
`, buf.String())
	})
}

func (s *LoggerTestSuite) TestInfo() {
	s.Run("Message", func() {
		buf := bytes.NewBufferString("")
		log := New(WithWriter(buf))

		log.Info("info")

		s.Equal(`   • info                     
`, buf.String())
	})

	s.Run("Message with field", func() {
		buf := bytes.NewBufferString("")
		log := New(WithWriter(buf))

		log.Info("info", log.WithField("bar", "baz"))

		s.Equal(`   • info                      bar=baz
`, buf.String())
	})

	s.Run("Message with fields", func() {
		buf := bytes.NewBufferString("")
		log := New(WithWriter(buf))

		log.Info("info", log.WithField("bar", "baz"), log.WithField("qux", "quux"))

		s.Equal(`   • info                      bar=baz qux=quux
`, buf.String())
	})
}

func (s *LoggerTestSuite) TestError() {
	s.Run("Message", func() {
		buf := bytes.NewBufferString("")
		log := New(WithWriter(buf))

		log.Error("error")

		s.Equal(`   ⨯ error                    
`, buf.String())
	})

	s.Run("Message with error", func() {
		buf := bytes.NewBufferString("")
		log := New(WithWriter(buf))

		log.Error("error", errors.New("foo"))

		s.Equal(`   ⨯ error                     error=foo
`, buf.String())
	})
}
