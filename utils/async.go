package utils

import (
	"log"
	"reflect"
	"sync"
)

// Result struct to store the outcome of each async function call
type Result struct {
	Output []interface{}
	Error  error
}

// AsyncExecute executes a target function asynchronously for each set of inputs
func AsyncExecute(targetFunc interface{}, inputs [][]interface{}) []Result {
	// Validate that the target function is indeed a function
	funcValue := reflect.ValueOf(targetFunc)
	if funcValue.Kind() != reflect.Func {
		log.Fatal("Provided target is not a function")
		return nil
	}

	// Prepare a WaitGroup and a channel for collecting results
	var wg sync.WaitGroup
	results := make(chan Result, len(inputs))

	// Launch goroutines for each set of inputs
	for _, params := range inputs {
		wg.Add(1)
		go func(params []interface{}) {
			defer wg.Done()

			// Convert inputs to reflect.Value slice
			inValues := make([]reflect.Value, len(params))
			for i, param := range params {
				inValues[i] = reflect.ValueOf(param)
			}

			// Call the function and capture its return values
			callResults := funcValue.Call(inValues)

			// Prepare the output and error (if any)
			var output []interface{}
			var callError error
			for _, res := range callResults {
				if err, ok := res.Interface().(error); ok && err != nil {
					callError = err
				} else {
					output = append(output, res.Interface())
				}
			}

			results <- Result{Output: output, Error: callError}
		}(params)
	}

	// Close the results channel when all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results into a slice
	var finalResults []Result
	for result := range results {
		finalResults = append(finalResults, result)
	}

	return finalResults
}
