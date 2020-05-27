# tf-operator
tf-operator is project based on operator-sdk, to manage rollout, upgrades and lifecyle management of Tungsten Fabric CNI
## Operator-SDK setup
There are multiple ways to setup the Operator SDK, while using we can as well download and use the available pre-compiled releases
```
wget https://dl.google.com/go/go1.13.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.13.4.linux-amd64.tar.gz
rm -rf go1.13.4.linux-amd64.tar.gz
echo "export GOPATH=\$HOME/go" >> ~/.bashrc
echo "export PATH=\$PATH:/usr/local/go/bin:\$HOME/go/bin" >> ~/.bashrc
source ~/.bashrc

# get pre-compiled operator-sdk v0.15.2
wget https://github.com/operator-framework/operator-sdk/releases/download/v0.15.2/operator-sdk-v0.15.2-x86_64-linux-gnu
# copy to bin
sudo mv operator-sdk-v0.15.2-x86_64-linux-gnu /usr/local/bin/operator-sdk
sudo chmod +x /usr/local/bin/operator-sdk
operator-sdk version
```

This can also be setup using:
```
scripts/initial-setup.sh
```

## Operator-SDK usage
operator-sdk can be used to add a new api to tf-operator
```
operator-sdk add api --api-version=tungsten.atsgen.com/v1alpha1 --kind TungstenCNI
```
Trigger following to regenerate k8s and crds after modification of any structure of introduced apis
```
operator-sdk generate k8s
operator-sdk generate crds
```
operator-sdk can be used to add a new api to tf-operator
```
operator-sdk add controller --api-version=tungsten.atsgen.com/v1alpha1 --kind=TungstenCNI
```

## Generate Operator image
you can build the tf operator image using following command
```
# expects golang to be available
make
```
alternatively operator-sdk can also build the base alpine image using below command
```
# assumed to be executed from within the top level of repo
operator-sdk build atsgen/tf-operator:v0.0.1
```

## Roll-Out Operator based installation
initialise and run tf-operator
```
# assumed to be executed from within the top level of repo
kubectl create -f deploy/namespace.yaml
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/crds/tungsten.atsgen.com_tungstencnis_crd.yaml
# default password set as 'atsgen'
# user can choose to skip creating secret from here and define
# there own password using
# kubectl create secret generic -n atsgen tungsten-auth --from-literal=password='atsgen'
kubectl create -f deploy/secret.yaml
kubectl create -f deploy/operator.yaml
```

tf operator assumes to enable tungsten fabric controller on master nodes only and vrouter on all the nodes.

Once everything is ready the cluster can be rolled out using
```
kubectl create -f deploy/crds/atsgen.com_v1alpha1_tungstencni_cr.yaml
```

