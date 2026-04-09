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
docker run -it atuador ./app <typeSensor> <id> <serverip>:<portTCP>
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

### 🔹 5. Exemplos de Uso
#### Cliente
