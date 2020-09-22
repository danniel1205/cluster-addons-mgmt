# cluster-addons-mgmt

## How to test on local kind cluster

- Create kind cluster: `kind create cluster --config hack/kind/kind.yaml --image kindest/node:v1.18.6`
- Install CRDs: `make install`
- Run controller: `make run`
- Create the addon CRD: `kubectl apply -f config/samples/addons_v1alpha1_certmanageraddon.yaml`
