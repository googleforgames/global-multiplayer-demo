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
#include "GameFramework/PlayerState.h"
#include "DroidshooterPlayerState.generated.h"

/**
 * 
 */
UCLASS()
class DROIDSHOOTER_API ADroidshooterPlayerState : public APlayerState
{
	GENERATED_BODY()
	
public:
	ADroidshooterPlayerState();

protected:
	UPROPERTY(ReplicatedUsing = OnRep_TotalHits, BlueprintReadWrite)
	int TotalHits;

	UPROPERTY(ReplicatedUsing = OnRep_Health)
	float Health;

	/** Maximum amount of health to allow for player. */
	UPROPERTY(EditAnywhere)
	float MaxHealth;

	UFUNCTION()
	void OnRep_TotalHits();

	UFUNCTION()
	void OnRep_Health();

public:

	virtual void GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const override;
	void PlayerHit();
	float GetHealth();
	void UpdateHealth(float HealthDelta);
};

