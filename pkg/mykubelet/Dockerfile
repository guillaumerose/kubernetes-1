FROM linuxkit/kubelet:11082f2f27278399856d2ec37b112c79bc68c5a6

ENTRYPOINT ["/usr/bin/kubelet.sh"]

COPY kubelet.sh /usr/bin/kubelet.sh
COPY kubelet.conf /kubernetes/kubelet.conf
