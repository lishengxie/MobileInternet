# MobileInternet

运行方法

1. 配置Hyperledger Fabric v2.2测试环境

   [https://hyperledger-fabric.readthedocs.io/en/release-2.2/test_network.html#before-you-begin](https://hyperledger-fabric.readthedocs.io/en/release-2.2/test_network.html#before-you-begin)

2. 安装Go环境(使用go1.14以上版本，可以使用go mod)

3. 将此仓库克隆到本地

   ```
   git clone https://github.com/lishengxie/MobileInternet.git
   ```

4. 进入Hyperledger Fabric v2.2测试环境```cd fabric-samples/test-network```

   * 将```installChaincode.sh```复制到```test-network```目录下
   * 将```deployCC2.sh```复制到目录下```test-network/scripts```目录下
   * 将```SmartReview```文件夹复制到```test-network/chaincode```目录下

5. 运行```installChaincode.sh```配置Hyperledger Fabric测试链及安装和初始化链码。

6. [配置/etc/hosts](https://blog.csdn.net/qq_18807043/article/details/108684285?utm_medium=distribute.pc_relevant.none-task-blog-baidujs_title-1&spm=1001.2101.3001.4242)
   配置/etc/hosts的目的，是让另1个终端SDK程序能够访问到3个节点：peer0.org1.example.com、peer0.org2.example.com、orderer.example.com。
   使用docker exec -it peer0.org1.example.com /bin/ash进入其中1个节点。

   ![在这里插入图片描述](https://img-blog.csdnimg.cn/20200919183019721.png#pic_center)

   在peer0.org1的节点中执行ping peer0.org1.example.com，获得该节点的IP地址。类似地，获得其他2个节点的IP地址。

   ![在这里插入图片描述](https://img-blog.csdnimg.cn/20200919183225449.png#pic_center)

   另开1个终端，修改/etc/hosts文件如下。并确认，这个终端，与上述3个节点的网络是相通的。

   ![在这里插入图片描述](https://img-blog.csdnimg.cn/20200919183413742.png#pic_center)

7. 进入```MobileIntenet```目录下，运行```go run main.go```，在浏览器中输入```localhost:9000```使用。

