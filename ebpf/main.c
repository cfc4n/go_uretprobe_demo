#include <uapi/linux/ptrace.h>
#include "include/bpf.h"
#include "include/bpf_helpers.h"
#include "include/bpf_tracing.h"

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
    int num, ret_num;
    num = (int)go_get_argument_by_stack(ctx, 2);

    ret_num = (int)go_get_argument_by_stack(ctx, 3);

    bpf_printk("countCC :: num:%d, ret_num:%d\n", num,  ret_num);
    return 0;
};

// func CountCC(a int) int
SEC("uretprobe/countcc")
int uretprobe_countcc(struct pt_regs *ctx)
{
    bpf_printk("new countCC[RET] detected\n");
    return 0;
};


char _license[] SEC("license") = "GPL";
__u32 _version SEC("version") = 0xFFFFFFFE;
