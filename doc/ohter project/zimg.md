# zimg

## 官方文档
- [安装](http://zimg.buaa.us/documents/install/)
- [使用](http://zimg.buaa.us/documents/guidebookcn/)

## 本人使用

### 安装

```bash
# 安装依赖包
sudo apt-get update
sudo apt-get install openssl cmake libevent-dev libjpeg-dev libgif-dev libpng-dev libwebp-dev libmagickcore-dev libmagickwand-dev libmemcached-dev
```

```bash
# 也可选用源代码安装依赖包
# http://zimg.buaa.us/documents/install/#build-zimg
# 下载zimg并make
git clone https://github.com/buaazp/zimg -b master --depth=1
cd zimg   
make  
# build 可选Storage Backends
# build memcached(optional)
wget http://www.memcached.org/files/memcached-1.4.19.tar.gz
tar zxvf memcached-1.4.19.tar.gz
cd memcached-1.4.19
./configure --prefix=/usr/local
make
make install
# build beansdb(optional)
git clone https://github.com/douban/beansdb
cd beansdb
./configure --prefix=/usr/local
make
# build benseye(optional)
git clone git@github.com:douban/beanseye.git
cd beanseye
make
# build SSDB(optional)
wget --no-check-certificate https://github.com/ideawu/ssdb/archive/master.zip
unzip master
cd ssdb-master
make
# build twemproxy(optional)
git clone git@github.com:twitter/twemproxy.git
cd twemproxy
autoreconf -fvi
./configure --enable-debug=log
make
src/nutcracker -h
```

[zimg build](http://zimg.buaa.us/documents/install/#build-zimg)手动安装依赖包