package confparse

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type IniParser struct {
	p    *parser
	c    *config
	w    *fsnotify.Watcher
	fn   func(ev fsnotify.Event)
	name string
}

// New creates and parse a new configuration from a file name
// returns a valid parsed object and a nil, or an error and nil object
// in the successful case the values are ready to be retrieved
func New(confname string) (*IniParser, error) {
	f, err := os.Open(confname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	p := &IniParser{
		p:    newParser(f),
		c:    newConfig(),
		w:    watcher,
		name: confname,
	}
	p.Parse()

	return p, nil
}

// Parse actually parses the object content, note the object is always in a
// valid state, must be called if the Parser has been created with New, in
// case it has been created with NewFromFile it has already been parsed.
func (i *IniParser) Parse() {
	var lastsection string

	for {
		item := i.p.Scan()

		switch {
		case item.Tok == EOF:
			return
		case item.Tok == KEY_VALUE:
			i.c.C[lastsection][item.Values[0]] = item.Values[1]
		case item.Tok == SECTION:
			lastsection = item.Values[0]
			i.c.C[item.Values[0]] = make(map[string]string, 0)

		}
	}
}

// Watch add a file system watcher on the config file itself and reloads the
// configuration and parses it every time a write event is received.
func (i *IniParser) Watch() error {
	defer i.w.Close()

	dir, file := filepath.Split(i.name)
	if dir == "" {
		cdir, err := os.Getwd()
		if err != nil {
			return err
		}
		dir = cdir
	}

	done := make(chan bool)

	go i.eventFilter(file)

	if err := i.w.Add(dir); err != nil {
		return err
	}

	<-done

	return nil
}

func (i *IniParser) eventFilter(file string) {
	var err error

	for {
		select {
		case ev := <-i.w.Events:
			_, evfile := filepath.Split(ev.Name)
			if file == evfile {
				if ev.Op&fsnotify.Write == fsnotify.Write {
					i, err = New(i.name)
					if err != nil {
						log.Println("Fatal: ", err)
						return
					}
				}
				i.fn(ev)
			}
		case err := <-i.w.Errors:
			log.Println("Watcher error: ", err)
		}
	}
}

// OnConfChange accept a single parameter, a function that gets run right after
// every event and receive a copy of it.
func (i *IniParser) OnConfChange(run func(ev fsnotify.Event)) {
	i.fn = run
}

// GetBool retrieves a bool value from named section with key name, returns
// either an error and an invalid value or a nil and a valid value.
func (i *IniParser) GetBool(sectionKey string) (bool, error) {
	keys := strings.Split(sectionKey, ".")
	value, err := i.c.getValue(keys[0], keys[1], i)
	if err != nil {
		return false, err
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return false, newParserError(err.Error(), keys[0], keys[1], i.errorLine(keys[1]))
	}

	return b, nil

}

// GetInt retrieves a int64 value from named section with key name, returns
// either an error and an invalid value or a nil and a valid value.
func (i *IniParser) GetInt(sectionKey string) (int64, error) {
	keys := strings.Split(sectionKey, ".")
	value, err := i.c.getValue(keys[0], keys[1], i)
	if err != nil {
		return -1, err
	}
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return -1, newParserError(err.Error(), keys[0], keys[1], i.errorLine(keys[1]))
	}

	return n, nil

}

// GetFloat retrieves a float64 value from named section with key name, returns
// either an error and an invalid value or a nil and a valid value.
func (i *IniParser) GetFloat(sectionKey string) (float64, error) {
	keys := strings.Split(sectionKey, ".")
	value, err := i.c.getValue(keys[0], keys[1], i)
	if err != nil {
		return -0.1, err
	}
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return -1, newParserError(err.Error(), keys[0], keys[1], i.errorLine(keys[1]))

	}

	return f, nil

}

// GetString retrieves a string value from named section with key name, returns
// either an error and an invalid value or a nil and a valid value.
func (i *IniParser) GetString(sectionKey string) (string, error) {
	keys := strings.Split(sectionKey, ".")
	value, err := i.c.getValue(keys[0], keys[1], i)
	if err != nil {
		return "", err
	}
	return value, nil
}

// GetDuration retrieves a time.Duration value from the named section/key, returns either an error and an -1 or a valid duration and a nil error.
func (i *IniParser) GetDuration(sectionKey string) (time.Duration, error) {
	keys := strings.Split(sectionKey, ".")
	value, err := i.c.getValue(keys[0], keys[1], i)
	if err != nil {
		return -1, err
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return -1, err
	}

	return duration, nil
}

// GetSlice retrieves a slice value from named section with key name, returns
// either an error and an invalid value or a nil and a valid value.
func (i *IniParser) GetSlice(sectionKey string) ([]string, error) {
	keys := strings.Split(sectionKey, ".")
	value, err := i.c.getValue(keys[0], keys[1], i)
	if err != nil {
		return nil, newParserError(err.Error(), keys[0], keys[1], i.errorLine(keys[1]))
	}

	return strings.Split(value, ","), nil

}

// GetSection retrieves an entire section coverting the values to the appropriate
// type is left to the user everything is a string
func (i *IniParser) GetSection(section string) (map[string]string, error) {
	return i.c.getSection(section, i)
}

func (i *IniParser) errorLine(word string) int {
	lineno, err := i.p.s.findLine(word)
	if err == io.EOF {
		return lineno
	}
	if err != nil {
		return -1
	}
	return lineno

}

type parser struct {
	s   *lexer
	buf struct {
		tok    token
		values []string
		n      int
	}
}

func newParser(r io.Reader) *parser {
	return &parser{s: newLexer(r)}
}

func (p *parser) scan() (item *itemType) {
	// If we have a token on the buffer, then return it.
	if p.buf.values == nil {
		p.buf.values = make([]string, 0)
	}

	if p.buf.n != 0 {
		p.buf.n = 0
		item.Tok = p.buf.tok
		item.Values = append(item.Values, p.buf.values...)
	}

	// Otherwise read the next token from the scanner.
	item = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok = item.Tok
	p.buf.values = append(p.buf.values, item.Values...)

	return
}

func (p *parser) unscan() { p.buf.n = 1 }

//Scan does not take into consideration white spaces ever nor should be
// called directly.
func (p *parser) Scan() (item *itemType) {
	item = p.scan()
	if item.Tok == WHITESPACE {
		item = p.scan()
	}
	return
}

type config struct {
	C map[string]map[string]string
}

func newConfig() *config {
	conf := &config{C: make(map[string]map[string]string, 0)}
	conf.C["default"] = make(map[string]string, 0)
	conf.C["default"]["version"] = "0.1"
	return conf

}

func (c *config) getValue(section, key string, i *IniParser) (string, error) {
	sec, ok := c.C[section]
	if !ok {
		return "", newParserError(SEC_NOT_FOUND.Error(), section, key,
			i.errorLine(key))
	}
	val, ok := sec[key]
	if !ok {
		return "", newParserError(KEY_NOT_FOUND.Error(), section, key,
			i.errorLine(key))
	}

	return val, nil

}

func (c *config) getSection(section string, i *IniParser) (map[string]string, error) {
	sectionm, ok := c.C[section]
	if !ok {
		return nil, newParserError(SEC_NOT_FOUND.Error(), section, "nokey",
			i.errorLine(section))
	}

	return sectionm, nil

}
