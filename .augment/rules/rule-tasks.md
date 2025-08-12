---
type: "agent_requested"
description: "Sempre que for criar ou executar uma tarefa use esse schema"
---
# Regras para Tasks

## Obrigatório
- Dividir em subtasks atômicas
- Sem testes/documentação
- Sem alias em imports
- `go build ./...` deve passar após cada task
- Seguir boas práticas Go (gofmt, tratamento de erros)

## Template
```yaml
id: <task-id>
title: <título>
desc: <descrição>
acceptance: [<critério>]
subtasks:
  - id: <sub1>
    action: <ação>
```

## Fluxo
1. Task → subtasks atômicas
2. Implementar 1 task
3. `go build ./...` → corrigir se falhar
4. Repetir até concluir