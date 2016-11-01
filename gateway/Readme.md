Gateway支持静态配置和动态配置两种方式。
-------------------
  - 静态配置
    使用 ./gateway -c 开启，这种模式下，直接调用和gateway可执行文件同级目录configs下的gateway.toml
  - 动态配置
    使用etcd + confd作为动态配置管理平台，请将config下的confd_gateway.toml复制到/etc/confd/conf.d下，将confd_gateway.conf.tmpl拷贝到/etc/confd/templates下，具体如何使用请参见confd的介绍。请务必确保最终的动态toml配置输出的目录是/etc/gomqtt/gateway.toml