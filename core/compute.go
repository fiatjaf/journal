package core

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
	"github.com/tidwall/buntdb"
)

var computeModule *gojq.Module
var methodsAvailable = make(map[string][]string)
var methodsCallable = make(map[string]*gojq.Code)

type moduleLoader struct{}

func (*moduleLoader) LoadModule(_ string) (*gojq.Module, error) {
	return computeModule, nil
}

func prepareComputation() {
	// read computation file as a module
	cfile, err := ioutil.ReadFile(COMPUTE_FILE)
	if err != nil {
		log.Fatal().Err(err).Str("path", COMPUTE_FILE).
			Msg("couldn't open compuation file")
	}
	computeModule, err = gojq.ParseModule(string(cfile))
	if err != nil {
		log.Fatal().Err(err).Str("path", COMPUTE_FILE).
			Msg("jq module parsing error")
	}

	loader := &moduleLoader{}

	for _, funcdef := range computeModule.FuncDefs {
		methodsAvailable[funcdef.Name] = funcdef.Args

		// compile gojq code for each method
		var argsStr string
		vars := make([]string, len(funcdef.Args))
		if len(funcdef.Args) == 0 {
			argsStr = ""
		} else {
			for i, _ := range funcdef.Args {
				vars[i] = "$_var" + strconv.Itoa(i)
			}
			argsStr = "(" + strings.Join(vars, ";") + ")"
		}

		p, _ := gojq.Parse(`import "compute" as m; m::` + funcdef.Name + argsStr)
		code, err := gojq.Compile(p,
			gojq.WithModuleLoader(loader),
			gojq.WithVariables(vars),
		)
		if err != nil {
			log.Fatal().Err(err).Str("func", funcdef.Name).Msg("failed to compile")
		}
		methodsCallable[funcdef.Name] = code
	}

	delete(methodsAvailable, "init")
}

func Compute(
	state interface{},
	method string,
	params []interface{},
) (interface{}, error) {
	log := log.With().
		Interface("state", state).Str("method", method).Interface("params", params).
		Logger()

	code, ok := methodsCallable[method]
	if ok {
		state, _ = code.Run(state, params...).Next()
		if err, isErr := state.(error); isErr {
			log.Warn().Err(err).Msg("compute error")
			return nil, err
		}
	}

	log.Debug().Interface("after", state).Msg("compute success")
	return state, nil
}

func ComputeAll() (state interface{}, err error) {
	state, err = Compute(make(map[string]interface{}), "init", []interface{}{})
	if err != nil {
		return nil, err
	}

	err = db.View(func(tx *buntdb.Tx) error {
		var err error

		tx.Ascend("datepos", func(_, v string) bool {
			var value LogEntry
			err = json.Unmarshal([]byte(v), &value)
			if err != nil {
				return err
			}

			method, ok := methodsAvailable[value.Method]
			if !ok {
				return true
			}

			value.Params["date"] = value.Time
			params := make([]interface{}, len(method))

			for i, argName := range methodsAvailable[value.Method] {
				v, _ := value.Params[argName]
				params[i] = v
			}

			state, err = Compute(state, value.Method, params)
			return true
		})

		return err
	})

	return state, err
}
