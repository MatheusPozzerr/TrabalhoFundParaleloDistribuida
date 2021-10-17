package calc

import (
	// "fmt"
	"errors"
	"hash/fnv"
	"strconv"
)

//Tipos exportaveis (inicio com letra maiuscula) podem ser registrados por servidores
type Args struct {
	// A, B float64
	TOKEN int //teste temporario
}

/*
	O array memory eh privado do objeto Arith, portanto, o unico modo de acessa-lo eh
	atraves dos metodos exportaveis Store e Load
	Este array pode ser de qualquer tipo, inclusive tipos criados neste arquivo
*/
type Arith struct {
	memory        []Account
	hashTokens    []int // tokens gerados
	receiveTokens []int //tokens de transações utilizadas

	account Account
}

type Account struct {
	Id       int
	Balance  float64
	isActive bool
}

type BalanceOperation struct {
	Id      int
	Ammount float64
}

func hash(s string) uint32 { //metodo para gerar hash mais aleatorios
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (a *Arith) GetHashToken(i *int, reply *int) error { //gera o token de sessão
	vlength := len(a.hashTokens) + 1
	s := strconv.Itoa(vlength)
	s1 := int(hash(s)) // gera o hash e converte em int
	a.hashTokens = append(a.hashTokens, s1)
	println("Token gerado %d", s1)
	*reply = s1
	return nil
}

func contains(s []int, e int) bool { //utilizado para validar o array de tokens e o recebido para transação
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (a *Arith) ValidateToken(args *Args, reply *bool) error { //valida o token para transação
	atualToken := int(args.TOKEN)
	var c bool = contains(a.receiveTokens, atualToken)
	if c == false { //token ainda não existe
		a.receiveTokens = append(a.receiveTokens, atualToken) //adiciona o token aos utilizados em transacao
	}
	*reply = c // true = transação já existe
	return nil
}

/*
	Metodos devem:
	-Pertencer a um tipo exportavel (Arith neste caso) e ser exportaveis
	-Possuir dois parametros de entrada. O primeiro pode ser qualquer tipo
	exportavel ou tipo nativo de go. O segundo de ser obrigatoriamente
	um ponteiro. O segundo argumento eh usado para o retorno do metodo.
	-Retornar um erro. Se for retornado algo alem de nil o cliente recebera
	apenas o erro, sem o ponteiro de reply
*/

func (a *Arith) NewAccount(args *Args, reply *int) error {
	var idAccount int = len(a.memory) + 1
	a.account.Id = idAccount
	a.account.Balance = 0
	a.account.isActive = true
	a.memory = append(a.memory, a.account)
	*reply = idAccount
	return nil
}

func (a *Arith) CloseAccount(idAccountClose *int, reply *int) error {
	var idAccountToClose int = *idAccountClose
	if idAccountToClose > len(a.memory) || idAccountToClose <= 0 {
		*reply = 1
		return nil
	}
	if a.memory[*idAccountClose-1].isActive == false{
		*reply = 2
		return nil
	}
	a.memory[*idAccountClose-1].isActive = false
	*reply = 0
	return nil
}

func (a *Arith) Deposit(operation *BalanceOperation, reply *int) error {
	if operation.Id > len(a.memory) || operation.Id <= 0{
		*reply = 1
		return nil
	}
	if !a.memory[operation.Id-1].isActive{
		return errors.New("erro: a conta foi encerrada")
	}
	a.memory[operation.Id-1].Balance += operation.Ammount
	*reply = 0
	return nil
}

func (a *Arith) Withdraw(operation *BalanceOperation, reply *int) error {
	if operation.Id > len(a.memory) || operation.Id <= 0{
		*reply = 1
		return nil
	}
	if !a.memory[operation.Id-1].isActive{
		return errors.New("erro: a conta foi encerrada")
	}

	a.memory[operation.Id-1].Balance -= operation.Ammount
	*reply = 0
	return nil
}

func (a *Arith) GetInfo(idAccount *int, reply *Account) error {
	if *idAccount > len(a.memory) || *idAccount <= 0{
		*reply = Account{}
		return errors.New("erro: a conta não existe")
	}
	if !a.memory[ *idAccount-1 ].isActive{
		return errors.New("erro: a conta foi encerrada")
	}

	*reply = a.memory[*idAccount-1]
	return nil
}
