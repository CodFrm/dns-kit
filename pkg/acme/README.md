# ACME

你可以看看我的这篇文章，学习一下 ACME 协议：[ACME 协议](https://github.com/CodFrm/blog/blob/main/docs/dev/backend/acme%E5%8D%8F%E8%AE%AE.md)

```go
func main() {
 ctx:=context.Background()
 acme, err := NewAcme("yz@ggnb.top", []string{"test2.ggnb.top"})
 if err!=nil{
        fmt.Println(err)
        return
 }
 challenge,err := acme.GetChallenge(ctx)
 if err != nil {
        fmt.Println(err)
        return
 }
 // 设置dns操作
 time.Sleep(10*time.Second)
 // 等待acme验证
 err = acme.WaitChallenge(ctx)
 if err != nil {
        fmt.Println(err)
        return
 }
 // 生成证书
 data, err := acme.GetCertificate(ctx)
 if err != nil {
        fmt.Println(err)
        return
 }
    fmt.Println(string(data))
}
```
