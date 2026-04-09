# iot-client-server
 
# 🏥 Sistema Distribuído para Monitoramento de UTI

Este projeto implementa um sistema distribuído para monitoramento de pacientes em uma Unidade de Terapia Intensiva (UTI), utilizando uma arquitetura baseada em broker para comunicação entre sensores, atuadores e aplicações cliente.

---

## 📌 Funcionalidades

- Monitoramento em tempo real de sinais vitais (BPM e SpO2)
- Comunicação híbrida com UDP (sensores) e TCP (clientes/atuadores)
- Controle remoto de atuadores (alarmes e ventiladores)
- Gerenciamento centralizado via broker
- Suporte a múltiplos dispositivos simultâneos
- Reconexão automática em caso de falhas

---

## 🧠 Arquitetura

O sistema é composto por:

- **Sensores**: enviam dados via UDP
- **Atuadores**: recebem comandos via TCP
- **Clientes**: monitoram dados e enviam comandos
- **Broker**: intermedia toda a comunicação

---

## ⚙️ Tecnologias Utilizadas

- Go (Golang)
- Docker
- TCP/UDP (Sockets)
- JSON (serialização)

---

## 📡 Protocolos de Comunicação

- **UDP**: envio contínuo de dados de sensores (baixa latência)
- **TCP**: comunicação confiável com clientes e atuadores

Formato das mensagens (JSON):

- Tipo
- ID
- Dado
- Ação (somente TCP)

---

## 📂 Estrutura do Projeto

```bash
iot-client-server/
├── .github/workflows/       
├── atuador/                 
│   ├── atuador.go
│   ├── Dockerfile           
│   ├── functionsAtuador.go  
│   └── go.mod               
├── broker/                  
│   ├── broker.go            
│   ├── Dockerfile           
│   ├── go.mod               
│   └── ...                  
├── cliente/                 
│   ├── cliente.go           
│   └── Dockerfile
├── sensor/                  
│   ├── sensor.go            
│   └── Dockerfile
├── docker-compose.yml       
├── makefile                 
└── README.md                
```

## 🚀 Como Executar

### Pré-requisitos

- Docker instalado

---

### 🔹 1. Clone o Repositório
```bash
git clone https://github.com/NycolasDeLima/iot-client-server.git
```

### 🔹 2. Build da imagem

- **Broker**:
```bash
cd broker
docker build -t broker .
```

- **Cliente**:
```bash
cd cliente
docker build -t cliente .
```

- **Atuador**:
```bash
cd atuador
docker build -t atuador .
```

- **Sensor**:
```bash
cd sensor
docker build -t sensor .
```
### 🔹 3. Execute os Containers

**Variáveis de Ambiente**
- serverIP: IP do Broker
- portUDP: Porta para comunicação UDP
- portTCP: Porta para comunicação TCP
- id: Identificador do Dispositivo
- typeSensor: Tipo do Sensor (bpm ou spo2)
- typeAtuador: Tipo do Atuador (alarme ou vmi)

- **Broker**:
```bash
cd broker
docker run -p <portUDP>:<portUDP>/udp -p <portTCP>:<portTCP>/tcp broker ./app <portUDP> <portTCP>
```

- **Cliente**:
```bash
cd cliente
docker run -it cliente ./app <id> <serverIP>:<portTCP>
```

- **Atuador**:
```bash
cd atuador
docker run -it atuador ./app <typeAtuador> <id> <serverip>:<portTCP>
```

- **Sensor**:
```bash
cd sensor
docker run -it sensor ./app <typeSensor> <id> <serverip>:<portUDP>
```

### 🔹 4. Execute os Containers (Makefile)
Caso esteja em um sistem Linux, é possível executar os containers facilmente por meio do Makefile.

**Variáveis do Makefile**
- N: número de containers. (default = 1)
- ip: IP do Broker. (default = localhost)
- types: Tipo do Sensor (bpm ou spo2). (default = bpm)
- typea: Tipo do Atuador (alarme ou vmi). (default = vmi)
- udp: Porta para comunicação UDP. (default = 8080)
- tcp: Porta para comunicação UDP. (default = 8000)

- **Build**:
```bash
make build
```

- **Broker**:
```bash
make broker tcp=<portTCP> udp=<portUDP>
```

- **Cliente**:
```bash
make cliente ip=<serverIP> tcp=<portTCP>
```

- **Atuador**:
```bash
make atuador N=<N> ip=<serverIP> typea=<typeAtuador> tcp=<portTCP>
```

- **Sensor**:
```bash
make atuador N=<N> ip=<serverIP> types=<typeAtuador> udp=<portUDP>
```

**OBS**: Todos os containers criados através do Makefile, com exceção do broker são executados em segundo plano.

---

## Exemplos de Uso

### SENSOR

- BPM:

<img width="822" height="180" alt="image" src="https://github.com/user-attachments/assets/27b247bb-4ebc-4edc-8f40-98114720cfab" />

- SpO2:

<img width="310" height="174" alt="image" src="https://github.com/user-attachments/assets/83936cc3-dda3-4640-a8f6-94ac72e396e8" />

### ATUADOR

- VMI:

<img width="739" height="208" alt="image" src="https://github.com/user-attachments/assets/99cb84e9-9174-43ca-83db-0e2ae3dc080f" />

- Alarme:

<img width="350" height="183" alt="image" src="https://github.com/user-attachments/assets/0d4e5ede-496a-4776-bc55-8aa3679b26f9" />

- Desconexão:

<img width="350" height="183" alt="image" src="https://github.com/user-attachments/assets/eb6a6514-99a6-419d-9108-1a0b96a9c40c" />


### CLIENTE

- Listar Sensores Conectados:

<img width="504" height="278" alt="image" src="https://github.com/user-attachments/assets/d2c2ed8b-46a9-48cc-9dcc-55bcc14c07aa" />


- Listar Atuadores Conectados:

<img width="507" height="243" alt="image" src="https://github.com/user-attachments/assets/fd1fa3ce-49a6-41be-b85c-8e37aa4ff6fc" />


- Ver Dados do Sensor:

<img width="402" height="50" alt="image" src="https://github.com/user-attachments/assets/0190bdb9-a317-4ea7-a67e-65a366c72052" />


- Enviar Comando ao Atuador:

<img width="404" height="442" alt="image" src="https://github.com/user-attachments/assets/8ccc4f54-871a-41e2-bf4c-6662b20330ad" />

<img width="532" height="236" alt="image" src="https://github.com/user-attachments/assets/cf2ff82f-4e57-4ebb-96ce-57ee82f53378" />


- Desconexão com Broker:

<img width="528" height="121" alt="image" src="https://github.com/user-attachments/assets/04db5e4e-758b-43ac-82e2-b60ad05b11e1" />

### BROKER

- Dispositivos conectados:

<img width="662" height="115" alt="image" src="https://github.com/user-attachments/assets/24ee5323-cb33-4d88-900b-615fc3400db3" />


- Dispositivos desconectados:

<img width="662" height="61" alt="image" src="https://github.com/user-attachments/assets/6a25aaa4-9b4a-469f-a461-2d0384c0b424" />


- Ações Cliente:

<img width="662" height="21" alt="image" src="https://github.com/user-attachments/assets/5f5371e2-247d-432e-b559-326fa4c62b6f" />

<img width="662" height="21" alt="image" src="https://github.com/user-attachments/assets/20315168-c718-441c-b886-32277c5d3d44" />

---

## ⚠️ Limitações

- Possível perda de pacotes no uso de UDP
- Interface baseada em terminal
- Falta de Validação dos campos de protocolo
- Verificação de Conexão apenas ao ler/enviar mensagens


