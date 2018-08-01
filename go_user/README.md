## 12306 用户数据导入服务

> ### TODO:
>
> - [x] 消费者实现 pipeline & duplex 通信
> - [ ] Message passing 替换为 Ring buffer
> - [ ] ~~生产者实现 Goroutines 化~~
> - [ ] 更新 README 文档, ongoing.

### How To

#### §. 编译 & 配置

1. build

    执行目录下`make`，在 release 文件夹下编译出所有可执行文件。

    所有构建出来的二进制文件均可通过`./${程序} -v`查看版本和构建时间，如：

    ```
    [root@master:release (master)]# ./import-user -v
    Git Commit Hash: 1.0-0-gf21b9db8493ce8
    UTC Build Time : 2018-07-11_08:41:03AM
    ```

2. 更改配置文件 config.conf

    ```
    [golang]
    # go 最大使用 cpu 数，为了限制资源占用
    go-max-procs 40
    # redis 并发连接数
    thread-pool 512
    # pipeline-size = pipeline-job * pivot
    pipeline-job 1500
    # 读取 File 的 Buffer(8K)，对性能影响较大，请慎重修改！
    file-readbuffer 8192

    [db]
    # redis
    redis-addr 127.0.0.1:19000
    redis-maxidle 512
    redis-maxactive 512
    redis-idletimeout 300
    ```


#### §. 导入

1. 导入用户数据 [import-user]

    * 命令：

        ```
        ./import-user data/WB10_user.bcp 2>> data_error
        ```

    * 示例：

        ```
        [root@KONG:release (master)]# ./import-user data/WB10_user.bcp 2>> data_error
        [16:46:09 CST 2018/07/11] [WARN] (main.main:49) 导入数据开始...
        [16:46:09 CST 2018/07/11] [WARN] (main.main:50) os args: [./import-user data/WB10_user.bcp]
        [16:46:09 CST 2018/07/11] [WARN] (main.main:57) config is: &{GO_MAX_PROCS:32 THREAD_POOL:2560 PIPELINE_JOB:8 FILE_READBUFFER:8192 Redis_Addr:192.168.3.4:19000 Redis_MaxIdle:8192 Redis_MaxActive:8192 Redis_IdleTimeout:300}
        [16:46:11 CST 2018/07/11] [WARN] (main.main:101) 导入数据 10000 lines 总耗时 1672 ms
        ```
    * data_error 为导入失败的原始数据，之后需要用合并工具执行合并操作：[用户异常处理](#§. 异常处理)

    * 导入数据只发生在需要全量重建用户数据的时候

    * 目前 140G 数据文件 (4亿+ 用户，有效 Redis Key 13亿+)，在 30 分钟内可重建完成，Redis 内存占用大小为 480G

    * 在线执行本步骤不会对线上所有数据造成影响，因为数据导入是按用户级别更新的

    * BCP 文件必须为 gb18030 编码

2. 导入用户绑定数据 [import-user-bind]

    * 命令：

        ```
        ./import-user-bind data/WBT20_user_bind.bcp 2>> data_error
        ```

    * 示例：

        ```
        [root@KONG:release (master)]# ./import-user-bind data/WBT20_user_bind.bcp 2>> data_error
        [16:46:09 CST 2018/07/11] [WARN] (main.main:49) 导入数据开始...
        [16:46:09 CST 2018/07/11] [WARN] (main.main:50) os args: [./import-user-bind data/WBT20_user_bind.bcp ]
        [16:46:09 CST 2018/07/11] [WARN] (main.main:57) config is: &{GO_MAX_PROCS:32 THREAD_POOL:2560 PIPELINE_JOB:8 FILE_READBUFFER:8192 Redis_Addr:192.168.3.4:19000 Redis_MaxIdle:8192 Redis_MaxActive:8192 Redis_IdleTimeout:300}
        [16:46:11 CST 2018/07/11] [WARN] (main.main:101) 导入数据 10000 lines 总耗时 1672 ms
        ```
    * data_error 为导入失败的原始数据，之后需要用合并工具执行合并操作：[用户绑定异常处理](#§. 异常处理)

    * 导入数据只发生在需要全量重建用户绑定数据的时候

    * 目前 5G 数据文件 (2000W+ 用户，有效 Redis Key 7000W+)，在 1 分钟内可重建完成，Redis 内存占用大小为 23G

    * 在线执行本步骤不会对线上所有数据造成影响，因为数据导入是按用户级别更新的

    * BCP 文件必须为 gb18030 编码

3. 导入超级用户数据 [import-super-user]

    * 命令：

        ```
        ./import-super-user data/WB10_super_user.bcp
        ```

    * 示例：

        ```
        [root@KONG:release (master)]# ./import-user-bind data/WB10_super_user.bcp
        [16:46:09 CST 2018/07/11] [WARN] (main.main:49) 导入数据开始...
        [16:46:09 CST 2018/07/11] [WARN] (main.main:50) os args: [./import-user-bind data/WB10_super_user.bcp]
        [16:46:09 CST 2018/07/11] [WARN] (main.main:57) config is: &{GO_MAX_PROCS:32 THREAD_POOL:2560 PIPELINE_JOB:8 FILE_READBUFFER:8192 Redis_Addr:192.168.3.4:19000 Redis_MaxIdle:8192 Redis_MaxActive:8192 Redis_IdleTimeout:300}
        [16:46:11 CST 2018/07/11] [WARN] (main.main:101) 导入数据 10000 lines 总耗时 1672 ms
        ```

4. 导入参数表数据 [import-para-define]

    * 命令：

        ```
        ./import-para-define data/WB10_para_define.bcp
        ```

    * 示例：

        ```
        [root@KONG:release (master)]# ./import-para-define data/WB10_para_define.bcp
        [16:46:09 CST 2018/07/11] [WARN] (main.main:49) 导入数据开始...
        [16:46:09 CST 2018/07/11] [WARN] (main.main:50) os args: [./import-para-define data/WB10_para_define.bcp]
        [16:46:09 CST 2018/07/11] [WARN] (main.main:57) config is: &{GO_MAX_PROCS:32 THREAD_POOL:2560 PIPELINE_JOB:8 FILE_READBUFFER:8192 Redis_Addr:192.168.3.4:19000 Redis_MaxIdle:8192 Redis_MaxActive:8192 Redis_IdleTimeout:300}
        [16:46:11 CST 2018/07/11] [WARN] (main.main:101) 导入数据 10000 lines 总耗时 1672 ms
        ```

5. 导入登录数量定义数据 [import-user-number]

    * 命令：

        ```
        ./import-user-number data/WB10_user_number.bcp
        ```

    * 示例：

        ```
        [root@KONG:release (master)]# ./import-para-define data/WB10_para_define.bcp
        [16:46:09 CST 2018/07/11] [WARN] (main.main:49) 导入数据开始...
        [16:46:09 CST 2018/07/11] [WARN] (main.main:50) os args: [./import-user-number data/WB10_user_number.bcp]
        [16:46:09 CST 2018/07/11] [WARN] (main.main:57) config is: &{GO_MAX_PROCS:32 THREAD_POOL:2560 PIPELINE_JOB:8 FILE_READBUFFER:8192 Redis_Addr:192.168.3.4:19000 Redis_MaxIdle:8192 Redis_MaxActive:8192 Redis_IdleTimeout:300}
        [16:46:11 CST 2018/07/11] [WARN] (main.main:101) 导入数据 10000 lines 总耗时 1672 ms
        ```

6. 导入非法手机号数据 [import-illegal-mobile]

    * 命令：

        ```
        ./import-illegal-mobile data/WB10_illegal_mobile.bcp
        ```

    * 示例：

        ```
        [root@KONG:release (master)]# ./import-para-define data/WB10_para_define.bcp
        [16:46:09 CST 2018/07/11] [WARN] (main.main:49) 导入数据开始...
        [16:46:09 CST 2018/07/11] [WARN] (main.main:50) os args: [./import-illegal-mobile data/WB10_illegal_mobile.bcp]
        [16:46:09 CST 2018/07/11] [WARN] (main.main:57) config is: &{GO_MAX_PROCS:32 THREAD_POOL:2560 PIPELINE_JOB:8 FILE_READBUFFER:8192 Redis_Addr:192.168.3.4:19000 Redis_MaxIdle:8192 Redis_MaxActive:8192 Redis_IdleTimeout:300}
        [16:46:11 CST 2018/07/11] [WARN] (main.main:101) 导入数据 10000 lines 总耗时 1672 ms
        ```


#### §. 异常处理

1. 用户表异常数据处理 [inspect-error]

    * 命令：

        ```
        ./inspect-error user_error_sorted user 1>>user_error_sorted_aborted 2>>user_error_sorted_merged
        ```

    * 示例：

        ```
        [root@master:release (master)]# ./inspect-error user_error_sorted user 1>>user_error_sorted_aborted 2>>user_error_sorted_merged
        ```
    * error_abort 为原始数据包含 "\0" 的异常记录，error_merged 为合并后的记录

    * 线上 4亿+ 条数据，共有 4216 条异常数据。其中：

        - [*已解决*]~~只有 2205 条含有 "\0" 的异常记录，建议人工复核后直接修改数据库，彻底解决该问题~~
        - 共有 2011 条断行记录，合并后共有 587 条记录，建议人工复核后导入

    * error_sorted 必须是排过序的异常数据，而且行首应该包含 "行号\x00"

2. 用户绑定表异常数据处理 [inspect-error]

    * 命令：

        ```
        ./inspect-error bind_error_sorted user_bind 1>>bind_error_sorted_aborted 2>>bind_error_sorted_merged
        ```

    * 示例：

        ```
        [root@master:release (master)]# ./inspect-error bind_error_sorted user_bind 1>>bind_error_sorted_aborted 2>>bind_error_sorted_merged
        ```

    * error_abort 为原始数据包含 "\0" 的异常记录，error_merged 为合并后的记录

    * 线上 2000W+ 条数据，共有 4416  条异常数据。其中：

        - 共有 4416 条断行记录，合并后共有 1330 条记录，建议人工复核后导入

    * error_sorted 必须是排过序的异常数据，而且行首应该包含 "行号\x00"


## 初始化导入(全量导入)

> **注：只需要配置 512 个 worker 即可**

### 1. 导入用户数据

```sh
cd release
# 导入数据并捕获异常记录到 user_error
./import-user /work/bcp_files/view_WB10_user_for_codis.bcp 2>>user_error
# 异常记录是乱序的，必须重新排序
sort -k 1 -n user_error >> user_error_sorted
# 处理异常记录
./inspect-error user_error_sorted user 1>>user_error_sorted_aborted 2>>user_error_sorted_merged
# 导入处理过的异常记录
./import-user user_error_sorted_merged
```

### 2. 导入用户绑定数据

```sh
# 导入数据并捕获异常记录到 bind_error
./import-user-bind /work/bcp_files/view_WBT20_user_bind_for_codis.bcp 2>>bind_error
# 异常记录是乱序的，必须重新排序
sort -k 1 -n bind_error >> bind_error_sorted
# 处理异常记录
./inspect-error bind_error_sorted user_bind 1>>bind_error_sorted_aborted 2>>bind_error_sorted_merged
# 导入处理过的异常记录
./import-user-bind bind_error_sorted_merged
```
