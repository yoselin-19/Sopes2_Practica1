# Sopes2 practica1 aplicacion web

## Aplicacion Web
---
Permite visualizar gráficas dinámicas que muestren el uso de la memoria RAM del servidor.

Además permite mostrará la información básica de los procesos que se ejecutan y permite terminar los procesos(kill) que se encuentran en ejecución.

## Tecnologia utilizada
---
- Go
- HTML/CSS

## Construir la imagen
---
    docker build -t monitor .
    docker run -it -d -p 3000:3000 --name=MONITOR monitor

## Para usar la imagen de docker hub
---
    docker pull yoseanne/monitor:[tag]
    docker run -it -d -p 3000:3000 --name=MONITOR yoseanne/monitor:[tag]