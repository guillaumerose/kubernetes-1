FROM gcr.io/google_containers/kube-controller-manager-amd64:v1.9.0
CMD kube-controller-manager\
         --controllers=*,bootstrapsigner,tokencleaner\
         --kubeconfig=/kubernetes/controller-manager.conf\
         --root-ca-file=/kubernetes/pki/ca.crt\
         --cluster-signing-cert-file=/kubernetes/pki/ca.crt\
         --cluster-signing-key-file=/kubernetes/pki/ca.key\
         --address=127.0.0.1\
         --leader-elect=true\
         --use-service-account-credentials=true\
         --service-account-private-key-file=/kubernetes/pki/sa.key

COPY controller-manager.conf /kubernetes/controller-manager.conf
COPY pki /kubernetes/pki
