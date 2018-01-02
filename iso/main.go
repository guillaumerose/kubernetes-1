package main

import (
	"encoding/json"
	"os"
	"github.com/rneugeba/iso9660wrap"
	log "github.com/sirupsen/logrus"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/certs"
	"fmt"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	certutil "k8s.io/client-go/util/cert"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/certs/pkiutil"
	"crypto/rsa"
	"crypto/x509"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	kubeconfigutil "k8s.io/kubernetes/cmd/kubeadm/app/util/kubeconfig"
	"k8s.io/client-go/tools/clientcmd"
	"io/ioutil"
	"path/filepath"
)

func main() {
	iso, err := newIsoDisk()
	if err != nil {
		panic(err)
	}
	fmt.Println(iso)
	iso.writeTo("./tmp/disk.iso")
}

type isoDisk struct {
	content ConfigFile
}

const rights = "0644"

func newIsoDisk() (*isoDisk, error) {
	caCert, caKey, err := certs.NewCACertAndKey()
	if err != nil {
		return nil, err
	}

	apiCert, apiKey, err := certs.NewAPIServerCertAndKey(&kubeadmapi.MasterConfiguration{
		API: kubeadm.API{
			AdvertiseAddress: "0.0.0.0",
		},
		Networking: kubeadm.Networking{
			ServiceSubnet: "10.96.0.0/12",
		},
		APIServerCertSANs: []string{"192.168.65.3", "127.0.0.1"},
	}, caCert, caKey)
	if err != nil {
		return nil, err
	}

	apiClientCert, apiClientKey, err := certs.NewAPIServerKubeletClientCertAndKey(caCert, caKey)
	if err != nil {
		return nil, err
	}

	saSigningKey, err := certs.NewServiceAccountSigningKey()
	if err != nil {
		return nil, err
	}
	publicKey, err := certutil.EncodePublicKeyPEM(&saSigningKey.PublicKey)
	if err != nil {
		return nil, err
	}

	frontProxyCACert, frontProxyCAKey, err := certs.NewFrontProxyCACertAndKey()
	if err != nil {
		return nil, err
	}

	frontProxyClientCert, frontProxyClientKey, err := certs.NewFrontProxyClientCertAndKey(frontProxyCACert, frontProxyCAKey)
	if err != nil {
		return nil, err
	}

	content := ConfigFile{
		"pki": directory(ConfigFile{
			kubeadmconstants.CACertName: file(string(certutil.EncodeCertPEM(caCert))),
			kubeadmconstants.CAKeyName:  file(string(certutil.EncodePrivateKeyPEM(caKey))),

			kubeadmconstants.APIServerCertName: file(string(certutil.EncodeCertPEM(apiCert))),
			kubeadmconstants.APIServerKeyName:  file(string(certutil.EncodePrivateKeyPEM(apiKey))),

			kubeadmconstants.APIServerKubeletClientCertName: file(string(certutil.EncodeCertPEM(apiClientCert))),
			kubeadmconstants.APIServerKubeletClientKeyName:  file(string(certutil.EncodePrivateKeyPEM(apiClientKey))),

			kubeadmconstants.ServiceAccountPrivateKeyName: file(string(certutil.EncodePrivateKeyPEM(saSigningKey))),
			kubeadmconstants.ServiceAccountPublicKeyName:  file(string(publicKey)),

			kubeadmconstants.FrontProxyCACertName: file(string(certutil.EncodeCertPEM(frontProxyCACert))),
			kubeadmconstants.FrontProxyCAKeyName:  file(string(certutil.EncodePrivateKeyPEM(frontProxyCAKey))),

			kubeadmconstants.FrontProxyClientCertName: file(string(certutil.EncodeCertPEM(frontProxyClientCert))),
			kubeadmconstants.FrontProxyClientKeyName:  file(string(certutil.EncodePrivateKeyPEM(frontProxyClientKey))),
		}),
	}

	masterEndpoint := "https://192.168.65.3:6443"

	var kubeConfigSpec = map[string]*kubeConfigSpec{
		kubeadmconstants.AdminKubeConfigFileName: {
			CACert:     caCert,
			APIServer:  masterEndpoint,
			ClientName: "kubernetes-admin",
			ClientCertAuth: &clientCertAuth{
				CAKey:         caKey,
				Organizations: []string{kubeadmconstants.MastersGroup},
			},
		},
		kubeadmconstants.KubeletKubeConfigFileName: {
			CACert:     caCert,
			APIServer:  masterEndpoint,
			ClientName: fmt.Sprintf("system:node:%s", "node1"),
			ClientCertAuth: &clientCertAuth{
				CAKey:         caKey,
				Organizations: []string{kubeadmconstants.NodesGroup},
			},
		},
		kubeadmconstants.KubeletKubeConfigFileName+"2": {
			CACert:     caCert,
			APIServer:  masterEndpoint,
			ClientName: fmt.Sprintf("system:node:%s", "node2"),
			ClientCertAuth: &clientCertAuth{
				CAKey:         caKey,
				Organizations: []string{kubeadmconstants.NodesGroup},
			},
		},
		kubeadmconstants.ControllerManagerKubeConfigFileName: {
			CACert:     caCert,
			APIServer:  masterEndpoint,
			ClientName: kubeadmconstants.ControllerManagerUser,
			ClientCertAuth: &clientCertAuth{
				CAKey: caKey,
			},
		},
		kubeadmconstants.SchedulerKubeConfigFileName: {
			CACert:     caCert,
			APIServer:  masterEndpoint,
			ClientName: kubeadmconstants.SchedulerUser,
			ClientCertAuth: &clientCertAuth{
				CAKey: caKey,
			},
		},
	}

	for path, spec := range kubeConfigSpec {
		config, err := buildKubeConfigFromSpec(spec)
		if err != nil {
			return nil, err
		}
		bin, err := clientcmd.Write(*config)
		if err != nil {
			return nil, err
		}
		content[path] = file(string(bin))
		ioutil.WriteFile(filepath.Join("./tmp", path), bin, 0644)
	}

	return &isoDisk{
		content: content,
	}, nil
}

func directory(content ConfigFile) Entry {
	return Entry{
		Entries: content,
	}
}

func file(content string) Entry {
	return Entry{Perm: rights, Content: &content}
}

func (d *isoDisk) writeTo(path string) error {
	if err := os.Remove(path); err != nil {
		log.Infof("Cannot remove %q: %v", path, err)
	}

	outfh, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	metadataJSON, err := json.Marshal(d.content)
	if err != nil {
		return err
	}

	ioutil.WriteFile("./tmp/metadata.json", metadataJSON, 0644)

	return iso9660wrap.WriteBuffer(outfh, metadataJSON, "config")
}

// ConfigFile represents the configuration file
type ConfigFile map[string]Entry

// Entry represents either a directory or a file
type Entry struct {
	Perm    string           `json:"perm,omitempty"`
	Content *string          `json:"content,omitempty"`
	Entries map[string]Entry `json:"entries,omitempty"`
}

// buildKubeConfigFromSpec creates a kubeconfig object for the given kubeConfigSpec
func buildKubeConfigFromSpec(spec *kubeConfigSpec) (*clientcmdapi.Config, error) {

	// If this kubeconfig should use token
	if spec.TokenAuth != nil {
		// create a kubeconfig with a token
		return kubeconfigutil.CreateWithToken(
			spec.APIServer,
			"kubernetes",
			spec.ClientName,
			certutil.EncodeCertPEM(spec.CACert),
			spec.TokenAuth.Token,
		), nil
	}

	// otherwise, create a client certs
	clientCertConfig := certutil.Config{
		CommonName:   spec.ClientName,
		Organization: spec.ClientCertAuth.Organizations,
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	clientCert, clientKey, err := pkiutil.NewCertAndKey(spec.CACert, spec.ClientCertAuth.CAKey, clientCertConfig)
	if err != nil {
		return nil, fmt.Errorf("failure while creating %s client certificate: %v", spec.ClientName, err)
	}

	// create a kubeconfig with the client certs
	return kubeconfigutil.CreateWithCerts(
		spec.APIServer,
		"kubernetes",
		spec.ClientName,
		certutil.EncodeCertPEM(spec.CACert),
		certutil.EncodePrivateKeyPEM(clientKey),
		certutil.EncodeCertPEM(clientCert),
	), nil
}

// clientCertAuth struct holds info required to build a client certificate to provide authentication info in a kubeconfig object
type clientCertAuth struct {
	CAKey         *rsa.PrivateKey
	Organizations []string
}

// tokenAuth struct holds info required to use a token to provide authentication info in a kubeconfig object
type tokenAuth struct {
	Token string
}

// kubeConfigSpec struct holds info required to build a KubeConfig object
type kubeConfigSpec struct {
	CACert         *x509.Certificate
	APIServer      string
	ClientName     string
	TokenAuth      *tokenAuth
	ClientCertAuth *clientCertAuth
}
