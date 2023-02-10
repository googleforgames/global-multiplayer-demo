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

#include "DroidshooterMainHUD.h"

ADroidshooterMainHUD::ADroidshooterMainHUD()
{
}

void ADroidshooterMainHUD::DrawHUD()
{
	Super::DrawHUD();
}

void ADroidshooterMainHUD::BeginPlay()
{
	Super::BeginPlay();

	if(PlayerHUDWidgetClass) 
	{
		PlayerHUDWidget = CreateWidget<UDroidshooterPlayerHUD>(GetWorld(), PlayerHUDWidgetClass);
		if (PlayerHUDWidget)
		{
			PlayerHUDWidget->AddToViewport();
			PlayerHUDWidget->SetHealth(25.f, 25.f);
			PlayerHUDWidget->SetPower(25.f, 25.f);
			PlayerHUDWidget->SetScore(0);
		}
	}
}

void ADroidshooterMainHUD::EndPlay(const EEndPlayReason::Type EndPlayReason)
{
	if(PlayerHUDWidget) 
	{
		PlayerHUDWidget->RemoveFromParent();
		// We can't destroy the widget directly, let the GC take care of it.
		PlayerHUDWidget = nullptr;
	}
	Super::EndPlay(EndPlayReason);

}

void ADroidshooterMainHUD::Tick(float DeltaSeconds)
{
	Super::Tick(DeltaSeconds);
}

void ADroidshooterMainHUD::SetHealth(float CurrentHealth, float MaxHealth)
{
	if (PlayerHUDWidget)
	{
		PlayerHUDWidget->SetHealth(CurrentHealth, MaxHealth);
	}
}

void ADroidshooterMainHUD::SetPower(float CurrentPower, float MaxPower)
{
	if (PlayerHUDWidget)
	{
		PlayerHUDWidget->SetPower(CurrentPower, MaxPower);
	}
}

void ADroidshooterMainHUD::SetScore(int32 score)
{
	if (PlayerHUDWidget)
	{
		PlayerHUDWidget->SetScore(score);
	}
}