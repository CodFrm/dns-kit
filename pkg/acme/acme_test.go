package acme

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAcme(t *testing.T) {
	acme, err := NewAcme("yz@ggnb.top", []string{"test2.ggnb.top"})
	assert.NoError(t, err)
	assert.NotNil(t, acme)
	err = acme.GetChallenge()
	assert.NoError(t, err)
	assert.NotNil(t, acme.challenges)
	// 设置dns操作
	// 等待acme验证
	err = acme.WaitChallenge()
	assert.NoError(t, err)
	// 生成证书
	data, err := acme.GetCertificate()
	assert.NoError(t, err)
	assert.NotNil(t, data)
}
