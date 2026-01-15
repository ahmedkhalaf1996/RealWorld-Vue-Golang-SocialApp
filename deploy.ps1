


param([string]$Version = "3.1.0")

$REGISTRY = "ahmedkhalaf666"
$BACKEND = "social-app-backend"
$FRONTEND = "social-app-frontend"

Write-Host "`n=== String Deployment ===" -ForegroundColor Green 
Write-Host "Version: $version" -ForegroundColor Cyan 

# check kubectl 
if (-not (Get-Command kubectl -ErrorAction SilentlyContinue)) {
    Write-Host "ERORR: kubectl not found!" -ForegroundColor Red
    exit 1
}

# check docker 
if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Host "ERORR: docker not found!" -ForegroundColor Red
    exit 1
}

# build images 
Write-Host "`nBuilding backend image..." -ForegroundColor Yellow
docker build -t "${REGISTRY}/${BACKEND}:${Version}" -f backend/api/Dockerfile .

Write-Host "`nBuilding frontend image..." -ForegroundColor Yellow
docker build -t "${REGISTRY}/${FRONTEND}:${Version}" -f frontend/Dockerfile ./frontend

# Push images 
Write-Host "nPushing backend image..." -ForegroundColor Yellow
docker push "${REGISTRY}/${BACKEND}:${Version}"

Write-Host "nPushing frontend image..." -ForegroundColor Yellow
docker push "${REGISTRY}/${FRONTEND}:${Version}"

# update deployment files 
Write-Host "nUpdationg deployment files..." -ForegroundColor Yellow
$backendYaml = Get-Content k8s/backend-deployment.yaml -Raw
$backendYaml = $backendYaml -replace "image:.*social-app-backend.*", "image: ${REGISTRY}/${BACKEND}:${Version}"
$backendYaml | Set-Content k8s/backend-deployment.yaml 

$frontendYaml = Get-Content k8s/frontend-deployment.yaml -Raw
$frontendYaml = $frontendYaml -replace "image:.*social-app-frontend.*", "image: ${REGISTRY}/${FRONTEND}:${Version}"
$frontendYaml | Set-Content k8s/frontend-deployment.yaml 


# deploy to kubernetes 
Write-Host "nDeploying to kubernetes..." -ForegroundColor Yellow 


# 1 Core infrastructure setup 
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/persistent-volumes.yaml

# 2 deploy Data layer mongodb redis kafka 
Write-Host "nDeploying  MongoDB..." -ForegroundColor Cyan 
kubectl apply -f k8s/mongodb-deployment.yaml 

Write-Host "nDeploying  Redis..." -ForegroundColor Cyan 
kubectl apply -f k8s/redis-deployment.yaml 

Write-Host "nDeploying  Kafka..." -ForegroundColor Cyan 
kubectl apply -f k8s/kafka-controller.yaml 
kubectl apply -f k8s/kafka-1.yaml 
kubectl apply -f k8s/kafka-2.yaml
kubectl apply -f k8s/kafka-ui.yaml 

# 3 Wait for infra to be ready 
Write-Host "nWating for infrastructure (60 seconds)..." -ForegroundColor Yellow 
Start-Sleep -Seconds 60

# 4 apply network policies before deploying applications 
Write-Host "Appling network policies..." -ForegroundColor Cyan 
kubectl apply -f k8s/network-policies.yaml 

# 5 deploy application layer (backend, frontend) 
Write-Host "nDeploying  applications..." -ForegroundColor Cyan 
kubectl apply -f k8s/backend-deployment.yaml 
kubectl apply -f k8s/frontend-deployment.yaml 

# 6 apply ingress and autoscaling 
Write-Host "Appling ingress..." -ForegroundColor Cyan 
kubectl apply -f k8s/ingress.yaml 


Write-Host "Appling HPA..." -ForegroundColor Cyan 
kubectl apply -f k8s/hpa.yaml 


# 7 Verification 
Write-Host "n=== Deployment Complete ===" -ForegroundColor Green 
Write-Host "nChecking deployment status..." -ForegroundColor Yellow 

kubectl get pods -n social-app 
kubectl get ingress -n social-app 
kubectl get networkpolicies -n social-app 

