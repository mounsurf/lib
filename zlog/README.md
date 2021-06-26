# glog
simple logger

```Go
package main

import (
    "github.com/mounsurf/lib/zlog"
)

func main() {
    zlog.Debug("123")
    zlog.Warn("456")
    zlog.Info("789")
    zlog.SetConfig(zlog.LevelDebug, "test.txt")
    zlog.Debug("aaa")
    zlog.Info("bbb")
    zlog.Warn("ccc")
}

```
