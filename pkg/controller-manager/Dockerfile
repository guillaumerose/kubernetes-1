FROM gcr.io/google_containers/kube-controller-manager-amd64:v1.9.0
CMD kube-controller-manager\
         --controllers=*,bootstrapsigner,tokencleaner\
         --kubeconfig=/var/config/controller-manager.conf\
         --root-ca-file=/var/config/pki/ca.crt\
         --cluster-signing-cert-file=/var/config/pki/ca.crt\
         --cluster-signing-key-file=/var/config/pki/ca.key\
         --address=127.0.0.1\
         --leader-elect=true\
         --use-service-account-credentials=true\
         --service-account-private-key-file=/var/config/pki/sa.key
