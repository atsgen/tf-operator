<<<<<<< HEAD
package utils_test
=======
package utils
>>>>>>> cfe3c72726bf1d76bbae6391e90b50a4cc933bf1

import (
        "testing"
        "os"
        "github.com/atsgen/tf-operator/pkg/values"
<<<<<<< HEAD
        "github.com/atsgen/tf-operator/pkg/utils"
)

func TestGetAdminPassword( t *testing.T ) {
        testPass := utils.GetAdminPassword()
        if testPass != utils.DefaultAdminPassword {
                t.Errorf("FAILED")
        }
}

func TestGetKubernetesAPIServer( t *testing.T ) {
        os.Setenv("KUBERNETES_SERVICE_HOST",  "test")
        testHost := utils.GetKubernetesAPIServer()
        if testHost!="test" { 
                t.Errorf("FAILED")
        }  
}

func TestGetContainerRegistry( t *testing.T ) {
        os.Setenv("CONTAINER_REGISTRY", "test")
        testRegistry := utils.GetContainerRegistry() 
        if testRegistry!="test" {
                t.Errorf("FAILED")
        }  
}

func TestGetContainerPrefix( t *testing.T ) {
        testPrefix := utils.GetContainerPrefix()
        if testPrefix != utils.ContainerPrefixContrail {
                t.Errorf("FAILED")
        }
}
 
func TestGetReleaseTag ( t *testing.T ) {
        testTag :=  utils.GetReleaseTag("auto")
        if testTag != values.TFCurrentRelease {    
=======
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
>>>>>>> cfe3c72726bf1d76bbae6391e90b50a4cc933bf1
                t.Errorf("FAILED")
        }
}

<<<<<<< HEAD
func TestGetKubernetesAPIPort( t *testing.T ) {
        testPort := utils.GetKubernetesAPIPort() 
        if testPort != "6443" {
                t.Errorf("FAILED")
        }
}

=======
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
>>>>>>> cfe3c72726bf1d76bbae6391e90b50a4cc933bf1

