FROM linuxkit/kubelet:d581c755f04a8a4060e9947cabe737d6b70fdd1b

ENTRYPOINT ["/usr/bin/kubelet.sh"]

COPY kubelet.sh /usr/bin/kubelet.sh
COPY kubelet.conf /kubernetes/kubelet.conf
