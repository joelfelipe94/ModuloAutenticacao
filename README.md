# Sistema de Autenticação

Implementar uma solução de estrangulamento de um módulo de um monolito para microsserviço.
Onde a primeira etapa é ter o sistema de autenticação separado. Este microsserviço apenas valida login e password.

## Escolha as tecnologias (java, go ou C#)
* Cache
    
    A cache não será implementada já que ela seria usada apenas para esconder a latência de
    acesso ao banco. E neste exemplo simplificado os usuarios estarão mockados. Já que
    em uma aplicação real a gestão das senhas implica em uma série de outras questões 
    como o uso de salt e uma função de hash adequada, comunicação com banco de dados e 
    cadastro dos usuários no sistema. 

* Mensageria

* GraphQL
    GraphQL é uma linguagem de consulta para APIs que oferece flexibilidade e 
    miminiza a quantidade de dados que transmitida entre cliente e servidor. 
    No entanto, como a presente tarefa consiste de uma validação de login e senha,
     essas capacidades não são necessárias neste contexto.
    
* Documentação da API
    A documentação da API será feita com o swagger. Já que ele possui um ecossistema de ferramentas bastante desenvolvido e por permitir que a implementação e a documentação avancem juntas. 

# Pré requisitos:
 * Maior robustez possível;
 * Tempo de resposta abaixo de 50ms
 * Apresentar um teste de carga com x requisições por segundo e y threds