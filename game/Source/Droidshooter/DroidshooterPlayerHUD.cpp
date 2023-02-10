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

#include "DroidshooterPlayerHUD.h"
#include "Components/ProgressBar.h"
#include "Components/TextBlock.h"

/*UDroidshooterPlayerHUD::UDroidshooterPlayerHUD(const FObjectInitializer& ObjectInitializer)
{
}*/


void UDroidshooterPlayerHUD::SetHealth(float CurrentHealth, float MaxHealth)
{
	if (HealthBar)
	{
		HealthBar->SetPercent(CurrentHealth / MaxHealth);
	}
}

void UDroidshooterPlayerHUD::SetPower(float CurrentPower, float MaxPower)
{
	if (PowerBar)
	{
		PowerBar->SetPercent(CurrentPower / MaxPower);
	}
}

void UDroidshooterPlayerHUD::SetScore(int32 score)
{
	if (ValueScore) 
	{
		ValueScore->SetText(FText::AsNumber(score));
	}
}