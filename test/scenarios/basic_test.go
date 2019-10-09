package operatorsource

/* Basic test for Jenkins Operator

Test will perform and verify actions supported for these users:

Operator (Application Manager) User
* Create OpenShift project
* Install Jenkins Operator through OperatorSource
* Create CR for Jenkins instance into OpenShift namespace
* Verify login to Jenkins, redirected through OpenShift
* Verify Jenkins master startup
* Verify Jenkins master configuration is free of warnings

Developer (Application Owner)
* Create pipeline project (oc new-app referencing github repo)
* Verify creating the project automatically creates a job and starts a pipeline
  for the project - the build pipeline should automatically install any jenkins
  agents needed to build (maven, nodejs)
* Verify Deployment to “my-app-dev” OpenShift project
* Verify Deployed App Endpoint/Route is accessible

*/

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const ShellToUse = "bash"

func Shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func Test_OperatorSource_oc_commands(t *testing.T) {

	defer CleanUp(t)
	t.Run("login", func(t *testing.T) { Login(t) })
	t.Run("new project", func(t *testing.T) { NewProject(t) })
	t.Run("install operator source", func(t *testing.T) { InstallOperatorSource(t) })
	t.Run("create operator subscription", func(t *testing.T) { CreateOperatorSubscription(t) })
	t.Run("ephemeral CR", func(t *testing.T) { EphemeralCR(t) })
}

func Login(t *testing.T) {
	// Start - Login to oc
	out, _, err := Shellout("oc login -u " + os.Getenv("OC_LOGIN_USERNAME") + " -p " + os.Getenv("OC_LOGIN_PASSWORD"))
	if err != nil {
		t.Fatalf("error: %v\n", err)
	} else {
		require.True(t, strings.Contains(out, "Login successful."), "Expecting successful login")
	}
}

func NewProject(t *testing.T) {
	// Create namespace/project
	out, _, err := Shellout("oc new-project " + os.Getenv("OC_PROJECT_NAME"))
	if err != nil {
		t.Fatalf("error: %v\n", err)
	} else {
		require.True(t, strings.Contains(out, "Now using project \""+os.Getenv("OC_PROJECT_NAME")+"\""), "Expecting successful project creation")
	}
}

func InstallOperatorSource(t *testing.T) {
	// Install operator from OperatorSource
	out, _, err := Shellout("oc apply -f ../../jenkins-operator-source.yaml")
	if err != nil {
		t.Fatalf("error: %v\n", err)
	} else {
		require.True(t, strings.Contains(out, "operatorsource.operators.coreos.com/openshift-jenkins-operator created"))
	}
}

func CreateOperatorSubscription(t *testing.T) {
}

func EphemeralCR(t *testing.T) {
}

func CleanUp(t *testing.T) {
	// Clean up resources

	// Delete namespace
	out, errout, err := Shellout("oc delete project " + os.Getenv("OC_PROJECT_NAME"))
	if err != nil {
		t.Logf("stdout: %s\n", out)
		t.Logf("stderr: %s\n", errout)
		t.Logf("error: %v\n", err)
	} else {
		t.Logf(out)
	}

	// Delete OperatorSource
	out, errout, err = Shellout("oc delete operatorsource.operators.coreos.com/openshift-jenkins-operator -n=openshift-marketplace")
	if err != nil {
		t.Logf("stdout: %s\n", out)
		t.Logf("stderr: %s\n", errout)
		t.Logf("error: %v\n", err)
	} else {
		t.Logf(out)
	}

}
