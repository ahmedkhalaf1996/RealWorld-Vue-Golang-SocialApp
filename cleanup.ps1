Write-Host "`n==================================" -ForegroundColor Red 
Write-Host "Complete Cleanup - Warning!" -ForegroundColor Red 
Write-Host "`n==================================" -ForegroundColor Red 

$confirm = Read-Host "Type 'DELETE' to continue (case-sensitive)"

if ($confirm -ne "DELETE") {
    Write-Host "`nCancelled - no changes made" -ForegroundColor Green 
    exit 
}

Write-Host "`n[1/6] Deleting social-app namespace resources..." -ForegroundColor Yellow 
kubectl delete all --all -n social-app 2>$null 
kubectl delete configmap --all -n social-app 2>$null 
kubectl delete secret --all -n social-app 2>$null 
kubectl delete pvc --all -n social-app 2>$null 
kubectl delete hpa --all -n social-app 2>$null 
kubectl delete networkpolicy --all -n social-app 2>$null 
kubectl delete ingress --all -n social-app 2>$null 

Write-Host "`n[2/6] Deleting social-app namespace..." -ForegroundColor Yellow 
kubectl delete namespace social-app 2>$null 

Write-Host "`n[3/6] wating for namespace deletion..." -ForegroundColor Yellow 

$timeout = 60
$elapsed = 0
while ((kubectl get namespace social-app 2>$null) -and ($elapsed -lt $timeout)) {
    Write-Host "Wating ... $elapsed seconds" -ForegroundColor Gray 
    Start-Sleep -Seconds 5 
    $elapsed += 5
}

Write-Host "`n[4/6] Deleting stray pods in default namespace..." -ForegroundColor Yellow 
kubectl delete pod redis -n default 2>$null 
kubectl delete pod --field-selector=status.phase--Completed -n default 2>$null

Write-Host "`n[5/6] Deleting orphaned persistent volumes..." -ForegroundColor Yellow 

$pvs = kubectl get pv -o json | ConvertFrom-Json
foreach ($pv in $pvs.items) {
    if ($pv.spec.claimRef.namespace -eq "social-app") {
        $pvName = $pv.metadata.name 
        Write-Host "Deleting PV: $pvName" -ForegroundColor Gray
        kubectl delete pv $pvName 2>$null 
    }
}

Write-Host "`n[6/6] Final Verfification..." -ForegroundColor Yellow 
Start-Sleep -Seconds 3


Write-Host "`n==================================" -ForegroundColor Green 
Write-Host "Cleanup is Completed!" -ForegroundColor Green 
Write-Host "`n==================================" -ForegroundColor Green 

Write-Host "`n namespaces:" -ForegroundColor Yellow
kubectl get namespaces

Write-Host "`n pods in defult:" -ForegroundColor Yellow
kubectl get pods -n default

Write-Host "`n Persisitent Volumes:" -ForegroundColor Yellow
kubectl get pv 


Write-Host "`n Social-app namespace (should be gone):" -ForegroundColor Yellow
kubectl get namespace social-app 2>$null

Write-Host "`n==================================" -ForegroundColor Green 
Write-Host "System is clean and ready for fresh deployment" -ForegroundColor Green 
Write-Host "`n==================================" -ForegroundColor Green 