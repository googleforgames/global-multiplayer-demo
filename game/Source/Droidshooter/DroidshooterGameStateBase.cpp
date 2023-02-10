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


#include "DroidshooterGameStateBase.h"
#include "Droidshooter.h"
#include "Net/UnrealNetwork.h"

ADroidshooterGameStateBase::ADroidshooterGameStateBase()
{
	TotalHits = 0;
}

void ADroidshooterGameStateBase::GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const
{
	Super::GetLifetimeReplicatedProps(OutLifetimeProps);
	DOREPLIFETIME_CONDITION(ADroidshooterGameStateBase, TotalHits, COND_SimulatedOnly);
}

void ADroidshooterGameStateBase::OnRep_TotalHits() {
	UE_LOG(LogDroidshooter, Log, TEXT("Client: TotalHits: %d"), TotalHits);
}


void ADroidshooterGameStateBase::PlayerHit()
{
	if (HasAuthority()) 
	{
		UE_LOG(LogDroidshooter, Log, TEXT("Player was hit (in DroidshooterGameState)"));
		++TotalHits;
		UE_LOG(LogDroidshooter, Log, TEXT("Server: TotalHits: %d"), TotalHits);
	}
}

uint16 ADroidshooterGameStateBase::GetTotalHits()
{
	return TotalHits;
}