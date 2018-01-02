FROM gcr.io/google_containers/kube-scheduler-amd64:v1.9.0
CMD kube-scheduler \
      --address=127.0.0.1\
      --leader-elect=true\
      --kubeconfig=/var/config/scheduler.conf
