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

	var intReply int
	var token int
	var validateToken bool
	var accountData calc.Account
	var text string
	operation := calc.BalanceOperation{Id: 1, Ammount: 200}
	for text != "EXIT"{
	println("Operações Caixa Automático, digite: ")
	fmt.Println("1: Para realizar depósito")
	fmt.Println("2: Para realizar retirada")
	fmt.Println("3: Para consultar saldo")
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
	if text == "2" {
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
