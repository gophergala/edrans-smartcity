Edrans Smart City (https://github.com/gophergala/edrans-smartcity) Overview:

Edrans Smart City is a proof of concept where a fantasy city is created and a central hub tracks the position of every emergency vehicle's all the time. As soon as any emergency is detected/called, the needed vehicle (ambulance, police car or fire pumper) is notified with a destiny location and a best path to get to it. 

Since this is a city and we have traffic, the shortest path might not be the fastest one to traverse. Moreover, not every road is neither made from the same material nor it measures the same that other one. For this and eventually even more reasons, each road has an associated weight: a 0-99 number which will be modelling these issues. For this particular proof of concept, these values will be randomized. Therefore, the best path will be calculated with a minimum-weight path algorithm.

Finally, since this is a "real" city, at the end of each block we have semaphores which might eventually make the vehicle to stop affecting this best path. Fortunately, the smart city has control over the semaphore lights in the whole city and clears the vehicle's path. 

Once the vehicles are at the emergency place, ambulances will comeback as soon as posible to hospital and police cars will go on patrol.

The idea of the project is to be able to read POST request for emergencies with a REST API from web and mobile clients. Its supposed to provide a web page where you can see a city map with the current status of their semaphores and emergency vehicles.

Once again, this is a proof of concept: its biggest weakness is that it assumes smart semaphores, something which is not real today. But hold on, Internet of Things (IoT) is reaching everywhere, this might not be so crazy!

This project is developed by:
- Dal Lago, Agustin (https://github.com/agudallago)
- Herlein, Christian (https://github.com/ChrisHerlein)


HOW TO USE:

Here is how to start using this prototype web-client.

1) Run server.go to start server. By default this will start listening in port 2489. If you want to listen in another port just run it with "--port=<new>" where <new> is the new port to listen. This


pOST /city para crear ciudad, que te diercciona al mapa
post /emergency para generar emergencias

_____________________________________________________________________________________________________________________


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
