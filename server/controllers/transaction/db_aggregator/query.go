package db_aggregator

import (
	"unicode"
)

// @Internal
// Saves which type requires which convention.
const QUERY_ARG_TYPE_COUNT = 6

var QUERY_ARG_CONVENTION map[QueryArgType]QueryArgConventionType

// @Internal
// Initialize QUERY_ARG_CONVENTION
func initQueryHelper() {
	QUERY_ARG_CONVENTION = make(map[QueryArgType]QueryArgConventionType, QUERY_ARG_TYPE_COUNT)
	QUERY_ARG_CONVENTION[Preload] = UpperConvention
	QUERY_ARG_CONVENTION[Where] = LowerConvention
	QUERY_ARG_CONVENTION[Update] = LowerConvention
	QUERY_ARG_CONVENTION[Select] = LowerConvention
	QUERY_ARG_CONVENTION[Lower] = LowerConvention
	QUERY_ARG_CONVENTION[Upper] = UpperConvention
}

// @Internal
// Convert a string arg to a snake_case.
func convertToLower(arg string) string {
	var result string
	for i, ch := range arg {
		if unicode.IsUpper(ch) {
			if i == 0 || arg[i-1] == '.' {
				result = result + string(unicode.ToLower(ch))
			} else {
				result = result + "_" + string(unicode.ToLower(ch))
			}
		} else {
			result = result + string(ch)
		}
	}

	return result
}

// @Internal
// Convert a string arg to a CamelCase.
func convertToUpper(arg string) string {
	var result string
	for i, ch := range arg {
		if i == 0 {
			result = string(unicode.ToUpper(ch))
			continue
		}
		if arg[i-1] == '_' && unicode.IsLetter(ch) {
			result = result + string(unicode.ToUpper(ch))
		} else if unicode.IsLetter(ch) {
			result = result + string(ch)
		}
	}

	return result
}

// @Internal
// Convert batch of string arg to a snake_case.
func convertToLowerBatch(args ...string) []string {
	result := make([]string, len(args))
	for i, arg := range args {
		result[i] = convertToLower(arg)
	}
	return result
}

// @Internal
// Convert batch of string arg to a snake_case.
func convertToUpperBatch(args ...string) []string {
	result := make([]string, len(args))
	for i, arg := range args {
		result[i] = convertToUpper(arg)
	}
	return result
}

// @Internal
// This is a helper function to get query name following the convention.
func getQueryArg(wordType QueryArgType, args ...string) []string {
	if QUERY_ARG_CONVENTION[wordType] == LowerConvention {
		return convertToLowerBatch(args...)
	} else if QUERY_ARG_CONVENTION[wordType] == UpperConvention {
		return convertToUpperBatch(args...)
	}

	return args
}

// @Internal
// This is a helper function to get query arg for select.
func getSelectQueryArg(arg string) string {
	if QUERY_ARG_CONVENTION[Select] == LowerConvention {
		return convertToLower(arg)
	} else if QUERY_ARG_CONVENTION[Select] == UpperConvention {
		return convertToUpper(arg)
	}

	return arg
}

// @Internal
// This is a helper function to get query arg for preload.
func getPreloadQueryArg(arg string) string {
	if QUERY_ARG_CONVENTION[Preload] == LowerConvention {
		return convertToLower(arg)
	} else if QUERY_ARG_CONVENTION[Preload] == UpperConvention {
		return convertToUpper(arg)
	}

	return arg
}
