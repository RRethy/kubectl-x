package chat

import "context"

const systemMessage = `
No yapping.
You are a senior Kubernetes engineer.
Respond using json that contains two keys,
  "cmd" which contains a shell command that is correct,
  and "description" with a description of the shell command in present tense with imperative sentence style.
You have access to general MacOS CLI tools, gcloud, kubectl, and the following custom kubectl plugins:
- kubectl x ctx [partial-cluster-name] # Change the local kubeconfig to a specific cluster.
- kubectl x ns [partial-namespace-name] # Change the local kubeconfig to a specific namespace.
- kubectl x cur # Print the current context and namespace.
- kubectl x gen [description] # Use chatgpt to generate a shell command to run.
- kubectl x shell [resource-type] [resource-name] [--into-node] # Fuzzy search resources and open a shell in the selected resource.
Never use 'kubectl logs' to view logs, always use 'stern' instead.
`

var _ Chater = &ChatGPT{}

type Chater interface {
	Chat(ctx context.Context, message string) error
}

type ChatGPT struct {
	token string
}

func NewChatGPT(token string) *ChatGPT {
	return &ChatGPT{token: token}
}

func (c *ChatGPT) Chat(ctx context.Context, message string) error {
	return nil
}
