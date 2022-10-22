package argparser

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type ArgParser struct {
	args []string
	i    int
}

func New(args []string) *ArgParser {
	return &ArgParser{
		args: args,
		i:    0,
	}
}

func (a *ArgParser) NextArg() string {
	if a.i >= len(a.args) {
		return ""
	}
	a.i++
	return a.args[a.i-1]
}

func (a *ArgParser) RemainingArgs() []string {
	i := a.i
	a.i = len(a.args)
	return a.args[i:]
}

func (a *ArgParser) NextOptions(defaults interface{}) error {
	if reflect.ValueOf(defaults).Kind() != reflect.Ptr {
		return fmt.Errorf("defaults struct must be pointer to struct")
	}
	lastMultipart := false
	lastMultipartTag := ""

	for i := a.i; i < len(a.args); i++ {
		if lastMultipart {
			// The last option required an argument to follow it
			lastMultipart = false
			if err := setTaggedValue(defaults, lastMultipartTag, a.args[i]); err != nil {
				return err
			}
		} else if !strings.HasPrefix(a.args[i], "-") {
			// This is not an option and the last option was not multipart
			a.i = i
			return nil
		} else {
			// This is an option. Is it single character or multi character
			if strings.HasPrefix(a.args[i], "--") {
				// Multi letter option
				sArg := strings.TrimPrefix(a.args[i], "--")
				isMultipart, err := isTaggedValueMultipart(defaults, sArg)
				if err != nil {
					return err
				}
				if isMultipart {
					lastMultipart = true
					lastMultipartTag = sArg
				} else {
					if err := setTaggedValue(defaults, sArg, "true"); err != nil {
						return err
					}
				}
			} else {
				// Single letter option(s)
				sArgs := toStrings(strings.TrimPrefix(a.args[i], "-"))
				if len(sArgs) > 1 {
					// Multiple single letter options (only booleans flags are allowed)
					for _, sArg := range sArgs {
						isMultipart, err := isTaggedValueMultipart(defaults, sArg)
						if err != nil {
							return err
						}
						if isMultipart {
							return fmt.Errorf("cannot use argument %s in that form as it requires an option", sArg)
						} else {
							if err := setTaggedValue(defaults, sArg, "true"); err != nil {
								return err
							}
						}
					}
				} else {
					// One single letter option
					sArg := sArgs[0]
					isMultipart, err := isTaggedValueMultipart(defaults, sArg)
					if err != nil {
						return err
					}
					if isMultipart {
						lastMultipart = true
						lastMultipartTag = sArg
					} else {
						if err := setTaggedValue(defaults, sArg, "true"); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	if lastMultipart {
		return fmt.Errorf("option %s was not supplied an argument", lastMultipartTag)
	}
	a.i = len(a.args)
	return nil
}

func toStrings(s string) []string {
	ss := make([]string, len(s))
	for si, sv := range s {
		ss[si] = string(sv)
	}
	return ss
}

func setTaggedValue(obj interface{}, tag string, v string) error {
	f, err := getTaggedReflectField(obj, tag)
	if err != nil {
		return err
	}
	wrongTypeMsg := fmt.Sprintf("could not convert arg '%s' to %s", v, f.Type())
	var vp interface{}
	var vpErr error
	switch f.Type() {
	// Special case: string does not need to be surrounded by ""
	case reflect.TypeOf("string"):
		f.Set(reflect.ValueOf(v))
		return nil

	// Special case: duration should be parsed with duration module to allow stuff like 5h20m1s
	case reflect.TypeOf(time.Duration(0)):
		if vp, err := time.ParseDuration(v); err != nil {
			return fmt.Errorf(wrongTypeMsg)
		} else {
			f.Set(reflect.ValueOf(vp))
			return nil
		}

	// Parse everything else with json parsing
	case reflect.TypeOf(int(0)):
		vp, vpErr = parseAndCast[int](v)
	case reflect.TypeOf(int8(0)):
		vp, vpErr = parseAndCast[int8](v)
	case reflect.TypeOf(int16(0)):
		vp, vpErr = parseAndCast[int16](v)
	case reflect.TypeOf(int32(0)):
		vp, vpErr = parseAndCast[int32](v)
	case reflect.TypeOf(int64(0)):
		vp, vpErr = parseAndCast[int64](v)

	case reflect.TypeOf(uint(0)):
		vp, vpErr = parseAndCast[uint](v)
	case reflect.TypeOf(uint8(0)):
		vp, vpErr = parseAndCast[uint8](v)
	case reflect.TypeOf(uint16(0)):
		vp, vpErr = parseAndCast[uint16](v)
	case reflect.TypeOf(uint32(0)):
		vp, vpErr = parseAndCast[uint32](v)
	case reflect.TypeOf(uint64(0)):
		vp, vpErr = parseAndCast[uint64](v)

	case reflect.TypeOf(false):
		vp, vpErr = parseAndCast[bool](v)

	case reflect.TypeOf(float32(0)):
		vp, vpErr = parseAndCast[float32](v)
	case reflect.TypeOf(float64(0)):
		vp, vpErr = parseAndCast[float64](v)

	default:
		return fmt.Errorf("parsing args of type %s is not supported yet (for arg '%s')", f.Type(), tag)
	}
	if vpErr != nil {
		return fmt.Errorf(wrongTypeMsg)
	} else {
		f.Set(reflect.ValueOf(vp))
	}
	return nil
}

func parseAndCast[T any](x string) (T, error) {
	var y T
	if err := json.Unmarshal([]byte(x), &y); err != nil {
		return y, fmt.Errorf("could not convert '%s' to %s", x, reflect.TypeOf(y))
	}
	return y, nil
}

func isTaggedValueMultipart(obj interface{}, tag string) (bool, error) {
	f, err := getTaggedReflectField(obj, tag)
	if err != nil {
		return false, err
	}
	return f.Type() != reflect.TypeOf(false), nil
}

func getTaggedReflectField(obj interface{}, tag string) (reflect.Value, error) {
	rVal := reflect.ValueOf(obj).Elem()
	rType := reflect.TypeOf(obj).Elem()
	for i := 0; i < rType.NumField(); i++ {
		fieldTags := strings.Split(rType.FieldByIndex([]int{i}).Tag.Get("flag"), "|")
		for _, f := range fieldTags {
			if f == tag {
				return rVal.FieldByIndex([]int{i}), nil
			}
		}
	}
	return reflect.ValueOf(nil), fmt.Errorf("cannot find flag %s", tag)
}
