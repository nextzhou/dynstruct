package dynstruct

import (
	"fmt"
	"io"
)

var _ fmt.Formatter = Value{}

func (v Value) Format(f fmt.State, c rune) {
	if c != 'v' {
		v.formatUnknown(f, c)
		return
	}
	switch {
	case f.Flag('+'):
		v.formatPlugString(f)
	case f.Flag('#'):
		v.formatGoString(f)
	default:
		v.formatString(f)
	}
}

func (v Value) formatString(w io.Writer) {
	w.Write([]byte{'{'})
	for i, field := range v.t.fields {
		if i > 0 {
			w.Write([]byte{' '})
		}
		fv := v.value[field.name]
		fmt.Fprint(w, fv)
	}
	w.Write([]byte{'}'})
}

func (v Value) formatPlugString(w io.Writer) {
	w.Write([]byte{'{'})
	for i, field := range v.t.fields {
		if i > 0 {
			w.Write([]byte{' '})
		}
		fmt.Fprint(w, field.name)
		w.Write([]byte{':'})
		fv := v.value[field.name]
		fmt.Fprintf(w, "%+v", fv)
	}
	w.Write([]byte{'}'})
}

func (v Value) formatGoString(w io.Writer) {
	w.Write([]byte(v.t.fullName))
	w.Write([]byte{'{'})
	for i, field := range v.t.fields {
		if i > 0 {
			w.Write([]byte{',', ' '})
		}
		fmt.Fprint(w, field.name)
		w.Write([]byte{':'})
		fv := v.value[field.name]
		fmt.Fprintf(w, "%#v", fv)
	}
	w.Write([]byte{'}'})
}

func (v Value) formatUnknown(f fmt.State, c rune) {
	f.Write([]byte{'{'})
	s := "%" + string(c)
	for i, field := range v.t.fields {
		if i > 0 {
			f.Write([]byte{' '})
		}
		fv := v.value[field.name]
		if i, ok := fv.(fmt.Formatter); ok {
			i.Format(f, c)
		} else {
			fmt.Fprintf(f, s, fv)
		}
	}
	f.Write([]byte{'}'})
}
