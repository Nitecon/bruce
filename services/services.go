package services

var ()

func init() {

}

func StartOSServiceExecution() ([]string, error) {
	var failedSvcs []string
	// TODO: Execute services and aggregate the list of ones that fail here

	return failedSvcs, nil
}

func RestoreFailedServices(svcs []string) error {
	// TODO: Replace the existing templates for the failed services and then restart them.

	return nil
}
