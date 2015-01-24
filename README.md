Edrans Smart City Overview:

Edrans Smart City pretende ser un servicio para agilizar las emergencias en una ciudad.
El servicio que se propone es abrir el camino de los vehiculos de emergencia (ambulancias,
patrulleros y camiones de bomberos), para que lleguen de forma rapida y eficiente a los
diversos lugares donde puedan registrarse siniestros.

Esta es una primer prueba de concepto. En la misma se propone una ciudad de fantasia,
simulando interacciones con los semaforos de la ciudad para despejar los caminos.


Edrans Smart City design patterns:

Componentes:
    - algoritmo shortest path variable *
    - interfaz rest
    - inicializacion de grafos random
    - cliente mobile (quizas)
    - cliente web (quizas)

Sobre el grafo:
    Nodo: 
        - Salidas (las llegadas no deberian ser incluidas -por eso el semaforo- y los objectos que esten en los limitees del grafo no deberian tener salidas)
        - un semaforo (el semaforo indica que entrada esta disponible)
        - id
        - coordenadas (quizas)

    Enlaces:
        - nombre (alusion al nombre de calle)
        - origen (nodo)
        - destino (nodo)
        - peso (en segundos a recorrer)
        Nota: los pesos de los enlaces deben ser construidos de forma tal que los mas proximos al centro de la ciudad tengan un peso mayor (zona de mas trafico)

    Semaforo:
        - Entradas (enlaces que llegan al nodo)
        - Entrada Activa
        - Tiempo de pausa (para cambiar de una entrada a otra)
        - Pausado (ciclo intencionalmente deshabilitado para que un vehiculo pase)

Sobre el algoritmo (supongo sera recursivo, cada recursividad debe tener su copia del PATH):
    - llegar de A a B
    - aceptar ID de objeto del grafo (multiples A/B)
    - aceptar coordenadas de objecto (multiples A/B)
    - calcular proximidad de coordenadas leidas con coordenadas del grafo

    SE DEBE SOPORTAR MULTIPLES CAMINOS PARA DIFERENTES SERVICIOS EN SIMULTANEO
    (Ej: un accidente autmovilistico puede requerir los tres servicios al mismo
    tiempo, y no es posible bloquear un servicio al momento de calcular el camino
    para otro servicio)

    Hay que checkear que no se haya pasado ya por un nodo,
    para evitar que el algoritmo "de vueltas en circulos".

Camino a seguir (PATH):
    - Array de los nodos
    - Array de los enlaces (len(enlaces) = len(nodos)-1)
    - Peso standar (la suma de todos los pesos de los enlaces)

CADA VEHICULO (AMBULANCIA, POLICIA, BOMBERO) DEBE TENER SU PROPIO "PESO MINIMO" PARA RECORRER UN ENLACE;
NO ES LA MISMA LA VELOCIDAD A LA QUE PUEDE IR UN PATRULLERO QUE A LA QUE PUEDE IR UN CAMION DE BOMBEROS

La libreria del algoritmo debe devolver error si:
- El id del inicio no existe
- El id del final no existe
- Las coordenadas no corresponden a nodos (en caso de llegar a implementar geolocalizacion)
