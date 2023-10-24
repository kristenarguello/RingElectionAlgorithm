package main

import (
  "bufio"
  "fmt"
  "log"
  "strconv"
  "strings"
  "sync"
  "time"
  "os"
)

const (
  Finish          = 0 // Controle solicita termino de execucao -- 5
  FailedProcess   = 1 // Controle manda mensagem parq que processo fique falho -- 2
  WorkingProcess  = 2 // Controle torna processo falho para funcional -- 3
  IdentifyFail    = 3 // Controle manda mensagem para um processo informando que aquele falhou -- 1
  Election        = 4 // Troca de mensagens entre processos para escolher um novo lider -- 4
  ElectionConfirmation = 5 // Troca de mensagens avisando de novo lider -- 0
)

/*
DELA - NOSSO
110 = 120
112 = 122
131 = 111
113 = 123
122 = 132
300 = 350
001 = 051
002 = 052
003 = 053
*/

//0 = 3
//1 = 0
//2 = 1
//3 = 2

type mensagem struct {
  tipo int
  idProcInicial int //no corpo
  idProc int //ja estao no corpo
  value int //corpo
  msgFail bool
}

var (
    chans = []chan mensagem {
        make(chan mensagem),
        make(chan mensagem),
        make(chan mensagem),
        make(chan mensagem),
    }
    controle = make(chan bool)
    wg sync.WaitGroup
)

func TranslateMsg(tipo int) string {
  switch tipo {
    case Finish:
    {
      return "FIM"
    }
    case FailedProcess:
    {
      return "Processo falhou"
    }
    case WorkingProcess:
    {
      return "Processo voltou"
    }
    case IdentifyFail:
    {
      return "Idenficacao de falha"
    }
    case Election:
    {
      return "Eleicao"
    }
    case ElectionConfirmation:
    {
      return "Confirmacao de eleicao"
    }
    default:
    {
      return "Invalid Type"
    }
  }
}

func ElectionControler(in chan bool) {
    defer wg.Done()

  file, err := os.Open("commands.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    var comando = []int{}
    for _, val := range strings.Fields(scanner.Text()) {
      c, err := strconv.Atoi(val)
      if err == nil {
        comando = append(comando, c)
      }
    }
    time.Sleep(time.Duration(comando[0]) * time.Second)
    fmt.Printf("\nEnviando comando %d para processo %d \n", comando[1], comando[2])
    
    var msg mensagem
    msg.tipo = comando[1]
    chans[comando[2]]<-msg
    fmt.Printf("Controle: confirma��o %t\n", <-in)
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

    fmt.Println("\n > Processo controlador conclu�do ?\n")
}

func NextIdProcess(id int) int {
  return (id + 1) % len(chans)
}

func ElectionStage(id int, in chan mensagem, initialLeaderId int) {
    defer wg.Done()

    var leaderId int = initialLeaderId
    var fail bool = false

  sendMessage := func(msg mensagem, nextProcId int) {
    msg.idProc = id
    chans[nextProcId]<-msg
  }

  sendElectionConfirmation := func() {
    fmt.Printf("%2d: Enviando confirmacao para demais processos \n", id)
    electedMessage := mensagem{tipo: ElectionConfirmation, idProcInicial: id, value: leaderId}
    sendMessage(electedMessage, NextIdProcess(id))
  }

  processLoop:for {
    message := <-in
    if (!message.msgFail) {
      fmt.Printf("%2d: recebi mensagem %s\n", id, TranslateMsg(message.tipo))
    }

    switch message.tipo {
      case Finish: // Finaliza processo
      {
        controle <- true
        break processLoop
      }
      case FailedProcess: // Controle marca processo como falho
      {
        fail = true
        fmt.Printf("%2d: Processo Falhou \n", id)
        controle <- true
      }
      case WorkingProcess: // Controle marca processo como funcional / inicia eleicao
      {
        fail = false
        fmt.Printf("%2d: Processo Voltou \n", id)
        controle <- true
        electionMessage := mensagem{tipo: Election, idProcInicial: id, value: id}
        sendMessage(electionMessage, NextIdProcess(id))
      }
      case IdentifyFail: // Controle indica para este processo que um processo falhou / Inicia eleicao
      {
        if fail {
          fmt.Printf("%2d: Processo est� falho - ignore \n", id)
          controle <- false
          break
        }
        fmt.Printf("%2d: Informacao de processo falho \n", id)
        controle <- true
        electionMessage := mensagem{tipo: Election, idProcInicial: id, value: id}
        sendMessage(electionMessage, NextIdProcess(id))
      }
      case Election: // Recebe mensagem de eleicao de outro processo
      {
        if fail { // Se processo falho - envia mensagem com erro (simulando uma 'nao-confirmacao' da mensagem)
          fmt.Printf("%2d: Processo est� falho - retornando msg de erro \n", id)
          message.msgFail = true
          sendMessage(message, message.idProc)
        } else if message.msgFail { // Se mensagem de eleicao recebida foi o retorno falho de um envio, tenta enviar novamente para outro processo
          newMsg := mensagem{tipo: message.tipo, idProcInicial: message.idProcInicial, value: message.value}
          nextId := NextIdProcess(message.idProc)
          if (nextId == id) { // Ocorre quando h� somente 1 processo ativo (o mesmo que iniciou a elei��o)
            leaderId = message.value
            fmt.Printf("%2d: Processo eleito: %d \n", id, leaderId)
            sendElectionConfirmation()
          } else {
            sendMessage(newMsg, NextIdProcess(message.idProc))
          }
        } else if message.idProcInicial == id { // Mensagem de eleicao andou todo anel - marca novo lider e envia id de novo lider para demais nodos
          leaderId = message.value
          fmt.Printf("%2d: Processo eleito: %d \n", id, leaderId)
          sendElectionConfirmation()
        } else { // Adiciona valor de id na mensagem (se id for maior id ja na mensagem) e envia mensagem para frente
          if id > message.value {
            message.value = id  
          }
          sendMessage(message, NextIdProcess(id))
        }
      }
      case ElectionConfirmation:
      {
        if fail { // Se processo falho - envia mensagem com erro (simulando uma 'nao-confirmacao' da mensagem)
          fmt.Printf("%2d: Processo est� falho - retornando msg de erro \n", id)
          message.msgFail = true
          sendMessage(message, message.idProc)
        } else if message.msgFail { // Se mensagem de confirmacao de eleicao recebida foi o retorno falho de um envio, tenta enviar novamente para outro processo
          newMsg := mensagem{tipo: message.tipo, idProcInicial: message.idProcInicial, value: message.value}
          nextId := NextIdProcess(message.idProc)
          if nextId == id {
            break
          }
          sendMessage(newMsg, NextIdProcess(message.idProc))
        } else if message.idProcInicial != id {
          fmt.Printf("%2d: Mensagem - Processo Eleito: %d \n", id, message.value)
          leaderId = message.value
          sendMessage(message, NextIdProcess(id))
        }
      }
      default:
      {
        fmt.Printf("%2d: Mensagem desconhecida\n", id)
      }
    }
  }

    fmt.Printf("%2d: terminei \n", id)
}

func main() {
    wg.Add(5) 

    go ElectionStage(0, chans[0], 0)
    go ElectionStage(1, chans[1], 0)
    go ElectionStage(2, chans[2], 0)
    go ElectionStage(3, chans[3], 0)

    go ElectionControler(controle)

    wg.Wait() // Wait for the goroutines to finish
}