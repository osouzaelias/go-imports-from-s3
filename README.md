# go-import-from-s3

## Sobre o projeto
Este projeto é uma prova de conceito para validar a importação de um 
arquivo do Bucket S3 para o DynamoDB, usando a feature `Imports from S3` 
conforme apresentado no [blog oficial da AWS.](https://aws.amazon.com/pt/blogs/database/amazon-dynamodb-can-now-import-amazon-s3-data-into-a-new-table/)

## Como executar

Você pode executar esse projeto rodando localmente na sua maquina, desde que os requisitos
abaixo sejam atendidos.

#### Requisitos

- Ler o post no [blog oficial da AWS.](https://aws.amazon.com/pt/blogs/database/amazon-dynamodb-can-now-import-amazon-s3-data-into-a-new-table/)
- Conta na AWS
- Credenciais configurada na maquina (aws config)
- Bucket S3 criado com permissão de leitura, escrita e exclusão 
- Permissão de leitura, escrita, e exclusão no DynamoDB
- Go 1.20.x instalado
- Git para clonar o repositório
- IDE GoLand ou VSCode para editar o projeto

#### Gerando massa para teste
Para gerar massa de teste, você pode usar esta ferramenta   
https://extendsclass.com/csv-generator.html  

#### Configurando o projeto
A configuração do projeto é bem simples. No diretório
`internal` tem o `config.go` nele você poderá definir as
configurações necessárias para o projeto rodar.

## Referencias

Documentação do serviço  
https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/S3DataImport.HowItWorks.html

Boas práticas de uso do serviço  
https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/S3DataImport.BestPractices.html

Referencia para estrutura do projeto    
https://github.com/golang-standards/project-layout