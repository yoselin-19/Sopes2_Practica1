#include <linux/module.h> 
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/list.h>
#include <linux/slab.h>
#include <linux/sched.h>
#include <linux/string.h>
#include <linux/fs.h>
#include <linux/seq_file.h>
#include <linux/proc_fs.h>
#include <asm/uaccess.h> 
#include <linux/hugetlb.h>
#include <linux/sched/signal.h>
#include <linux/uidgid.h>

#define modulo_cpu "cpu_grupo14"

struct task_struct *task;
struct task_struct *task_child;
struct list_head *list;

/*INFO DEL MODULO DE CPU*/
MODULE_LICENSE("GPL");
MODULE_AUTHOR("Yoselin Lemus 201403819 - Brandon Alvarez 201403862 - Ruben Osorio 201403703");
MODULE_DESCRIPTION("Modulo con descripción del CPU");

int contadorPadre = 1;
int contadorHijos = 1;

static int escribiendoArchivo(struct seq_file *mifile, void *v)
{
    seq_printf(mifile, "{\n");
    seq_printf(mifile, "    \"cpu\" : [\n");
    contadorPadre = 1;
    for_each_process(task)
    {
        if(contadorPadre == 1){
            seq_printf(mifile, "        {\n");
        }
        else{
            seq_printf(mifile, ",\n");
            seq_printf(mifile, "        {\n");
        }
        seq_printf(mifile, "        \"pid\" : %d,\n", task->pid);
		seq_printf(mifile, "        \"uid\" : %d,\n", current_uid());
        seq_printf(mifile, "        \"nombre\" : \"%s\",\n", task->comm);
        seq_printf(mifile, "        \"estado\" : %ld,\n", task->state);
        seq_printf(mifile, "        \"hijos\" : [\n");

        contadorHijos = 1;

        list_for_each(list, &task->children)
        {
            task_child = list_entry(list, struct task_struct, sibling);

            if(contadorHijos == 1){
                seq_printf(mifile, "                    {\n");
            }
            else{
                seq_printf(mifile, ",\n");
                seq_printf(mifile, "                {\n");
            }            
            seq_printf(mifile, "                    \"pid\" : %d,\n", task_child->pid);
            seq_printf(mifile, "                    \"nombre\" : \"%s\",\n", task_child->comm);
            seq_printf(mifile, "                    \"estado\" : %ld\n", task_child->state);
            seq_printf(mifile, "                }");
            contadorHijos++;
        }
        seq_printf(mifile, "\n");
        seq_printf(mifile, "            ]\n");//es para terminar los hijos
        seq_printf(mifile, "        }");

        contadorPadre++;
    }
    seq_printf(mifile, "\n");
    seq_printf(mifile, "    ]\n");
    seq_printf(mifile, "}");
    return 0;
}


static int alAbrirArchivo(struct inode *inodo, struct file *mifile)
{
    return single_open(mifile, escribiendoArchivo, NULL);
}

static struct proc_ops operacionesDeArchivo={
    .proc_open = alAbrirArchivo,
    .proc_release = single_release,
    .proc_read = seq_read,
    .proc_lseek = seq_lseek,
};


static int iniciandoCPU(void)
{
    proc_create(modulo_cpu, 0, NULL, &operacionesDeArchivo);
    printk(KERN_INFO "Yoselin Lemus - Brandon Alvarez - Ruben Osorio\n");
    return 0;
}

static void terminandoCPU(void)
{
    remove_proc_entry(modulo_cpu, NULL);
    printk(KERN_INFO "Sistemas Operativos 2\n");
}

module_init(iniciandoCPU);
module_exit(terminandoCPU);

