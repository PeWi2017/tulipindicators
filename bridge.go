package tulipindicators

/*
 #include "indicators.h"
 #include <stdio.h>

 int bridgeStartFunction(ti_indicator_start_function f, TI_REAL const *options) {

	return f(options);
 }

 int bridgeIndicatorFunction(ti_indicator_function f,
    int size,
    TI_REAL const *const *inputs,
    TI_REAL const *options,
    TI_REAL *const *outputs) {
		return f(size, inputs, options, outputs);
 }
*/
import (
	"C"
)
import (
	"fmt"
)

func indicator(
	numOutputs int,
	startFunc /* unsafe.Pointer, */ C.ti_indicator_start_function,
	indicatorFunc /* unsafe.Pointer, */ C.ti_indicator_function,
	inputs [][]float64,
	options []float64,
) ([][]float64, error) {

	castOptions := castToCDoubleArray(options)
	defer freeCDoubleArray(castOptions)

	castInputs, inputs := castToC2dDoubleArray(inputs)
	defer freeC2dDoubleArray(castInputs, len(inputs))

	outputSizeDiff := C.bridgeStartFunction(startFunc, castOptions)

	inputSize := len(inputs[0])
    outputSize := len(inputs[0]) - int(outputSizeDiff)

	if outputSize < 1 {
		return nil, fmt.Errorf("insufficient inputs")
	}

	outputs := make([][]float64, numOutputs)

	for i := range outputs {
		outputs[i] = make([]float64, outputSize)
	}

	castOutputs, outputs := castToC2dDoubleArray(outputs)
	defer freeC2dDoubleArray(castOutputs, len(outputs))

	castOutputSize := C.int(inputSize)

	doResponse, doError := C.bridgeIndicatorFunction(
		indicatorFunc,
		castOutputSize,
		castInputs,
		castOptions,
		castOutputs,
	)

	if doError != nil {
		//skipping error because the output *is* actually valid.  SCARY
		//referencing golang github issue 23468
		return nil, doError
	}

	if doResponse == C.TI_INVALID_OPTION {
		return nil, fmt.Errorf("invalid Option for TulipIndicator")
	}

	outputs = extractOutputs(castOutputs, outputs)

	return outputs, nil
}
