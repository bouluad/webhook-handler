# 🐙 GitHub Webhook Handler Gateway

This project implements a secure, reliable Go-based gateway deployed in Azure Kubernetes Service (AKS) to bridge GitHub Enterprise webhooks with on-premise tools (like Jenkins) via Azure Service Bus (ASB).

## ✨ Features

* **Signature Validation:** Securely validates GitHub's `X-Hub-Signature-256` HMAC to ensure payload integrity and authenticity.
* **Immediate Acknowledge:** Responds to GitHub with `HTTP 200 OK` immediately after validation, preventing timeout issues.
* **Asynchronous Queueing:** Pushes the raw, validated JSON payload to an Azure Service Bus queue in a non-blocking background process.
* **Minimal Footprint:** Built with Go and Alpine for a small, efficient Docker image suitable for AKS.

## 🚀 Getting Started

### Prerequisites

1.  Go 1.21+
2.  Docker
3.  Kubernetes cluster (AKS) with Ingress Controller configured.
4.  Azure Service Bus Queue provisioned.
5.  A GitHub Webhook Secret.

### 1. Configuration (Environment Variables)

The application requires the following environment variables, which must be provided via a Kubernetes Secret in production:

| Variable | Description |
| :--- | :--- |
| `PORT` | The port the service runs on (Default: `8080`). |
| `GITHUB_WEBHOOK_SECRET` | The shared secret configured in your GitHub webhook settings. |
| `AZURE_SERVICE_BUS_CONN_STRING` | Connection string for your Azure Service Bus namespace. |
| `AZURE_SERVICE_BUS_QUEUE_NAME` | The name of the queue where payloads are sent. |

### 2. Build and Push Docker Image

1.  **Build the image:**
    ```bash
    docker build -t webhook-handler:latest .
    ```
2.  **Tag for ACR (Azure Container Registry):**
    ```bash
    docker tag webhook-handler:latest your-acr-name.azurecr.io/webhook-handler:latest
    ```
3.  **Push the image to ACR:**
    ```bash
    docker push your-acr-name.azurecr.io/webhook-handler:latest
    ```

### 3. Deployment to AKS

1.  **Create the Kubernetes Secret:**
    * **NEVER** commit secrets in plain text.
    * Edit `k8s/01-secret.yaml` with your base64 encoded values.
    ```bash
    kubectl apply -f k8s/01-secret.yaml
    ```

2.  **Deploy the Application and Service:**
    ```bash
    kubectl apply -f k8s/02-deployment.yaml
    kubectl apply -f k8s/03-service.yaml
    ```

3.  **Configure Ingress:**
    * Update `k8s/04-ingress.yaml` with your correct host/domain.
    ```bash
    kubectl apply -f k8s/04-ingress.yaml
    ```

## 🧪 Testing

1.  **Health Check:** Access the `/healthz` endpoint via your Ingress URL: `https://webhook.yourdomain.com/healthz` (Should return `OK` with HTTP 200).
2.  **End-to-End Test:** Configure a webhook in your GitHub Enterprise instance pointing to `https://webhook.yourdomain.com/webhook` and verify the following:
    * GitHub receives an immediate HTTP 200 response.
    * Logs in your Pod show "Signature validated" and "Successfully pushed payload to Service Bus."
    * The raw JSON payload appears in the configured Azure Service Bus queue.

---

I've provided all the requested files and instructions, creating a complete and well-documented project.

Is there any further assistance you require with this project, such as troubleshooting tips or the design of the **on-premise queue consumer**?
