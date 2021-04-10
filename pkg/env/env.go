package env

import (
	"k8s.io/klog/v2"
	"os"
	"strconv"
)

const CustomResourceDefinitionExecutionTimeout = "CRD_EXECUTION_TIMEOUT_IN_SECOND"
const DefaultExecutionDuration = 30


// GetExecutionTimeoutDuration returns the timeout duration of the execution.
func GetExecutionTimeoutDuration() (int64, error) {
	v := os.Getenv(CustomResourceDefinitionExecutionTimeout)
	if v != "" {
		if i, err := strconv.Atoi(v); err != nil {
			return int64(i), err
		} else {
			return DefaultExecutionDuration, nil
		}
	}
	klog.Infof("No env variable getting from '%s' in the container, use default value: %d", CustomResourceDefinitionExecutionTimeout, DefaultExecutionDuration)
	return DefaultExecutionDuration, nil
}
