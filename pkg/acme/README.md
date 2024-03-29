# ACME

```go
func main() {
	acme, err := NewAcme("yz@ggnb.top", []string{"test2.ggnb.top"})
	if err!=nil{
        fmt.Println(err)
        return
    }
	err = acme.GetChallenge()
	if err != nil {
        fmt.Println(err)
        return
    }
	// 设置dns操作
	time.Sleep(10*time.Second)
	// 等待acme验证
	err = acme.WaitChallenge()
	if err != nil {
        fmt.Println(err)
        return
    }
	// 生成证书
	data, err := acme.GetCertificate()
	if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(string(data))
}
```