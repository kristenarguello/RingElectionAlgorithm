// Código exemplo para o trabaho de sistemas distribuidos (eleicao em anel)
// By Cesar De Rose - 2022

// mudar o estado local = ideia do trabalho

package main

import (
	"fmt"
	"sync"
)

/* TIPOS DE MENSAGENS NO CANAL INT (tipo)
2 = processo falho
3 = processo ativo
4 = consumir leitura -> mensagem de término
*/

type mensagem struct {
	tipo  int    // tipo da mensagem para fazer o controle do que fazer (eleição, confirmacao da eleicao)
	corpo [3]int // conteudo da mensagem para colocar os ids (usar um tamanho ocmpativel com o numero de processos no anel)
}

var (
	chans = []chan mensagem{ // vetor de canias para formar o anel de eleicao - chan[0], chan[1] and chan[2] ...
		make(chan mensagem),
		make(chan mensagem),
		make(chan mensagem),
		make(chan mensagem),
	} //como é um vetor formando o anel, basta fazer um for que percorre todos, criando a logica de anel (algoritmo)
	controle = make(chan int)
	wg       sync.WaitGroup // wg is used to wait for the program to finish
)

func ElectionControler(in chan int) {
	defer wg.Done()
	fmt.Println("Controle: iniciando rotina controle de eleição")

	var temp mensagem

	// comandos para o anel iniciam aqui

	// mudar o processo 0 - canal de entrada 3 - para falho (defini mensagem tipo 2 pra isto)
	temp.tipo = 2
	chans[3] <- temp
	fmt.Printf("Controle: mudar o processo 0 para falho\n")
	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	// mudar o processo 1 - canal de entrada 0 - para falho (defini mensagem tipo 2 pra isto)
	temp.tipo = 2
	chans[0] <- temp
	fmt.Printf("Controle: mudar o processo 1 para falho\n")
	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	// matar os outros processos com mensagens não conhecidas (só pra cosumir a leitura)
	temp.tipo = 4
	fmt.Println("Controle: encerrando todos os processos enviando mensagem de termino (codigo 4)")
	for _, c := range chans {
		c <- temp
	}
	// chans[1] <- temp = mesma logica que aqui, so que percorre todos os canais ao invés de um só por vez
	// chans[2] <- temp

	fmt.Println("\n		Controle: Processo controlador concluído ")
}

func ElectionStage(TaskId int, in chan mensagem, out chan mensagem, leader int) {
	defer wg.Done()

	// variaveis locais que indicam se este processo é o lider e se esta ativo
	var actualLeader int
	var bFailed bool = false // todos inciam sem falha

	actualLeader = leader // indicação do lider veio por parâmatro

	var stop bool = false // variável para parar o processo quando chegar no sinal de término de processo
	for !stop {
		temp := <-in // ler mensagem
		fmt.Printf("Election: %2d: recebi mensagem %d, [ %d, %d, %d ]\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2])

		switch temp.tipo {
		case 2:
			{
				bFailed = true
				fmt.Printf("Election: %2d: falho %v \n", TaskId, bFailed)
				fmt.Printf("Election: %2d: lider atual %d\n", TaskId, actualLeader)
				controle <- -5
			}
		case 3:
			{
				bFailed = false
				fmt.Printf("Election: %2d: falho %v \n", TaskId, bFailed)
				fmt.Printf("Election: %2d: lider atual %d\n", TaskId, actualLeader)
				controle <- -5
			}
		case 4: //chegou em um processo com mensagem inserida para consumir leitura
			{
				fmt.Println("Election: chegou em processo com mensagem de termino (codigo 4)")
				stop = true // matar o processo
			}
		default:
			{
				fmt.Printf("Election: %2d: não conheço este tipo de mensagem\n", TaskId)
				fmt.Printf("%2d: lider atual %d\n", TaskId, actualLeader)
			}
		}
	}
	
	fmt.Printf("Election: %2d: terminei \n", TaskId)
}

func main() {

	wg.Add(5) // Add a count of four, one for each goroutine

	// criar os processo do anel de eleicao
	go ElectionStage(0, chans[3], chans[0], 0) // este é o lider
	go ElectionStage(1, chans[0], chans[1], 0) // não é lider, é o processo 0
	go ElectionStage(2, chans[1], chans[2], 0) // não é lider, é o processo 0
	go ElectionStage(3, chans[2], chans[3], 0) // não é lider, é o processo 0

	fmt.Println("\n		Main: Anel de processos criado")

	// criar o processo controlador

	go ElectionControler(controle)
	fmt.Println("\n		Main: Processo controlador criado")

	wg.Wait() // Wait for the goroutines to finish\
	fmt.Println("\n		Main: programa encerrado")
}
