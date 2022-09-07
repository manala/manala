package report

func NewReport(message string) *Report {
	return &Report{
		message: message,
		fields:  map[string]interface{}{},
		reports: []*Report{},
	}
}

func NewErrorReport(err error) *Report {
	report := &Report{
		err:     err,
		fields:  map[string]interface{}{},
		reports: []*Report{},
	}

	for err != nil {
		if _err, ok := err.(Reporter); ok {
			_err.Report(report)
		}
		if _err, ok := err.(interface{ Unwrap() error }); ok {
			err = _err.Unwrap()
		} else {
			err = nil
		}
	}

	return report
}

type Reporter interface {
	Report(*Report)
}

type Report struct {
	message string
	fields  map[string]interface{}
	err     error
	trace   string
	reports []*Report
}

func (report *Report) String() string {
	if report.message != "" {
		return report.message
	}

	if report.err != nil {
		return report.err.Error()
	}

	return ""
}

func (report *Report) Compose(options ...Option) {
	for _, option := range options {
		option(report)
	}
}

func (report *Report) Message() string {
	return report.message
}

func (report *Report) Fields() map[string]interface{} {
	return report.fields
}

func (report *Report) Err() error {
	return report.err
}

func (report *Report) Trace() string {
	return report.trace
}

func (report *Report) Reports() []*Report {
	return report.reports
}

func (report *Report) Add(rep *Report) {
	report.reports = append(report.reports, rep)
}

type Option func(report *Report)

func WithMessage(message string) Option {
	return func(report *Report) {
		report.message = message
	}
}

func WithField(key string, value interface{}) Option {
	return func(report *Report) {
		report.fields[key] = value
	}
}

func WithErr(err error) Option {
	return func(report *Report) {
		report.err = err
	}
}

func WithTrace(trace string) Option {
	return func(report *Report) {
		report.trace = trace
	}
}
