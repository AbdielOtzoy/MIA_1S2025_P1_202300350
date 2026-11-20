# EXT2 file system simulator (v1) — Go + Next.js

The first version of my file system simulator. It creates virtual disks as binary `.mia` files, partitions them with MBR/EBR, formats partitions as **EXT2**, and manages folders, files, users and groups through a web terminal. Reports are rendered with Graphviz.

Course project for *Manejo e Implementación de Archivos* (Universidad de San Carlos de Guatemala, 1st semester 2025). Individual work.

**The complete version is [ext3-filesystem-simulator](https://github.com/AbdielOtzoy/ext3-filesystem-simulator).** It adds EXT3 journaling with crash simulation and recovery, `chmod`/`chown`, copy/move/rename/find, and a graphical file explorer.

## Commands in this version

`mkdisk`, `rmdisk`, `fdisk`, `mount`, `mounted`, `mkfs`, `mkdir`, `mkfile` (with simple indirect blocks), `cat`, `login`, `logout`, `mkgrp`, `rmgrp`, `mkusr`, `rmusr`, `chgrp`, and `rep` (MBR, disk, inode, block, bitmaps, superblock, tree, file, `ls`).

## Run it

```bash
cd backend && go run .                          # Go Fiber API on :8000
cd frontend && pnpm install && pnpm dev         # Next.js terminal on :3000
```

Requires Go 1.23+, Node.js with pnpm, and Graphviz. Example disks are in `backend/disks/`.

Docs (Spanish): [technical manual](docs/ManualTecnico_ext2.pdf) · [user manual](docs/ManuaUsuario_ext2.pdf)
