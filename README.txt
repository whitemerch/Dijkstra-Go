Voici notre projet de GOLANG portant sur l'algorithme de Dijkstra et sur la résolution de graphe quelque soit leur taille.

Dans le dossier Graph_generator se trouve notre fichier de générateur de graphe avec des liens aléatoires selon le nombre de noeuds qu'on lui donne. 

Dans le dossier Client se trouve notre fichier go client qui est responsable de parser les informations d'un fichier txt et de les envoyer au serveur. D'ailleurs, veuillez une fois le fichier graph généré le copier-coller dans le dossier Client pour plus d'organisation.

Dans le dossier Server se trouve deux dossiers;le premier Server qui représente notre serveur qui contient l'algorithme Dijkstra et qui renverra les solutions à notre client. 
Ce premier dossier est la version Server mais sans goroutine alors que le deuxième Server_go est notre server mais avec goroutines.
Nous avons fait ce choix d'opter pour 2 versions pour bien marquer la différence entre les deux selon le nombre de noeuds dans un graphe.

;)