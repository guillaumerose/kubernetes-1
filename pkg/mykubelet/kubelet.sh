#!/bin/sh
# Kubelet outputs only to stderr, so arrange for everything we do to go there too
exec 1>&2

if [ ! -e /var/lib/cni/.opt.defaults-extracted ] ; then
    mkdir -p /var/lib/cni/bin
    tar -xzf /root/cni.tgz -C /var/lib/cni/bin
    touch /var/lib/cni/.opt.defaults-extracted
fi

if [ ! -e /var/lib/cni/.cni.conf-extracted ] && [ -d /var/config/cni ] ; then
    mkdir -p /var/lib/cni/conf
    cp /var/config/cni/* /var/lib/cni/conf/
    touch /var/lib/cni/.cni.configs-extracted
fi


exec kubelet --kubeconfig=/kubernetes/kubelet.conf \
    	      --allow-privileged=true \
    	      --cluster-dns=10.96.0.10 \
    	      --cluster-domain=cluster.local \
    	      --cgroups-per-qos=false \
    	      --enforce-node-allocatable= \
    	      --network-plugin=cni \
    	      --cni-conf-dir=/etc/cni/net.d \
    	      --cni-bin-dir=/opt/cni/bin \
    	      --cadvisor-port=0 \
    	      --fail-swap-on=false \
    	      --hostname-override 46316edb4192
