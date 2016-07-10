GoLinux: set of userspace tools for Linux written in Go.

Hobby project with the goal of learning the basics of Go.
First idea was to do a simple shell in Go (called gosh), but
this has more potential.

First step, implementation of a small init that can be started
by the kernel. Currently it has a built in shell for some basic
commands to interact with the system.

To build the init binary, run:

go build init.go


Create an initramfs with the init file:

mkdir initramfs

rm initramfs.gz
cp init initramfs/init

cd initramfs
find . | cpio -H newc -o > ../initramfs.cpio
cd ..

cat initramfs.cpio | gzip > initramfs.gz
rm initramfs.cpio
rm -rf  initramfs

Start a QEMU VM with a Linux kernel binary and the initramfs:

qemu-system-x86_64 -nographic -kernel /boot/vmlinuz-linux -initrd initramfs.gz -append "console=ttyS0"
