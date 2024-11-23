# How to Run

``` 
colima start --with-kubernetes

kubectl apply -f deployment.yaml &&
kubectl apply -f service.yaml &&
kubectl apply -f configmap.yaml &&
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/cloud/deploy.yaml &&
kubectl apply -f service-ingress.yaml &&
kubectl apply -f ingress.yaml &&
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml &&
kubectl apply -f hpa.yaml



```

# First use

```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/cloud/deploy.yaml

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: admin-user
  namespace: kubernetes-dashboard
EOF

cat <<EOF | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: admin-user-binding
subjects:
  - kind: ServiceAccount
    name: admin-user
    namespace: kubernetes-dashboard
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io
EOF



```

# To access the dashboard

```
kubectl proxy
kubectl -n kubernetes-dashboard create token admin-user

```