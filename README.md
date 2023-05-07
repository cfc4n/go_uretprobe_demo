# go uretprobe演示

## 测试环境
1. UBUNTU 20.04 x86_64
2. Golang 1.18.8
3. Clang version 10.0.0-4ubuntu1

### 测试方法
#### 编译程序
```shell
make
mkdir -p ebpf/bin
clang -D__KERNEL__ -D__ASM_SYSREG_H \
	-D__TARGET_ARCH_x86 \
	-D__BPF_TRACING__ \
	-Wno-unused-value \
	-Wno-pointer-sign \
	-Wno-compare-distinct-pointer-types \
	-Wunused \
	-Wall \
	-Werror \
	-I/lib/modules/$(uname -r)/build/include \
	-I/lib/modules/$(uname -r)/build/include/uapi \
	-I/lib/modules/$(uname -r)/build/include/generated/uapi \
	-I/lib/modules/$(uname -r)/build/arch/x86/include \
	-I/lib/modules/$(uname -r)/build/arch/x86/include/uapi \
	-I/lib/modules/$(uname -r)/build/arch/x86/include/generated \
	-O2 -emit-llvm \
	ebpf/main.c \
	-c -o - | llc -march=bpf -filetype=obj -o ebpf/bin/probe.o
go run github.com/shuLhan/go-bindata/cmd/go-bindata -pkg main -prefix "ebpf/bin" -o "probe.go" "ebpf/bin/probe.o"
go build -gcflags="-N -l" -o bin/main .
go build -gcflags="-N -l" -o bin/demo ./tests/
```

#### 运行eBPF Hook程序
注意，这里有`-e`参数，会启用uprobe + offset 偏移地址来代替uretprobe。
若不加`-e`参数，则使用uretprobe的方式HOOK，也就是说被HOOK的程序`./tests/tests` 运行时，会崩溃，即触发[Go uretprobe的BUG](https://github.com/iovisor/bcc/issues/1320)。
```shell
sudo bin/main -e
Use uprobe+offset address instead of uretprobe:true
assembly code of the hooked function:
0000	LEA R12, [RSP+Reg(0)-0x18]
0005	CMP R12, [R14+0x10]
// ...
0172	ADD RSP, 0x98
0179	RET
017A	MOV [RSP+Reg(0)+0x8], RAX
017F	NOP
0180	CALL .-190149
0185	MOV RAX, [RSP+Reg(0)+0x8]
018A	JMP .-399
2023/05/07 15:29:29 Golang uretprobe hook main.CountCC [RET] at 0x179
2023/05/07 15:29:29 successfully started, head over to /sys/kernel/debug/tracing/trace_pipe
^C
```

#### 查看系统`trace_pipe`
```shell
root@vm-server-2004:/home/cfc4n/project/go_uretprobe_demo# cat /sys/kernel/debug/tracing/trace_pipe
        tests-297431  [000] .... 428400.595371: 0: new countCC detected
           tests-297431  [000] .... 428400.595413: 0: countCC :: num:51, ret_num:0
           tests-297433  [001] .... 428401.595235: 0: new countCC detected
           tests-297433  [001] .N.. 428401.595268: 0: countCC :: num:51, ret_num:0
           tests-298182  [000] .... 428526.185500: 0: new countCC detected
           tests-298182  [000] .... 428526.185542: 0: countCC :: num:51, ret_num:0
           tests-298182  [000] .... 428527.184804: 0: new countCC detected
           tests-298182  [000] .... 428527.184866: 0: countCC :: num:51, ret_num:0
```

#### 运行被HOOK程序
```shell
./bin/demo
NewTestFunc
0
CountCC return :51
NewTestFunc
0
CountCC return :51
^C
```

可以看到，

## 技术文章

技术文章将在「榫卯江湖」微信公众号上发布。
![](https://image.cnxct.com/2022/03/wechat-white-search-no-alpha-1536x722.png)
