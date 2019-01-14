package jsonscan

import "fmt"

type Classes = int

const _E Classes = -1

const (
	C_SPACE Classes = iota // space
	C_WHITE                // OTHER WHITESPACE
	C_LCURB                // {
	C_RCURB                // }
	C_LSQRB                // [
	C_RSQRB                // ]
	C_COLON                // :
	C_COMMA                // ,
	C_QUOTE                // "
	C_BACKS                // \
	C_SLASH                // /
	C_PLUS                 // +
	C_MINUS                // -
	C_POINT                // .
	C_ZERO                 // 0
	C_DIGIT                // 123456789
	C_LOW_A                // A
	C_LOW_B                // B
	C_LOW_C                // C
	C_LOW_D                // D
	C_LOW_E                // E
	C_LOW_F                // F
	C_LOW_L                // L
	C_LOW_N                // N
	C_LOW_R                // R
	C_LOW_S                // S
	C_LOW_T                // T
	C_LOW_U                // U
	C_ABCDF                // ABCDF
	C_E                    // E
	C_ETC                  // EVERYTHING ELSE
	NR_CLASSES
)

var AsciiClass = []Classes{
	/*
	   This array maps the 128 ASCII characters into character classes.
	   The remaining Unicode characters should be mapped to C_ETC.
	   Non-whitespace control characters are errors.
	*/
	_E, _E, _E, _E, _E, _E, _E, _E,
	_E, C_WHITE, C_WHITE, _E, _E, C_WHITE, _E, _E,
	_E, _E, _E, _E, _E, _E, _E, _E,
	_E, _E, _E, _E, _E, _E, _E, _E,

	C_SPACE, C_ETC, C_QUOTE, C_ETC, C_ETC, C_ETC, C_ETC, C_ETC,
	C_ETC, C_ETC, C_ETC, C_PLUS, C_COMMA, C_MINUS, C_POINT, C_SLASH,
	C_ZERO, C_DIGIT, C_DIGIT, C_DIGIT, C_DIGIT, C_DIGIT, C_DIGIT, C_DIGIT,
	C_DIGIT, C_DIGIT, C_COLON, C_ETC, C_ETC, C_ETC, C_ETC, C_ETC,

	C_ETC, C_ABCDF, C_ABCDF, C_ABCDF, C_ABCDF, C_E, C_ABCDF, C_ETC,
	C_ETC, C_ETC, C_ETC, C_ETC, C_ETC, C_ETC, C_ETC, C_ETC,
	C_ETC, C_ETC, C_ETC, C_ETC, C_ETC, C_ETC, C_ETC, C_ETC,
	C_ETC, C_ETC, C_ETC, C_LSQRB, C_BACKS, C_RSQRB, C_ETC, C_ETC,

	C_ETC, C_LOW_A, C_LOW_B, C_LOW_C, C_LOW_D, C_LOW_E, C_LOW_F, C_ETC,
	C_ETC, C_ETC, C_ETC, C_ETC, C_LOW_L, C_ETC, C_LOW_N, C_ETC,
	C_ETC, C_ETC, C_LOW_R, C_LOW_S, C_LOW_T, C_LOW_U, C_ETC, C_ETC,
	C_ETC, C_ETC, C_ETC, C_LCURB, C_ETC, C_RCURB, C_ETC, C_ETC,
}

type state = int

const (
	GO state = iota // start
	OK              // ok
	OB              // object
	KE              // key
	CO              // colon
	VA              // value
	AR              // array
	ST              // string
	ES              // escape
	U1              // u1
	U2              // u2
	U3              // u3
	U4              // u4
	MI              // minus
	ZE              // zero
	IN              // integer
	FR              // fraction
	FS              // fraction
	E1              // e
	E2              // ex
	E3              // exp
	T1              // tr
	T2              // tru
	T3              // true
	F1              // fa
	F2              // fal
	F3              // fals
	F4              // false
	N1              // nu
	N2              // nul
	N3              // null
	NR_STATES
)

var sn = []string{
	"GO",
	"OK",
	"OB",
	"KE",
	"CO",
	"VA",
	"AR",
	"ST",
	"ES",
	"U1",
	"U2",
	"U3",
	"U4",
	"MI",
	"ZE",
	"IN",
	"FR",
	"FS",
	"E1",
	"E2",
	"E3",
	"T1",
	"T2",
	"T3",
	"F1",
	"F2",
	"F3",
	"F4",
	"N1",
	"N2",
	"N3",
}

var StateTransitionTable = [NR_STATES][NR_CLASSES]state{
	/*
	   The state transition table takes the current state and the current symbol,
	   and returns either a new state or an action. An action is represented as a
	   negative number. A JSON text is accepted if at the end of the text the
	   state is OK and if the mode is MODE_DONE.
	                white                                      1-9                                   ABCDF  etc
	            space |  {  }  [  ]  :  ,  "  \  /  +  -  .  0  |  a  b  c  d  e  f  l  n  r  s  t  u  |  E  |*/
	/*start  GO*/ {GO, GO, -6, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E},
	/*ok     OK*/ {OK, OK, _E, -8, _E, -7, _E, -3, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E},
	/*object OB*/ {OB, OB, _E, -9, _E, _E, _E, _E, ST, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E},
	/*key    KE*/ {KE, KE, _E, _E, _E, _E, _E, _E, ST, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E},
	/*colon  CO*/ {CO, CO, _E, _E, _E, _E, -2, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E},
	/*value  VA*/ {VA, VA, -6, _E, -5, _E, _E, _E, ST, _E, _E, _E, MI, _E, ZE, IN, _E, _E, _E, _E, _E, F1, _E, N1, _E, _E, T1, _E, _E, _E, _E},
	/*array  AR*/ {AR, AR, -6, _E, -5, -7, _E, _E, ST, _E, _E, _E, MI, _E, ZE, IN, _E, _E, _E, _E, _E, F1, _E, N1, _E, _E, T1, _E, _E, _E, _E},
	/*string ST*/ {ST, _E, ST, ST, ST, ST, ST, ST, -4, ES, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST},
	/*escape ES*/ {_E, _E, _E, _E, _E, _E, _E, _E, ST, ST, ST, _E, _E, _E, _E, _E, _E, ST, _E, _E, _E, ST, _E, ST, ST, _E, ST, U1, _E, _E, _E},
	/*u1     U1*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, U2, U2, U2, U2, U2, U2, U2, U2, _E, _E, _E, _E, _E, _E, U2, U2, _E},
	/*u2     U2*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, U3, U3, U3, U3, U3, U3, U3, U3, _E, _E, _E, _E, _E, _E, U3, U3, _E},
	/*u3     U3*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, U4, U4, U4, U4, U4, U4, U4, U4, _E, _E, _E, _E, _E, _E, U4, U4, _E},
	/*u4     U4*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, ST, ST, ST, ST, ST, ST, ST, ST, _E, _E, _E, _E, _E, _E, ST, ST, _E},
	/*minus  MI*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, ZE, IN, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E},
	/*zero   ZE*/ {OK, OK, _E, -8, _E, -7, _E, -3, _E, _E, _E, _E, _E, FR, _E, _E, _E, _E, _E, _E, E1, _E, _E, _E, _E, _E, _E, _E, _E, E1, _E},
	/*int    IN*/ {OK, OK, _E, -8, _E, -7, _E, -3, _E, _E, _E, _E, _E, FR, IN, IN, _E, _E, _E, _E, E1, _E, _E, _E, _E, _E, _E, _E, _E, E1, _E},
	/*frac   FR*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, FS, FS, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E},
	/*fracs  FS*/ {OK, OK, _E, -8, _E, -7, _E, -3, _E, _E, _E, _E, _E, _E, FS, FS, _E, _E, _E, _E, E1, _E, _E, _E, _E, _E, _E, _E, _E, E1, _E},
	/*e      E1*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, E2, E2, _E, E3, E3, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E},
	/*ex     E2*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, E3, E3, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E},
	/*exp    E3*/ {OK, OK, _E, -8, _E, -7, _E, -3, _E, _E, _E, _E, _E, _E, E3, E3, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E},
	/*tr     T1*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, T2, _E, _E, _E, _E, _E, _E},
	/*tru    T2*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, T3, _E, _E, _E},
	/*true   T3*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, OK, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E},
	/*fa     F1*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, F2, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E},
	/*fal    F2*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, F3, _E, _E, _E, _E, _E, _E, _E, _E},
	/*fals   F3*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, F4, _E, _E, _E, _E, _E},
	/*false  F4*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, OK, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E},
	/*nu     N1*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, N2, _E, _E, _E},
	/*nul    N2*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, N3, _E, _E, _E, _E, _E, _E, _E, _E},
	/*null   N3*/ {_E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, _E, OK, _E, _E, _E, _E, _E, _E, _E, _E},
}

type mode uint8

const (
	MODE_ARRAY mode = iota
	MODE_DONE
	MODE_KEY
	MODE_OBJECT
)

var ms = []string{
	"ARRAY",
	"DONE",
	"KEY",
	"OBJECT",
}

type scanner struct {
	pos   int
	state state
	data  []byte
	modes modeStack
}

func newScanner(data []byte) *scanner {
	s := &scanner{
		state: GO,
		data:  data,
	}
	s.modes.push(MODE_DONE)
	return s
}

type modeStack []mode

func (s *modeStack) push(m mode) {
	*s = append(*s, m)
}

func (s *modeStack) pop(m mode) bool {
	l := len(*s)
	if l == 0 {
		return false
	}
	if (*s)[l-1] != m {
		return false
	}
	*s = (*s)[:len(*s)-1]
	return true
}

func (s modeStack) top() mode {
	l := len(s)
	if l == 0 {
		return mode(0xff)
	}
	return s[l-1]
}

func (s modeStack) len() int {
	return len(s)
}

type KV struct {
	Key   string
	Value []byte
}

var je = fmt.Errorf("invalid json")

func (s *scanner) scan() ([]KV, error) {
	var keyBeg, keyEnd, valBeg, valEnd int
	var nextClass Classes
	var kvs []KV
	for idx, c := range s.data {
		if c >= 128 {
			nextClass = C_ETC
		} else {
			nextClass = AsciiClass[int(c)]
			if nextClass == _E {
				return nil, je
			}
		}

		nextState := StateTransitionTable[s.state][nextClass]

		if s.modes.len() == 2 {
			switch s.state {
			case OB, KE:
				if nextState == ST {
					keyBeg = idx
				}
			case VA:
				if nextState != VA {
					valBeg = idx
				}
			case ZE, IN, FS, E3:
				if nextState == CO || nextState == OK {
					valEnd = idx
					kvs = append(kvs, s.getKV(keyBeg, keyEnd, valBeg, valEnd))
				}
			case T3, F4, N3:
				if nextState == OK {
					valEnd = idx + 1
					kvs = append(kvs, s.getKV(keyBeg, keyEnd, valBeg, valEnd))
				}
			}
		}

		if nextState >= 0 {
			s.state = nextState
		} else {
			switch nextState {
			case -9:
				if !s.modes.pop(MODE_KEY) {
					return nil, je
				}
				s.state = OK
			case -8:
				if !s.modes.pop(MODE_OBJECT) {
					return nil, je
				}
				if s.modes.len() == 2 {
					valEnd = idx + 1
					kvs = append(kvs, s.getKV(keyBeg, keyEnd, valBeg, valEnd))
				}
				s.state = OK
			case -7:
				if !s.modes.pop(MODE_ARRAY) {
					return nil, je
				}
				if s.modes.len() == 2 {
					valEnd = idx + 1
					kvs = append(kvs, s.getKV(keyBeg, keyEnd, valBeg, valEnd))
				}
				s.state = OK
			case -6:
				s.modes.push(MODE_KEY)
				s.state = OB
			case -5:
				s.modes.push(MODE_ARRAY)
				s.state = AR
			case -4:
				switch s.modes.top() {
				case MODE_KEY:
					s.state = CO
					if s.modes.len() == 2 {
						keyEnd = idx + 1
					}
				case MODE_ARRAY, MODE_OBJECT:
					if s.modes.len() == 2 {
						valEnd = idx + 1
						kvs = append(kvs, s.getKV(keyBeg, keyEnd, valBeg, valEnd))
					}
					s.state = OK
				default:
					return nil, je
				}
			case -3:
				switch s.modes.top() {
				case MODE_OBJECT:
					if !s.modes.pop(MODE_OBJECT) {
						return nil, je
					}
					s.modes.push(MODE_KEY)
					if s.modes.len() == 2 && (s.state == ZE || s.state == IN || s.state == FS || s.state == E3) {
						valEnd = idx
						kvs = append(kvs, s.getKV(keyBeg, keyEnd, valBeg, valEnd))
					}
					s.state = KE
				case MODE_ARRAY:
					s.state = VA
				default:
					return nil, je
				}
			case -2:
				if !s.modes.pop(MODE_KEY) {
					return nil, je
				}
				s.modes.push(MODE_OBJECT)
				s.state = VA
			default:
				return nil, je
			}
		}
	}
	ok := s.state == OK && s.modes.pop(MODE_DONE)
	if !ok {
		return nil, je
	}
	return kvs, nil
}

func (s *scanner) getKV(keyBeg, keyEnd, valBeg, valEnd int) KV {
	return KV{
		Key:   string(s.data[keyBeg+1 : keyEnd-1]),
		Value: s.data[valBeg:valEnd],
	}
}

func Scan(data []byte) ([]KV, error) {
	return newScanner(data).scan()
}
