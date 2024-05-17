package matcher_test

import (
	"testing"

	"github.com/frankban/quicktest"

	"mig/pkg/matcher"
)

func matchAll(line string) string {
	for _, m := range matcher.AllMatchers {
		r := m(line)
		if r != "" {
			return r
		}
	}

	return ""
}

// Test function
func TestMatchers(t *testing.T) {
	examples := []struct {
		comment quicktest.Comment
		in      string
		out     string
	}{
		//{
		//	comment: quicktest.Commentf("Simple import"),
		//	in:      "\t\"github.com/pkg/errors\"",
		//	out:     "\t\"github.com/kanisterio/errkit\"",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.New"),
		//	in:      "errors.New(\"foo\")",
		//	out:     "errkit.New(\"foo\")",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.New with fmt.Sprintf"),
		//	in:      "errors.New(fmt.Sprintf(\"foo\"))",
		//	out:     "errkit.New(fmt.Sprintf(\"foo\")) // TODO: Fixme",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.Wrap"),
		//	in:      "errors.Wrap(err, \"foo\")",
		//	out:     "errkit.Wrap(err, \"foo\")",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.Wrap with fmt.Sprintf"),
		//	in:      "errors.Wrap(err, fmt.Sprintf(\"foo\"))",
		//	out:     "errkit.Wrap(err, fmt.Sprintf(\"foo\")) // TODO: Fixme",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.Wrapf"),
		//	in:      "errors.Wrapf(err, \"foo\")",
		//	out:     "errkit.Wrap(err, \"foo\")",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.Wrapf with fmt.Sprintf"),
		//	in:      "errors.Wrapf(err, fmt.Sprintf(\"foo\"))",
		//	out:     "errkit.Wrap(err, fmt.Sprintf(\"foo\")) // TODO: Fixme",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.Wrapf with stderr 1"),
		//	in:      "errors.Wrapf(err, \"foo %s\", stderr)",
		//	out:     "errkit.Wrap(err, \"foo\", \"stderr\", stderr)",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.Wrapf with stderr 2"),
		//	in:      "errors.Wrapf(err, \"foo: %s\", stderr)",
		//	out:     "errkit.Wrap(err, \"foo\", \"stderr\", stderr)",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.Wrapf with stderr 3"),
		//	in:      "errors.Wrapf(err, \"foo:%s\", stderr)",
		//	out:     "errkit.Wrap(err, \"foo\", \"stderr\", stderr)",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.Wrapf with stderr 4"),
		//	in:      "errors.Wrapf(err, \"foo. %s\", stderr)",
		//	out:     "errkit.Wrap(err, \"foo\", \"stderr\", stderr)",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.Wrapf with stderr `5"),
		//	in:      "errors.Wrapf(err, \"foo. %s \", stderr)",
		//	out:     "errkit.Wrap(err, \"foo\", \"stderr\", stderr)",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.Wrapf with stderr 6"),
		//	in:      "errors.Wrapf(err, \"foo. Error:%s\", stderr)",
		//	out:     "errkit.Wrap(err, \"foo\", \"stderr\", stderr)",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.Wrapf with stderr 7"),
		//	in:      "errors.Wrapf(err, \"Error %s: foo.\", stderr)",
		//	out:     "errkit.Wrap(err, \"Error foo.\", \"stderr\", stderr)",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.Wrapf with stderr 8"),
		//	in:      "errors.Wrapf(err, \"Error %s foo.\", stderr)",
		//	out:     "errkit.Wrap(err, \"Error foo.\", \"stderr\", stderr)",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.Wrapf with stderr 8"),
		//	in:      "errors.Wrapf(err, \"Error %s, foo.\", stderr)",
		//	out:     "errkit.Wrap(err, \"Error foo.\", \"stderr\", stderr)",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.Wrapf with stderr and app name 1"),
		//	in:      "errors.Wrapf(err, \"Error foo: %s, app: %s\", stderr, a.name)",
		//	out:     "errkit.Wrap(err, \"Error foo\", \"stderr\", stderr, \"app\", a.name)",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.Wrapf with stderr and app name 2"),
		//	in:      "\treturn errors.Wrapf(err, \"foo. %s app=%s\", stderr, cb.name)",
		//	out:     "\treturn errkit.Wrap(err, \"foo\", \"stderr\", stderr, \"app\", cb.name)",
		//},
		{
			comment: quicktest.Commentf("Simple errors.Wrapf with stderr and stdout"),
			in:      "\treturn errors.Wrapf(err, \"Error %s, foo. stdout is %s\", stderr, stdout)",
			out:     "\treturn errkit.Wrap(err, \"Error foo.\", \"stdout\", stdout, \"stderr\", stderr)",
		},
		//{
		//	comment: quicktest.Commentf("Simple errors.Wrapf with multiple named parameters 1"),
		//	in:      "\treturn errors.Wrapf(err, \"foo. app=%s name=%s\", a.name, a.dbSubnetGroup)",
		//	out:     "\treturn errkit.Wrap(err, \"foo.\", \"app\", a.name, \"name\", a.dbSubnetGroup)",
		//},
		//{
		//	comment: quicktest.Commentf("Simple errors.Wrapf with multiple named parameters 2"),
		//	in:      "\treturn errors.Wrapf(err, \"foo. app=%s chart=%s release=%s\", c.name, c.chart.Chart, c.chart.Release)",
		//	out:     "\treturn errkit.Wrap(err, \"foo.\", \"app\", c.name, \"chart\", c.chart.Chart, \"release\", c.chart.Release)",
		//},
	}

	c := quicktest.New(t)
	for _, example := range examples {
		result := matchAll(example.in)
		c.Assert(result, quicktest.Equals, example.out, example.comment)
	}
}
