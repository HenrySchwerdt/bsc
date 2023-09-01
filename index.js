var x = 0;
var sum = 0;
while (x <= 5) {
    var xI = 0;
    while(xI <= 5) {
        xI = xI + 1;
        sum = sum + xI;
    }
    x = x + 1;
    sum = sum + x;
}
console.log(sum);
// global _start
// _start:
//     push rbp
//     mov rbp, rsp
//     mov rax, 0
//     push rax
//     mov rax, 0
//     push rax
// .L1:
//     mov rax, QWORD [rbp - 8]
//     push rax
//     mov rax, 3
//     push rax
//     pop rbx
//     pop rax
//     cmp rax, rbx
//     setle al
//     movzx rax, al
//     push rax
//     pop rax
//     test rax, rax
//     jz .E1
//     ; body outer
//     mov rax, 0
//     push rax
// .L2:
//     mov rax, QWORD [rbp - 24]
//     push rax
//     mov rax, 3
//     push rax
//     pop rbx
//     pop rax
//     cmp rax, rbx
//     setle al
//     movzx rax, al
//     push rax
//     pop rax
//     test rax, rax
//     jz .E2
//     ; start of body
//     mov rax, QWORD [rbp - 24]
//     push rax
//     mov rax, 1
//     push rax
//     pop rbx
//     pop rax
//     add rax, rbx
//     push rax
//     pop rax
//     mov QWORD [rbp - 24], rax
//     mov rax, QWORD [rbp - 16]
//     push rax
//     mov rax, QWORD [rbp - 24]
//     push rax
//     pop rbx
//     pop rax
//     add rax, rbx
//     push rax
//     pop rax
//     mov QWORD [rbp - 16], rax
//     ; end body
//     jmp .L2
// .E2:
//     mov rax, QWORD [rbp - 8]
//     push rax
//     mov rax, 1
//     push rax
//     pop rbx
//     pop rax
//     add rax, rbx
//     push rax
//     pop rax
//     mov QWORD [rbp - 8], rax
//     mov rax, QWORD [rbp - 16]
//     push rax
//     mov rax, QWORD [rbp - 8]
//     push rax
//     pop rbx
//     pop rax
//     add rax, rbx
//     push rax
//     pop rax
//     mov QWORD [rbp - 16], rax
//     ; end body outer
//     jmp .L1
// .E1:
//     mov rax, QWORD [rbp - 16]
//     push rax
//     mov rax, 60
//     pop rdi
//     syscall
//     mov rsp, rbp
//     pop rbp
//     mov rax, 60
//     mov rdi, 0
//     syscall

