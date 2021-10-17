package main

import (
	"../calc"
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
)

func main() {
	/*
	   Inicializa o cliente na porta 4040 do localhost
	   utilizando o protocolo tcp. Se o servidor estiver
	   em outra maquina deve ser utilizado IP:porta no
	   segundo parametro.
	*/
	c, err := rpc.DialHTTP("tcp", "localhost:4040")
	if err != nil {
		log.Fatal("Error dialing: ", err)
	}

	// var option string

	//Variavel para receber os resultados
	// var reply float64
	var intReply int
	var token int
	var validateToken bool
	var accountData calc.Account
	var text string
	//Estrutura para enviar os numeros
	operation := calc.BalanceOperation{Id: 1, Ammount: 200}
	for text != "EXIT"{
	println("Operações da Agencia, digite: ")
	fmt.Println("1: Para realizar a abertura de conta")
	fmt.Println("2: Para realizar fechamento de conta")
	fmt.Println("3: Para realizar depósito")
	fmt.Println("4: Para realizar retirada")
	fmt.Println("5: Para consultar saldo")
	fmt.Println("6: Para sair do programa")

	e := c.Call("Arith.GetHashToken", 0, &token)
	if e != nil {
		log.Fatal(e)
	}
	// token = 873244444
	args := calc.Args{TOKEN: token}

	reader := bufio.NewReader(os.Stdin)

	text, e := reader.ReadString('\n')
	if e != nil {
		log.Fatal(e)
	}

	text = strings.Replace(text, "\n", "", -1) //Unix

	if(text=="6"){
	println("Finalizando programa...")
	os.Exit(0)
	}
	
	if text == "1" {
		err = c.Call("Arith.ValidateToken", args, &validateToken)
		if err != nil {
			log.Println("Arith error: ", err)
		}
		if validateToken == true {
			log.Println("Token de validação inválido.")
		} else {
			err = c.Call("Arith.NewAccount", args, &intReply)
			if err != nil {
				log.Fatal("Arith error: ", err)
			}
			fmt.Printf("Conta criada, o id da conta é: %d \n", intReply)
		}
	}
	if text == "2" {
		println("Digite o Id da conta a ser encerrada")
		err = c.Call("Arith.ValidateToken", args, &validateToken)
		if err != nil {
			log.Println("Arith error: ", err)
		}
		if validateToken == true {
			log.Println("Token inválido.")
		} else {
			var idConta int
			fmt.Scanf("%d", &idConta)
			err = c.Call("Arith.CloseAccount", idConta, &intReply)
			if err != nil {
				log.Fatal("Arith error: ", err)
			}
			if intReply == 1 {
				fmt.Printf("Conta %d não encontrada.\n", idConta)
			}
			if intReply == 2 {
				fmt.Printf("Conta %d já foi encerrada.\n", idConta)
			}
			if intReply == 0 {
				fmt.Printf("Conta %d encerrada.\n", idConta)
			}
		}
	}

	if text == "3" {
		println("Digite o id da conta")
		err = c.Call("Arith.ValidateToken", args, &validateToken)
		if err != nil {
			log.Println("Arith error: ", err)
		}
		if validateToken == true {
			log.Println("Token inválido.")
		} else {
			var idConta int
			fmt.Scanf("%d", &idConta)
			println("Digite a quantia a ser depositada")
			var ammount float64
			fmt.Scanf("%f", &ammount)
			operation = calc.BalanceOperation{Id: idConta, Ammount: ammount}
			err = c.Call("Arith.Deposit", operation, &intReply)
			if err != nil {
				log.Println(err)
			} else {
				if intReply == 0 {
				fmt.Printf("Deposito de %.2f na conta %d realizado com sucesso \n", ammount, idConta)
				}
				if intReply == 1 {
				println("Erro ao realizar deposito")
				}
			}
		}
	}
	if text == "4" {
		println("Digite o id da conta")
		err = c.Call("Arith.ValidateToken", args, &validateToken)
		if err != nil {
			log.Println("Arith error: ", err)
		}
		if validateToken == true {
			log.Println("Token inválido.")
		} else {
			var idConta int
			fmt.Scanf("%d", &idConta)
			println("Digite a quantia a ser retirada")
			var ammount float64
			fmt.Scanf("%f", &ammount)
			operation = calc.BalanceOperation{Id: idConta, Ammount: ammount}
			err = c.Call("Arith.Withdraw", operation, &intReply)
			if err != nil {
				log.Println(err)
			} else {
				if intReply == 0 {
					fmt.Printf("Retirada de %.2f na conta %d realizado com sucesso\n", ammount, idConta)
				}
				if intReply == 1 {
					println("Erro ao realizar retirada")
				}
			}
		}
	}

	if text == "5" {
		println("Digite o id da conta")
		err = c.Call("Arith.ValidateToken", args, &validateToken)
		if err != nil {
			log.Println("Arith error: ", err)
		}
		if validateToken == true {
			log.Println("Token inválido.")
		} else {
			var idConta int
			fmt.Scanf("%d", &idConta)
			err = c.Call("Arith.GetInfo", idConta, &accountData)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Printf("O saldo da conta %d é: %2.f \n" , accountData.Id, accountData.Balance)
			}
		}
	}
}
}
