package kubernetes

import (
	"context"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/dns-kit/pkg/platform"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

type Kubernetes struct {
	cli *kubernetes.Clientset
	cfg *rest.Config
}

func NewKubernetes(kubeConfig string) (platform.CDNManager, error) {
	cfg, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeConfig))
	if err != nil {
		return nil, err
	}
	cli, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &Kubernetes{cli: cli, cfg: cfg}, nil
}

func (k *Kubernetes) GetCDNList(ctx context.Context) ([]*platform.CDNItem, error) {
	//TODO implement me
	panic("implement me")
}

func (k *Kubernetes) GetCDNDetail(ctx context.Context, domain *platform.CDNItem) (*platform.CDNItem, error) {
	//TODO implement me
	panic("implement me")
}

func (k *Kubernetes) SetCDNHttpsCert(ctx context.Context, domain *platform.CDNItem, cert, key string) error {
	// 从domain.id中取出命名空间和域名
	secretName := strings.ReplaceAll(domain.Domain, ".", "-")
	secretName = strings.ReplaceAll(secretName, "*.", "")
	secretName += "-tls"
	// 判断secret是否存在
	getSecret, err := k.cli.CoreV1().Secrets(domain.ID).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	secret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: domain.ID,
		},
		StringData: map[string]string{
			"tls.crt": cert,
			"tls.key": key,
		},
		Type: "kubernetes.io/tls",
	}
	if getSecret != nil {
		_, err = k.cli.CoreV1().Secrets(domain.ID).Update(ctx, secret, metav1.UpdateOptions{})
	} else {
		_, err = k.cli.CoreV1().Secrets(domain.ID).Create(ctx, secret, metav1.CreateOptions{})
	}
	if err != nil {
		logger.Ctx(ctx).Error("create secret failed", zap.Error(err))
		return err
	}
	return nil
}

func (k *Kubernetes) UserDetails(ctx context.Context) (*platform.User, error) {
	_, err := k.cli.CoreV1().Secrets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	logger.Ctx(ctx).Debug("list namespace")
	return &platform.User{
		ID:       k.cfg.Host,
		Username: k.cfg.Host,
	}, nil
}
