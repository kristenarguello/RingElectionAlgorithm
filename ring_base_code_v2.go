// Código exemplo para o trabaho de sistemas distribuidos (eleicao em anel)
// By Cesar De Rose - 2022

package main

import (
	"fmt"
	"sync"
	"time"
)

/* TIPOS DE MENSAGENS NO CANAL INT (tipo)
0 = avisando de novo lider = ElectionConfirmation
1 = controle manda mensagem para um processo informando que aquele falhou = IndentifyFail
2 = processo falho = FailedProcess
3 = processo ativo = WorkingProcess
4 = troca de mensagens entre processos para escolher um lider = Election
5 =  consumir leitura -> mensagem de término = Finish
*/

type mensagem struct {
	tipo  int    // tipo da mensagem para fazer o controle do que fazer (eleição, confirmacao da eleicao)
	corpo [3]int // conteudo da mensagem para colocar os ids (usar um tamanho ocmpativel com o numero de processos no anel)
		// idProcInicial, idProc, value
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

	//mudar aqui pra fazer os testes

	// comandos para o anel iniciam aqui

	/*
		processo 0 = canal 3
		processo 1 = canal 0
		processo 2 = canal 1
		processo 3 = canal 2
	*/
	
	//definir o processo 0 como falho
    time.Sleep(time.Duration(1) * time.Second)
	fmt.Printf("\nEnviando comando %d para processo %d \n", 2,0)
	fmt.Print("Controle: definir o processo 1 como falho")
	var msg mensagem
	msg.tipo = 2
	chans[0] <- msg
	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	//definir o processo 2 como falho
    time.Sleep(time.Duration(1) * time.Second)
	fmt.Printf("\nEnviando comando %d para processo %d \n", 2,2)
	fmt.Print("Controle: definir o processo 3 como falho")
	msg = mensagem{}
	msg.tipo = 2
	chans[2] <- msg
	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	//definir o processo 1 como aviso que aquele processo falhou
    time.Sleep(time.Duration(1) * time.Second)
	fmt.Printf("\nEnviando comando %d para processo %d \n", 1,1)
	fmt.Print("Controle: definir o processo 2 como aviso que aquele processo falhou")
	msg = mensagem{}
	msg.tipo = 1
	chans[1] <- msg
	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	//definir o processo 3 como falho
    time.Sleep(time.Duration(1) * time.Second)
	fmt.Printf("\nEnviando comando %d para processo %d \n", 2,3)
	fmt.Print("Controle: definir o processo 0 como falho")
	msg = mensagem{}
	msg.tipo = 2
	chans[3] <- msg
	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	//definir o processo 2 como ativo
    time.Sleep(time.Duration(3) * time.Second)
	fmt.Printf("\nEnviando comando %d para processo %d \n", 3,2)
	fmt.Print("Controle: definir o processo 3 como ativo")
	msg = mensagem{}
	msg.tipo = 3
	chans[2] <- msg
	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	//matar o processo 0
    time.Sleep(time.Duration(5) * time.Second)
	fmt.Printf("\nEnviando comando %d para processo %d \n", 5,0)
	fmt.Print("Controle: matar o processo 1")
	msg = mensagem{}
	msg.tipo = 5
	chans[0] <- msg
	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	//matar o processo 1
    time.Sleep(time.Duration(0) * time.Second)
	fmt.Printf("\nEnviando comando %d para processo %d \n", 5,1)
	fmt.Print("Controle: matar o processo 2")
	msg = mensagem{}
	msg.tipo = 5
	chans[1] <- msg
	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	//matar o processo 2
    time.Sleep(time.Duration(0) * time.Second)
	fmt.Printf("\nEnviando comando %d para processo %d \n", 5,2)
	fmt.Print("Controle: matar o processo 3")
	msg = mensagem{}
	msg.tipo = 5
	chans[2] <- msg
	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	//matar o processo 3
    time.Sleep(time.Duration(0) * time.Second)
	fmt.Printf("\nEnviando comando %d para processo %d \n", 5,3)
	fmt.Print("Controle: matar o processo 0")
	msg = mensagem{}
	msg.tipo = 5
	chans[3] <- msg
	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação


	// // mudar o processo 0 - canal de entrada 3 - para falho (defini mensagem tipo 2 pra isto)
	// temp.tipo = 2
	// chans[3] <- temp
	// fmt.Printf("Controle: mudar o processo 0 para falho\n")
	// fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	// // mudar o processo 1 - canal de entrada 0 - para falho (defini mensagem tipo 2 pra isto)
	// temp.tipo = 2
	// chans[0] <- temp
	// fmt.Printf("Controle: mudar o processo 1 para falho\n")
	// fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	// // matar os outros processos com mensagens não conhecidas (só pra cosumir a leitura)
	// temp.tipo = 4 falho falho ativo deu deadlock
	// fmt.Println("Controle: encerrando todos os processos enviando mensagem de termino (codigo 4)")
	// for _, c := range chans {
	// 	c <- temp
	// }
	// chans[1] <- temp = mesma logica que aqui, so que percorre todos os canais ao invés de um só por vez
	// chans[2] <- temp

	fmt.Println("\n		Controle: Processo controlador concluído ")
}

func ElectionStage(TaskId int, in chan mensagem, out chan mensagem, leader int, anterior chan mensagem) {
	defer wg.Done()

	// variaveis locais que indicam se este processo é o lider e se esta ativo
	var actualLeader int
	var bFailed bool = false // todos inciam sem falha

	actualLeader = leader // indicação do lider veio por parâmatro

	var stop bool = false // variável para parar o processo quando chegar no sinal de término de processo
	for !stop {
		temp := <-in // ler mensagem
		fmt.Printf("\nElection: %2d: recebi mensagem %d, [ %d, %d, %d ]\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2])

		switch temp.tipo {
		case 0: //electionconfirmation = ELECTION WINNER CONFIRMATION
			{
				if bFailed { //nao faz nada se ja estiver falho
					fmt.Printf("%2d: processo falho, lider nao eleito\n", TaskId)					
				} else { //falhou = nao faz nada e continua no anel
					actualLeader = temp.corpo[0]
					fmt.Printf("%2d: processo eleito: %d\n", TaskId, actualLeader)
				}
				out <- temp
				// nao precisa de controle pq nao sai do anel pra chegar no controle
			}
		case 1: // controle indica para este processo que um processo falhou = inicia processo de eleição
			{
				//VOTE REQUEST
				if bFailed {
					fmt.Printf("\t%2d: nao participou da votacao, processo falho\n", TaskId)
				} else {
					temp.corpo[TaskId] = TaskId
					fmt.Printf("%2d: processo votou")
				} 
				out <- temp
				// nao precisa de controle pq nao sai do anel pra chegar no controle

			}
		case 2: // PROCESSO FALHO
			{
				bFailed = true
				fmt.Printf("Election: %2d: falho %v \n", TaskId, bFailed)
				controle <- 1
			}
		case 3: // PROCESSO ATIVO	
			{
				bFailed = false
				fmt.Printf("Election: %2d: falho %v \n", TaskId, bFailed)
				// fmt.Printf("Election: %2d: lider atual %d\n", TaskId, actualLeader)
				// fmt.Printf("Election: %2d: iniciando eleicao por volta a falha\n", TaskId)
				//eleicao aqui
				//controle marca o processo como funcional = inicia a eleição novament
				controle <- 1
			}
		case 4: // processo de eleicao (recebe mensagem para iniciar a eleicao) = INITIATE ELECTION
			{	
				if !bFailed {
					fmt.Printf("Election: %2d: processo falho\n", TaskId)
					//election
					//COLOCAR AQUI A ELEICAO 

				} else {
					fmt.Printf("Election: %2d: processo falho, eleicao nao iniciada\n", TaskId)
					controle <- 0
				}
				
			}
			case 5: // finaliza o processo
			{
				fmt.Println("Election: chegou em processo com mensagem de termino (codigo 4)")
				controle <- 1
				stop = true // matar o processo
			}
		default:
			{
				fmt.Printf("Election: %2d: não conheço este tipo de mensagem\n", TaskId)
				fmt.Printf("%2d: lider atual %d\n", TaskId, actualLeader)
			}
		}
		fmt.Printf("Election: %2d: lider atual %d\n", TaskId, actualLeader)

	}
	
	fmt.Printf("Election: %2d: terminei \n", TaskId)
}

func main() {

	wg.Add(5) // Add a count of four, one for each goroutine

	// criar os processo do anel de eleicao
	go ElectionStage(0, chans[0], chans[1], 0, chans[3]) // este é o lider
	fmt.Println("0:", chans[0])
	go ElectionStage(1, chans[1], chans[2], 0, chans[0]) // não é lider, é o processo 0
	fmt.Println("1:", chans[1])
	go ElectionStage(2, chans[2], chans[3], 0, chans[1]) // não é lider, é o processo 0
	fmt.Println("2:", chans[2])
	go ElectionStage(3, chans[3], chans[0], 0, chans[2]) // não é lider, é o processo 0
	fmt.Println("3:", chans[3])

	fmt.Println("\n		Main: Anel de processos criado")

	// criar o processo controlador

	go ElectionControler(controle)
	fmt.Println("\n		Main: Processo controlador criado")

	wg.Wait() // Wait for the goroutines to finish\
	fmt.Println("\n		Main: programa encerrado")
}
