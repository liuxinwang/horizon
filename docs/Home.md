Welcome to the horizon wiki!

### 巡检设计
数据来源：prometheus（自行安装）
* 服务器指标
  * CPU
  * Memory
    * 物理
    * SWAP
  * Disk
  * IOPS
  * Network
* 实例指标
  * 稳定性
    * 死锁
    * 慢SQL
    * 锁等待
    * 大表
    * 并发运行数
  * 参数
    * 连接使用率
    * 缓存命中率
  * 性能
    * QPS/TPS
  * 安全
    * 空/弱口令账号
  * 高可用
    * 高可用
    * 主从/延迟
  * 备份
    * 备份
* 业务指标

* 评分
<table>
    <tr>
        <td>指标一级</td>
        <td>指标二级</td>
        <td>指标三级</td>
        <td>评分规则</td>
    </tr>
    <tr>
        <td rowspan="6">服务器指标</td>
        <td>CPU</td>
        <tr>
            <td rowspan="2">Memory</td>
            <td>物理</td>
            <tr>
                <td>SWAP</td>
            </tr>
        </tr>
        <tr>
            <td>Disk</td>
        </tr>
        <tr>
            <td>IOPS</td>
        </tr>
        <tr>
            <td>Network</td>
        </tr>
    </tr>
    <tr>
        <td rowspan="12">实例指标</td>
        <td rowspan="5">稳定性</td>
        <td>死锁</td>
        <tr>
            <td>慢SQL</td>
        </tr>
        <tr>
            <td>锁等待</td>
        </tr>
        <tr>
            <td>大表</td>
        </tr>
        <tr>
            <td>并发运行数</td>
        </tr>
        <tr>
            <td rowspan="2">参数</td>
            <td>连接使用率</td>
            <tr>
                <td>缓存命中率</td>
            </tr>
        </tr>
        <tr>
            <td rowspan="1">性能</td>
            <td>QPS/TPS</td>
        </tr>
        <tr>
            <td rowspan="1">安全</td>
            <td>空/弱口令账号</td>
        </tr>
        <tr>
            <td rowspan="2">高可用</td>
            <td>高可用</td>
            <tr>
                <td>主从/延迟</td>
            </tr>
        </tr>
        <tr>
            <td rowspan="1">备份</td>
            <td>备份</td>
        </tr>
    </tr>
</table>