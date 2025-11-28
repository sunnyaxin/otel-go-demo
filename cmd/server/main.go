package main

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
)

func main() {
    s := g.Server()

    s.BindHandler("/hello", func(r *ghttp.Request) {
        r.Response.Write("hello world")
    })

    s.SetPort(8080)
    s.Run()
}
