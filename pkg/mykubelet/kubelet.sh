#!/bin/sh
# Kubelet outputs only to stderr, so arrange for everything we do to go there too
exec 1>&2

mkdir -p /var/lib/cni/bin
tar -xzf /root/cni.tgz -C /var/lib/cni/bin

mkdir -p /var/lib/cni/conf

cat <<EOF >/var/lib/cni/conf/10-default.conflist
{
  "cniVersion": "0.3.1",
  "name": "default",
  "plugins": [
    {
      "type": "bridge",
      "bridge": "cni0",
      "isDefaultGateway": true,
      "ipMasq": true,
      "hairpinMode": true,
      "ipam": {
        "type": "host-local",
        "subnet": "10.2.0.0/16",
        "gateway": "10.2.0.1"
      },
      "dns": {
        "nameservers": ["10.2.0.1"]
      }
    },
    {
      "type": "portmap",
      "capabilities": {
        "portMappings": true
      },
      "snat": true
    }
  ]
}
EOF
cat <<EOF >/var/lib/cni/conf/99-loopback.conf
{
  "cniVersion": "0.2.0",
  "type": "loopback"
}
EOF


exec kubelet --kubeconfig=/var/config/kubelet.conf \
    	      --allow-privileged=true \
    	      --cluster-dns=10.96.0.10 \
    	      --cluster-domain=cluster.local \
    	      --cgroups-per-qos=false \
    	      --enforce-node-allocatable= \
    	      --network-plugin=cni \
    	      --pod-cidr=10.2.0.0/16 \
    	      --cni-conf-dir=/etc/cni/net.d \
    	      --cni-bin-dir=/opt/cni/bin \
    	      --cadvisor-port=0 \
    	      --fail-swap-on=false \
    	      --hostname-override node1
