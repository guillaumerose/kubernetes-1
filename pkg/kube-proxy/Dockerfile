FROM gcr.io/google_containers/kube-proxy-amd64:v1.9.0
CMD kube-proxy \
      --kubeconfig=/kubernetes/kube-proxy.conf \
      --hostname-override 46316edb4192

COPY kube-proxy.conf /kubernetes/kube-proxy.conf
