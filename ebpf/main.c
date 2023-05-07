#include <uapi/linux/ptrace.h>
#include "include/bpf.h"
#include "include/bpf_helpers.h"
#include "include/bpf_tracing.h"

#define GO_PARAM1(x) ((x)->ax)
#define GO_PARAM2(x) ((x)->bx)
#define GO_PARAM3(x) ((x)->cx)
#define GO_PARAM4(x) ((x)->di)
#define GO_PARAM5(x) ((x)->si)
#define GO_PARAM6(x) ((x)->r8)
#define GO_PARAM7(x) ((x)->r9)
#define GO_PARAM8(x) ((x)->r10)
#define GO_PARAM9(x) ((x)->r11)
#define GOROUTINE(x) ((x)->r14)
#define GO_SP(x) ((x)->sp)

void* go_get_argument_by_stack(struct pt_regs *ctx, int index) {
    void* ptr = 0;
    bpf_probe_read(&ptr, sizeof(ptr), (void *)(PT_REGS_SP(ctx)+(index*8)));
    return ptr;
}


// func CountCC(a int) int
SEC("uprobe/countcc")
int uprobe_countcc(struct pt_regs *ctx)
{
    bpf_printk("new countCC detected\n");
    int num;
    num = (int)GO_PARAM1(ctx);
    bpf_printk("countCC :: num:%d, ret_num:%d\n", num);
    return 0;
};

SEC("uretprobe/countcc")
int uretprobe_countcc(struct pt_regs *ctx)
{
    bpf_printk("new countCC[RET] detected\n");
    return 0;
};


char _license[] SEC("license") = "GPL";
__u32 _version SEC("version") = 0xFFFFFFFE;
