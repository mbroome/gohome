# gohome
Home automation in go.

### git clone
Make sure you git clone the repo into the proper go src tree path.  It is expected
to live under $GOPATH/src/github.com/mbroome/gohome.  Generally,
$GOPATH will be set to your home dir + '/go'.  In my case, /home/mbroome/go.

This is known to work with golang 1.11 though it might work with other version

### deps
```shell
cd $GOPATH/src/github.com/mbroome/gohome
glide up -v
```

### build
```shell
make
```

### run
```shell
./gohome
```


