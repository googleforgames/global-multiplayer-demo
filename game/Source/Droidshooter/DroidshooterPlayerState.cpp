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


#include "DroidshooterPlayerState.h"
#include "Droidshooter.h"
#include "Net/UnrealNetwork.h"
#include "DroidshooterGameStateBase.h"
#include "DroidshooterMainHUD.h"

ADroidshooterPlayerState::ADroidshooterPlayerState() {
	TotalHits = 0;

	// Health
	MaxHealth = 25.f;
	Health = MaxHealth;
}

void ADroidshooterPlayerState::OnRep_TotalHits() {
	if (ADroidshooterGameStateBase* GS = GetWorld()->GetGameState<ADroidshooterGameStateBase>()) {
		uint16 GSHits = GS->GetTotalHits();
		UE_LOG(LogDroidshooter, Log, TEXT("Player has TotalHits: %d/%d"), TotalHits, GSHits);

		ADroidshooterMainHUD* MyHUD = GetWorld()->GetFirstPlayerController()->GetHUD<ADroidshooterMainHUD>();
		if (MyHUD)
		{
			MyHUD->SetScore(TotalHits);
		}
	}
}

void ADroidshooterPlayerState::OnRep_Health()
{
	UE_LOG(LogDroidshooter, Log, TEXT("Player's health is %6.4lf"), Health);

	ADroidshooterMainHUD* MyHUD = GetWorld()->GetFirstPlayerController()->GetHUD<ADroidshooterMainHUD>();
	if (MyHUD)
	{
		MyHUD->SetHealth(Health, MaxHealth);
	}
}


void ADroidshooterPlayerState::PlayerHit()
{
	if (!HasAuthority())
		return;

	++TotalHits;
}

void ADroidshooterPlayerState::UpdateHealth(float HealthDelta)
{
	if (!HasAuthority())
		return;

	Health = FMath::Clamp(Health + HealthDelta, 0.f, MaxHealth);

	if (Health == 0.f)
	{
		UE_LOG(LogDroidshooter, Log, TEXT("Player is dead! (in DroidshooterPlayerState)"));
	}
}

float ADroidshooterPlayerState::GetHealth() {
	return Health;
}


void ADroidshooterPlayerState::GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const
{
	Super::GetLifetimeReplicatedProps(OutLifetimeProps);
	DOREPLIFETIME_CONDITION(ADroidshooterPlayerState, TotalHits, COND_OwnerOnly);
	DOREPLIFETIME_CONDITION(ADroidshooterPlayerState, Health, COND_OwnerOnly);
}
