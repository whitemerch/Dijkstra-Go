package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

//ro un graph.go <size> <filename.txt>
//Récupère les arguments à l'execution du fichier, vérifie qu'il y ait bien le bon nombre d'arguments
//récupère le nom du fichier a créer ainsi que le nombre de noeuds souhaités
func getArgs() (int, string) {
	if len(os.Args) != 3 {
		fmt.Println("Erreur : utilisez l'appel suivant : go run graph.go <size> <graph.txt>")
		os.Exit(1)
	} else {
		size, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Erreur : utilisez l'appel suivant : go run graph.go <size> <graph.txt>")
			os.Exit(1)
		} else {
			filename := os.Args[2]
			if size <= 2 {
				fmt.Println("Il faut un graphe plus grand que ça")
				os.Exit(1)
			} else {
				return size, filename
			}
		}
	}
	return -1, "" // ne devrait jamais retourner
}

//génère un poids aléatoire entre 1 et 100 (minimum 1)
func randWeight() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(100) + 1
}

//génère une lettre aléatoire appartenant à l'alphabet (le nom d'un des noeuds)
func randLetter(alphabet []int) int {
	rand.Seed(time.Now().UnixNano())
	return alphabet[rand.Intn(len(alphabet))]
}

//supprime l'élement du tableau fourni en fonction de son index
func remove(slice []int, s int) []int {
	end := append(slice[:s], slice[s+1:]...)
	return end
}

//méthode qui appelle remove, pour supprimer un élement connu mais d'index inconnu
func remove_element(slice []int, elt int) []int {
	i := 0
	for slice[i] != elt {
		i++
	}
	return remove(slice, i)
}

//génère le String qui représente le graph qui va être créé, elle prend en paramètre le nombre de noeuds voulus
func generateTie(size int) string {
	var alphabet []int //création de l'alphabet, un tableau d'entier qui contient le nom de tout les noeuds
	for i := 0; i < size; i++ {
		alphabet = append(alphabet, i)
	}

	neighb := make(map[int][]int) //création de la map de voisinage, qui attribue a un noeud tous les noeuds voisins (tous sauf lui même au départ)
	for _, letter := range alphabet {
		neighb[letter] = make([]int, len(alphabet))
		copy(neighb[letter], alphabet)
		remove(neighb[letter], letter)
		neighb[letter] = neighb[letter][:len(neighb[letter])-1]
	}
	//fmt.Printf("neighb %d \n", neighb) //DEBUG
	disp_from := alphabet //contient la liste des lettres de l'alphabet qui peuvent servir comme point de départ d'un lien
	var from, to int      //points de départ et d'arrivée
	var toWrite string
	draw := true                  //initialisations des bouléens, qui déterminent si l'algo doit continuer et si il faut tirer un nouveau point de départ
	for i := 0; i < 4*size; i++ { //on veut ici avoir 4 fois plus de liens que de noeuds, pour fournir un minimum le graph
		if len(disp_from) == 0 {
			break
		}
		alea := rand.Float64()
		if draw || alea < 0.3 { // on détermine que le point de départ sera réutilisé (si possible) avec une probabilité de 1-aléa
			from = randLetter(disp_from) //sinon on tire un nouveau point de départ parmi les disponibles
			draw = false
		}
		//fmt.Println(from) //DEBUG
		to = randLetter(neighb[from]) //on prend une destination parmi les disponibles
		//fmt.Println(to) //DEBUG
		if len(neighb[from]) > 1 { //si le nombre de points disponibles est plus grand que 1
			remove_element(neighb[from], to) //on retire la lettre destination du tableau de voisins
			neighb[from] = neighb[from][:len(neighb[from])-1]
			if len(neighb[to]) > 1 {
				remove_element(neighb[to], from)
				neighb[to] = neighb[to][:len(neighb[to])-1] // retrait de l'inverse selon les memes conditions
				//} else {
				//fmt.Println(disp_from) //DEBUG
				//  remove_element(disp_from, to)
				//  disp_from = disp_from[:len(disp_from)-1]
			}
		} else { //si li n'y a plus de voisins au départ
			if len(disp_from) > 1 { //mais qu'il reste des points de départ possible
				remove_element(disp_from, from) //on le retire de la liste des départs
				disp_from = disp_from[:len(disp_from)-1]
				if len(neighb[to]) > 1 {
					remove_element(neighb[to], from) //on enlève le reverse
					neighb[to] = neighb[to][:len(neighb[to])-1]
				} else {
					remove_element(disp_from, to)
					disp_from = disp_from[:len(disp_from)-1]
				}
				draw = true //on indique qu'il faudra tirer un nouveau point de départ obligatoirement
			}
		}
		//fmt.Printf("neighb %d \n", neighb) //DEBUG
		weight := randWeight()
		toWrite += fmt.Sprintf("%d %d %d\n", from, to, weight)
		alea = rand.Float64()
		if alea < 0.25 {
			toWrite += fmt.Sprintf("%d %d %d\n", to, from, weight)
		} else {
			toWrite += fmt.Sprintf("%d %d %d\n", to, from, randWeight())
		}
	}
	toWrite += ". . ."
	return toWrite
}

//fonction qui va écrire le graph et utiliser écrire le fichier à l'endroit désigné
func writeGraph(size int, path string) {
	fmt.Printf("Création du fichier %v et génaration d'un graph de taille %d \n", path, size)
	f, err := os.OpenFile(path, //ouvre le fichier donné en argument (méthode d'ouverture de fichier généralisée (plus précise que os.Open ou os.Create))
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(generateTie(size)); err != nil {
		log.Println(err)
	}
}

func main() {
	s := time.Now()
	writeGraph(getArgs())
	fmt.Printf("Éxécution en  : %s\n", time.Since(s))
}
