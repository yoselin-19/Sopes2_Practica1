#include <linux/module.h> 
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/list.h>
#include <linux/types.h>
#include <linux/slab.h>
#include <linux/string.h>
#include <linux/sched.h>
#include <linux/string.h>
#include <linux/fs.h>
#include <linux/seq_file.h>
#include <linux/proc_fs.h>
#include <asm/uaccess.h> 
#include <linux/hugetlb.h>


#define modulo_memoria "memo_201403819_201403862"

struct sysinfo informacion;

/*INFO DEL MODULO DE MEMORIA*/
MODULE_LICENSE("GPL");
MODULE_AUTHOR("Yoselin Lemus 201403819 - Brandon Alvarez 201403862 - Ruben Osorio 201403703");
MODULE_DESCRIPTION("Modulo con descripción de la memoria RAM");


static int escribiendoArchivo(struct seq_file *mifile, void *v){
    si_meminfo(&informacion);
    long memoria_total = (informacion.totalram * 4);
    long memoria_libre = (informacion.freeram * 4);
    long memoria_consumida = memoria_total - memoria_libre;

    seq_printf(mifile, "--------------------------------\n");
    seq_printf(mifile, "- Integrantes:                 -\n");
    seq_printf(mifile, "- Yoselin Lemus   - 201403819  -\n");
    seq_printf(mifile, "- Brandon Alvarez - 201403862  -\n");
    seq_printf(mifile, "- Ruben Osorio    - 201403703  -\n");
    seq_printf(mifile, "- Memoria Total: %lu MB        -\n", memoria_total / 1024);
    seq_printf(mifile, "- Memoria Libre: %lu MB        -\n", memoria_libre / 1024);
    seq_printf(mifile, "- Memoria Consumida: %lu MB        -\n", memoria_consumida / 1024);
    seq_printf(mifile, "- Memoria en Uso: %i %%        -\n", (memoria_libre * 100)/memoria_total);
    seq_printf(mifile, "--------------------------------\n");
    return 0;
}

static int alAbrirArchivo(struct inode *inodo, struct file, *mifile){
    return single_open(mifile, escribiendoArchivo, NULL);
}

// static struct file_operations operacionesDeArchivo =
// {
//     .owner = THIS_MODULE,
//     .open = alAbrirArchivo,
//     .read  = seq_read,
//     .llseek = seq_lseek,
//     .release = single_release,
// };

static struct proc_ops operacionesDeArchivo={
    .proc_open = alAbrirArchivo,
    .proc_release = single_release,
    .proc_read = seq_read,
    .proc_lseek = seq_lseek,
    .proc_write = escribiendoArchivo,
};

static int iniciandoModulo(void)
{
    proc_create(modulo_memoria, 0, NULL, &operacionesDeArchivo);
    printk(KERN_INFO "Hola mundo, somos el grupo 14 y este es el monitor de memoria");
    return 0;
}

static void finalizandoModulo(void)
{
    remove_proc_entry(modulo_memoria, NULL);
    printk(KERN_INFO "Sayonara mundo, somos el grupo 14 y este fue el monitor de memoria");
}


module_init(iniciandoModulo);
module_exit(finalizandoModulo);
