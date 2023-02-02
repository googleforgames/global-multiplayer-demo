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

#pragma once

#include "CoreMinimal.h"
#include "GameFramework/Actor.h"
#include "DroidshooterShot.generated.h"

UCLASS()
class DROIDSHOOTER_API ADroidshooterShot : public AActor
{
	GENERATED_BODY()

public:
	ADroidshooterShot();

	/** Collision handling function. */
	UFUNCTION()
	void OnHit(UPrimitiveComponent* HitComponent, AActor* OtherActor, UPrimitiveComponent* OtherComponent,
		FVector NormalImpulse, const FHitResult& Hit);

	/** Sphere to use for root component and collisions. */
	UPROPERTY(EditAnywhere)
	class USphereComponent* Collision;

	/** Niagara FX component to hold system for flying visual. */
	UPROPERTY(EditAnywhere)
	class UNiagaraComponent* FlySystemComponent;

	/** Projectile component to move the actor. */
	UPROPERTY(EditAnywhere)
	class UProjectileMovementComponent* Movement;

	/** Niagara FX system for hit visual. */
	UPROPERTY(EditAnywhere)
	class UNiagaraSystem* HitSystem;

	/** How much to change health by when hitting another player. */
	UPROPERTY(EditAnywhere)
	float HealthDelta;

	/** How much to change power by when using this shot. */
	UPROPERTY(EditAnywhere)
	float PowerDelta;
};