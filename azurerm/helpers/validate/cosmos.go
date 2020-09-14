package validate

import (
	"fmt"
	"regexp"
)

func CosmosAccountName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// Portal: The value must contain only alphanumeric characters or the following: -
	if matched := regexp.MustCompile("^[-a-z0-9]{3,50}$").Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%s name must be 3 - 50 characters long, contain only letters, numbers and hyphens.", k))
	}

	return warnings, errors
}

func CosmosEntityName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) < 1 || len(value) > 255 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 1 and 255 characters: %q", k, value))
	}

	return warnings, errors
}

func CosmosThroughput(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(int)

	if value < 400 {
		errors = append(errors, fmt.Errorf(
			"%s must be a minimum of 400", k))
	}

	if value%100 != 0 {
		errors = append(errors, fmt.Errorf(
			"%q must be set in increments of 100", k))
	}

	return warnings, errors
}

func CosmosMaxThroughput(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be int", k))
		return
	}

	if v < 4000 {
		errors = append(errors, fmt.Errorf(
			"%s must be a minimum of 4000", k))
	}

	if v > 1000000 {
		errors = append(errors, fmt.Errorf(
			"%s must be a maximum of 1000000", k))
	}

	if v%1000 != 0 {
		errors = append(errors, fmt.Errorf(
			"%q must be set in increments of 1000", k))
	}

	return warnings, errors
}

func CosmosIndexingPolicy(v interface{}, k string) (warnings []string, errors []error) {
	indexingPolicy, ok := v.(map[string]interface{})
	if !ok {
		return warnings, []error{fmt.Errorf("expected type of %q to be map", k)}
	}

	if indexingPolicy == nil {
		return nil, nil
	}

	// Any indexing policy has to include the root path /* as either an included or an excluded path.
	rootPath := "/*"
	rootPathIncluded := false
	var possiblePaths []interface{}

	if includedPaths, exists := indexingPolicy["included_path"]; exists {
		possiblePaths = append(possiblePaths, includedPaths.([]interface{}))
	}

	if excludedPaths, exists := indexingPolicy["excluded_path"]; exists {
		possiblePaths = append(possiblePaths, excludedPaths.([]interface{}))
	}

	for _, i := range possiblePaths {
		if possiblePath, ok := i.(map[string]interface{}); ok {
			if path, ok := possiblePath["path"].(string); ok {
				if path == rootPath {
					rootPathIncluded = true
					break
				}
			}
		}
	}

	if !rootPathIncluded {
		errors = append(errors, fmt.Errorf("the root path '%q' must be included as either an included_path or excluded_path under %q", rootPath, k))
	}

	return warnings, errors
}
