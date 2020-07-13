package utils_test

import (
        "testing"
        "os"
        "github.com/atsgen/tf-operator/pkg/values"
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
                t.Errorf("FAILED")
        }
}

func TestGetKubernetesAPIPort( t *testing.T ) {
        testPort := utils.GetKubernetesAPIPort() 
        if testPort != "6443" {
                t.Errorf("FAILED")
        }
}


