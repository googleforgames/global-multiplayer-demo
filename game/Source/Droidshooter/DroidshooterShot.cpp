// Copyright 2023 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#include "DroidshooterShot.h"
#include "Droidshooter.h"
#include "DroidshooterPlayerPawn.h"
#include "Components/SphereComponent.h"
#include "GameFramework/ProjectileMovementComponent.h"
#include "NiagaraComponent.h"
#include "NiagaraFunctionLibrary.h"
#include "DroidshooterGameMode.h"
#include "DroidshooterPlayerState.h"


ADroidshooterShot::ADroidshooterShot()
{
	Collision = CreateDefaultSubobject<USphereComponent>(TEXT("Collision"));
	SetRootComponent(Collision);
	Collision->SetSphereRadius(25.f);
	Collision->OnComponentHit.AddDynamic(this, &ADroidshooterShot::OnHit);
	Collision->SetCollisionProfileName(TEXT("BlockAllDynamic"));

	FlySystemComponent = CreateDefaultSubobject<UNiagaraComponent>(TEXT("Niagara"));
	FlySystemComponent->SetupAttachment(Collision);

	Movement = CreateDefaultSubobject<UProjectileMovementComponent>(TEXT("Movement"));
	Movement->InitialSpeed = 20000.f;
	Movement->MaxSpeed = 20000.f;
	Movement->ProjectileGravityScale = 0.f;

	// Destroy after moving 40k units (life span * speed) to match the net cull distance
	// in the player pawn.
	InitialLifeSpan = 2.f;
	HealthDelta = -2.f;
	PowerDelta = -1.f;
}

void ADroidshooterShot::OnHit(UPrimitiveComponent* HitComponent, AActor* OtherActor, UPrimitiveComponent* OtherComponent,
	FVector NormalImpulse, const FHitResult& Hit)
{
	UE_LOG(LogDroidshooter, Log, TEXT("Shot hit %s %s"), *OtherActor->GetName(),
		IsNetMode(NM_Client) ? TEXT("Client") : TEXT("Server"));

	ADroidshooterPlayerPawn* Target = Cast<ADroidshooterPlayerPawn>(OtherActor);
	ADroidshooterPlayerPawn* Shooter = GetInstigator<ADroidshooterPlayerPawn>();

	if (Target && Target != Shooter)
	{
		if (Target->GetLocalRole() == ROLE_Authority) {
			ADroidshooterGameMode* GM = GetWorld()->GetAuthGameMode< ADroidshooterGameMode>();
			ADroidshooterPlayerState* PSShooter = Cast< ADroidshooterPlayerState>(Shooter->GetPlayerState());
			ADroidshooterPlayerState* PSTarget = Cast< ADroidshooterPlayerState>(Target->GetPlayerState());

			if (PSTarget) {
				PSTarget->UpdateHealth(-2.f);
			}

			
			if (PSTarget->GetHealth() == 0.f && GM)
			{
				// Total server kills
				GM->PlayerHit();

				// Kill and respawn player
				AController* Controller = Target->GetController();
				Target->Destroy();
				GM->Respawn(Controller);

				if (PSShooter) {
					// Increase shooters score
					PSShooter->PlayerHit(); 
				}
			}
		}
	}

	if (HitSystem)
	{
		UNiagaraFunctionLibrary::SpawnSystemAtLocation(GetWorld(), HitSystem,
			Collision->GetComponentLocation(), Collision->GetComponentRotation());
	}

	Destroy();
}