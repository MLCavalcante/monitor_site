package main //indica que este é o pacote principal e seu código começará a ser executado por ele

import (
	"bufio" // importa o pacote bufio, que contém funções para entrada e saída de dados
	"fmt"   //importa o pacote fmt, que contém funções para entrada e saída de dados
	"io"    //importa o pacote io, que contém funções para entrada e saída de dados
	"io/ioutil" //importa o pacote ioutil, que contém funções para entrada e saída de dados
	"net/http" //importa o pacote net/http, que contém funções para trabalhar com a web
	"os"       //importa o pacote os, que contém funções para interagir com o sistema operacional
	"strconv"  //importa o pacote strconv, que contém funções para conversão de tipos
	"strings"  // importa o pacote strings, que contém funções para manipular strings
	"time"     // importa o pacote time, que contém funções para trabalhar com datas e horários
)

const monitoramentos = 2 //constante que indica o número de monitoramentos que serão feitos
const delay = 5 //constante que indica o tempo de espera entre cada monitoramento

func main (){ //indica que este é o ponto de entrada do programa
    exibeIntroducao()
	
	for{
		
		exibeMenu()
		
		comando := leComando()


		switch comando { //Em Go não tem break, ele sai do case automaticamente
		case 1:
			iniciarMonitoramento()
			
		case 2:
			imprimeLogs()
		case 0:
			fmt.Println("Saindo do Programa...")
			os.Exit(0) //0 indica que a saída ocorreu sem erros
		default:	
			fmt.Println("Não conheço este comando")
			os.Exit(-1) //qualquer número diferente de 0 indica que a saída ocorreu com erros
		}
	}



}

func exibeIntroducao() {
	nome := "Gopher"
	versao := 1.2 
	fmt.Println("Olá, sr.", nome)
	fmt.Println("Este programa está na versão", versao)

}
func exibeMenu() {
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir Logs")
	fmt.Println("0 - Sair do Programa")
}

func leComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("O comando escolhido foi:", comandoLido)

	return comandoLido
}

func iniciarMonitoramento() ([]int, []string) {
	fmt.Println("Monitorando...")
	sites := leSitesDoArquivo()
	
	var statusCodes []int
        for i := 0; i < monitoramentos; i++ {
			for _, site := range sites { //O _ é usado para ignorar o índice(i) do elemento. 	                               
				testaSite(site, &statusCodes)//! operador & nos permite obter o endereço de memória de uma variável	
			}
			time.Sleep(delay * time.Second) //time.Sleep() é uma função que faz o programa esperar por um determinado tempo
			fmt.Println("")
		}
	return statusCodes, sites
}                                                    //! O operador * nos permite acessar e modificar o valor 
			                                    //! do slice apontado pelo ponteiro. Isso nos ajuda a manter 
					                            //!o slice dinâmico e atualizado, independentemente de onde o slice é usado no código

func testaSite(site string, statusCodes *[]int){ 
     res, err := http.Get(site)
		if err != nil { //nil é o valor nulo em Go :P
			fmt.Println("Erro ao obter resposta do site:", site)
			*statusCodes = append(*statusCodes, -1) // aqui adicionamos o valor -1 ao slice de statusCodes
			return   
		}

		fmt.Println("Site:", site) // o site que está sendo processado no momento...
	
		if res.StatusCode == 200 {
			fmt.Println("Site:", site, "foi carregado com sucesso!")
			registraLog(site, true)	
		} else {
			fmt.Println("Site:", site, "está com problemas. Status Code:", res.StatusCode)
			registraLog(site, false)	
		}
		*statusCodes = append(*statusCodes, res.StatusCode)
}

func leSitesDoArquivo() []string {
	var sites []string

	arquivo, err := os.Open("sites.txt") //os.Open() é uma função que abre um arquivo

	if err != nil {
			fmt.Println("Ocorreu um erro:", err)
	}

	leitor := bufio.NewReader(arquivo)
    for {
		linha, err := leitor.ReadString('\n') //lê uma string até encontrar o caractere \n
		linha = strings.TrimSpace(linha) //remove os espaços em branco da string
		sites = append(sites, linha)
		if err == io.EOF {
			break
		}	
	}
	arquivo.Close()    
	return sites
}

func registraLog(site string, status bool){ // função que registra os logs de monitoramento
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) //os.OpenFile() é uma função que abre um arquivo
	if err != nil {
		fmt.Println(err)	
	}

	arquivo.WriteString(time.Now().Format("02/01/06 15:04:05") +"-"+ site + "- online" + strconv.FormatBool(status) + "\n" )
	arquivo.Close()

}

func imprimeLogs(){ // função que imprime os logs de monitoramento no terminal
	arquivo, err := ioutil.ReadFile("log.txt") //ioutil.ReadFile() é uma função que lê um arquivo
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(arquivo))
}