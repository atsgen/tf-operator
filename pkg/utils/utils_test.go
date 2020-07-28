package utils

import (
        "testing"
        "os"
        "github.com/atsgen/tf-operator/pkg/values"
)

// Trying to set a value and check for its functionality 
// as part of the testing.

//  Verification for Admin password
func TestAdminPassword( t *testing.T ) {
        testPass := GetAdminPassword()
        SetAdminPassword("password")
        setPass := GetAdminPassword()
        if ((testPass != DefaultAdminPassword) || (setPass != "password")) {
                t.Errorf("FAILED")
        }
        os.Unsetenv(AdminPasswordEnvVar)
}

// Testing for the K8S provider for the deployment
func TestGetKubernetesProvider( t *testing.T ) {
        sampleProvider := GetKubernetesProvider()
        os.Setenv(KubernetesProviderEnvVar, "test")
        testProvider := GetKubernetesProvider()
        if (IsOpenShiftCluster() || (testProvider != "test") || (sampleProvider != "")) {
                t.Errorf("FAILED")
        }
        // Removing the set value
        os.Unsetenv(KubernetesProviderEnvVar)
}

// Checking the status of multus
func TestOpenShiftMultus( t *testing.T ) {
        sampleMultus, _ := IsOpenShiftMultusEnabled()
        SetOpenShiftMultusStatus(true)
        testMultus, _ := IsOpenShiftMultusEnabled()
        SetOpenShiftMultusStatus(false)
        falseMultus, _ := IsOpenShiftMultusEnabled()
        if ((sampleMultus != false) || (testMultus != true) || (falseMultus != false)) {
                t.Errorf("FAILED")
        }
        os.Unsetenv(OpenShiftMultusStatusEnvVar)
}
 
// Testing for the return of out of cluster configuration for K8S API
func TestGetKubernetesAPIServer( t *testing.T ) {
        sampleHost := GetKubernetesAPIServer()
        os.Setenv(KubernetesServiceHostEnvVar, "test")
        testHost := GetKubernetesAPIServer()
        if ((testHost != "test") || (sampleHost != "")) {
                t.Errorf("FAILED")
        }
        // Removing the set value
        os.Unsetenv(KubernetesServiceHostEnvVar)
}

// Testing for container registry 
func TestGetContainerRegistry( t *testing.T ) {
        defRegistry := GetContainerRegistry()
        os.Setenv(ContainerRegistryEnvVar, "test")
        testRegistry := GetContainerRegistry()
        if ((testRegistry != "test") || (defRegistry != DefaultContainerRegistry)) {
                t.Errorf("FAILED")
        }
        os.Unsetenv(ContainerRegistryEnvVar)
}

// Checking for the container prefix
func TestGetContainerPrefix( t *testing.T ) {
        defPrefix := GetContainerPrefix()
        os.Setenv(ContainerPrefixEnvVar, "test")
        testPrefix := GetContainerPrefix()
        if ((defPrefix != ContainerPrefixContrail) || (testPrefix != "test"))  {
                t.Errorf("FAILED")
        }
        os.Unsetenv(ContainerPrefixEnvVar)
}

// Testing for the return of release tag, both default and configured
func TestGetReleaseTag ( t *testing.T ) {
        testTag := GetReleaseTag("")
        configTag := GetReleaseTag("test")
        if ((testTag != values.TFCurrentRelease) || (configTag != "test")) {
                t.Errorf("FAILED")
        }
}

// Checking for K8S API port
func TestGetKubernetesAPIPort( t *testing.T ) {
        defPort := GetKubernetesAPIPort()
        os.Setenv(KubernetesServicePortEnvVar, "1234")
        testPort := GetKubernetesAPIPort()
        if ((defPort != "6443") || (testPort != "1234")) {
                t.Errorf("FAILED")
        }
        os.Unsetenv(KubernetesServicePortEnvVar)
}

// Testing for Tungstenfabric HA status 
func TestTungstenFabricHADisablility( t *testing.T ) {
        os.Setenv(KubernetesProviderEnvVar, "OpenShift")
        sampleHa := IsTungstenFabricHADisabled()
        // Whatever the environment variable value is 
        // HA is disabled
        os.Setenv(DisableTungstenHAEnvVar, "test")
        testHa := IsTungstenFabricHADisabled()
        // Currently HA is supported only for OpenShift cluster
        if ((testHa != true) || (sampleHa != false)) {
                t.Errorf("FAILED")
        }
        os.Unsetenv(DisableTungstenHAEnvVar)
        os.Unsetenv(KubernetesProviderEnvVar)
}

// Testing for the namespace in which operator is running
func TestGetOperatorNamespace( t *testing.T ) {
        sampleNamespace := GetOperatorNamespace()
        os.Setenv(OperatorNamespaceEnvVar, "test")
        testNamespace := GetOperatorNamespace()
        if ((testNamespace != "test") || (sampleNamespace != "")) {
                t.Errorf("FAILED")
        }
        os.Unsetenv(OperatorNamespaceEnvVar)
}

// Testing for resource hack
func TestResourceHack( t *testing.T ) {
        testHack := IsResourceHackDisabled()
        // Irrespective of the environment variable value
        // if it is defined, resource hack is disabled
        os.Setenv(DisableResourceHackEnvVar, "test")
        trueHack := IsResourceHackDisabled()
        if ((testHack != FalseStr) || (trueHack != TrueStr)) {
                t.Errorf("FAILED")
        }
        os.Unsetenv(DisableResourceHackEnvVar)
}
