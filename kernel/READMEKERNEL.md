Mediano plazo, pasos funcionales que debe cumplir:
cuando llegue una syscall de io, pasar el proceso de running a bloqueado.
luego, en esa funcion q maneja la syscall abrir un hilo y empezar a tomar el tiempo
segun el archivo de configuracion.seguramente en dos if contemple los dos posibles caminos:
si la notificacion del modulo de io(de que termino su IO para ese proceso) llega antes que termine el
timer el proceso pasa de bloqueado a ready, de lo contrario pasara a suspendidoBLoqueado
y en ese caso ocurren dos cosas mas dentro de ese if principa, avisar a memoria
que ese proceso lo "saque" y lo pase a swap; tambien eso lleva a tener mas
espacio libre entonces el largo plazo debera intentar pasar un proceso de new a ready(avisandole)
al largo plazo mediante una tuberia/señal, aca el largo plazo tiene que ver 
primero la cola de suspReady y luego la de New.
En ese punto tenemos el proceso o en la cola de ready (porque la io termino antes
que el timer) o en suspBloqueado porque se acabo el timer estando en bloqueado.
Entonces de todas formas kernel recibira en algun momento la noticia
de que termino la IO de un proceso dado por el endpoint correspondiente.
aca kernel va a evaluar seguro en un if donde primero buscara el pid en la
cola Bloqueados y luego si no lo encuentra(porque paso a ready antes del timer) buscara 
en la cola de SuspBloqueado para pasarlo finalmente a la cola de SuspReady, para que finalmente el 
largo plazo pueda planificarlo segun su algoritmo y pasarlo a ready.
Una vez que el proceso llegue a SUSP. READY tendrá el mismo comportamiento, es decir, utilizará el 
mismo algoritmo que la cola NEW teniendo más prioridad que esta última. De esta manera, ningún proceso que 
esté esperando en la cola de NEW podrá ingresar al sistema si hay al menos un proceso en SUSP. READY