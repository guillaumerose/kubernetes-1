FROM gcr.io/google_containers/kube-apiserver-amd64:v1.9.0
CMD kube-apiserver --service-account-key-file=/kubernetes/pki/sa.pub \
        --kubelet-client-certificate=/kubernetes/pki/apiserver-kubelet-client.crt \
        --kubelet-client-key=/kubernetes/pki/apiserver-kubelet-client.key \
        --requestheader-client-ca-file=/kubernetes/pki/front-proxy-ca.crt \
        --proxy-client-cert-file=/kubernetes/pki/front-proxy-client.crt \
        --allow-privileged=true \
        --requestheader-username-headers=X-Remote-User \
        --requestheader-group-headers=X-Remote-Group \
        --service-cluster-ip-range=10.96.0.0/12 \
        --tls-private-key-file=/kubernetes/pki/apiserver.key \
        --secure-port=6443 \
        --requestheader-extra-headers-prefix=X-Remote-Extra- \
        --requestheader-allowed-names=front-proxy-client \
        --advertise-address=0.0.0.0 \
        --client-ca-file=/kubernetes/pki/ca.crt \
        --tls-cert-file=/kubernetes/pki/apiserver.crt \
        --insecure-port=0 \
        --enable-bootstrap-token-auth=true \
        --kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname \
        --admission-control=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,NodeRestriction,ResourceQuota  \
        --proxy-client-key-file=/kubernetes/pki/front-proxy-client.key \
        --authorization-mode=Node,RBAC \
        --etcd-servers=http://127.0.0.1:2379

COPY pki /kubernetes/pki