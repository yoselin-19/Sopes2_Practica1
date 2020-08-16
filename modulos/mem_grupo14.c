#include <linux/module.h> 
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/list.h>
#include <linux/types.h>
#include <linux/slab.h>
#include <linux/sched.h>
#include <linux/string.h>
#include <linux/fs.h>
#include <linux/seq_file.h>
#include <linux/proc_fs.h>
#include <asm/uaccess.h> 
#include <linux/hugetlb.h>

#define modulo_memoria "mem_grupo14"

struct sysinfo informacion;

/*INFO DEL MODULO DE MEMORIA*/
MODULE_LICENSE("GPL");
MODULE_AUTHOR("Yoselin Lemus 201403819 - Brandon Alvarez 201403862 - Ruben Osorio 201403703");
MODULE_DESCRIPTION("Modulo con descripción de la memoria RAM");


static int escribiendoArchivo(struct seq_file *mifile, void *v){
    #define S(x) ((x) << (PAGE_SHIFT -10))
    si_meminfo(&informacion);
    seq_printf(mifile, "Memoria total: %lu MB\n",S(informacion.totalram/1024));
    seq_printf(mifile, "Memoria consumida: %lu MB\n",S(informacion.totalram/1024) - S(informacion.freeram/1024));
    seq_printf(mifile, "Memoria utilizada: %lu %%\n",S((informacion.freeram)*100)/S(informacion.totalram));        
    return 0;
}

static int alAbrirArchivo(struct inode *inodo, struct file *mifile){
    return single_open(mifile, escribiendoArchivo, NULL);
}

static struct proc_ops operacionesDeArchivo={
    .proc_open = alAbrirArchivo,
    .proc_release = single_release,
    .proc_read = seq_read,
    .proc_lseek = seq_lseek,
};

static int iniciandoModulo(void)
{
    proc_create(modulo_memoria, 0, NULL, &operacionesDeArchivo);
    printk(KERN_INFO "Hola mundo, somos el grupo 14 y este es el monitor de memoria\n");
    return 0;
}

static void finalizandoModulo(void)
{
    remove_proc_entry(modulo_memoria, NULL);
    printk(KERN_INFO "Sayonara mundo, somos el grupo 14 y este fue el monitor de memoria\n");
}


module_init(iniciandoModulo);
module_exit(finalizandoModulo);
