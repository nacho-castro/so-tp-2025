# ORDEN AL LEVANTAR MODULOS

1. MEMORIA
2. KERNEL
3. CPU
4. IO

make corto       # Prueba de planificación de corto plazo
make lym         # Prueba de mediano/largo plazo
make swap        # Prueba de swap
make cache       # Prueba de caché
make eg          # Prueba de estabilidad general
make clean       # Detiene todo lo que esté corriendo

# 🤔 ¿Cuáles serán los valores que usaremos para las conexiones entre módulos? 🤔

## 🧠 CPU
### IP: 127.0.0.1
### Port: 8080
## ⚒️ Kernel 
### IP: 127.0.0.1
### Port: 8081
## 🔌 IO
### IP: 127.0.0.2
### Port: 8082
## 🧰 Memoria
### IP: 127.0.0.1
### Port: 8083

