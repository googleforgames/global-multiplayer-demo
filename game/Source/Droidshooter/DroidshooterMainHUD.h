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
#include "GameFramework/HUD.h"
#include "DroidshooterPlayerHUD.h"
#include "Components/WidgetComponent.h"
#include "DroidshooterMainHUD.generated.h"

/**
 * 
 */
UCLASS()
class DROIDSHOOTER_API ADroidshooterMainHUD : public AHUD
{
	GENERATED_BODY()

public:
	ADroidshooterMainHUD();

	/*
	* HUD
	*/

	virtual void DrawHUD() override;

	virtual void BeginPlay() override;

	virtual void EndPlay(const EEndPlayReason::Type EndPlayReason) override;

	virtual void Tick(float DeltaSeconds) override;

	/** Update HUD with current health. */
	UFUNCTION()
	void SetHealth(float CurrentHealth, float MaxHealth);

	/** Update HUD with current power. */
	UFUNCTION()
	void SetPower(float CurrentPower, float MaxPower);

	/** Update HUD with current score. */
	UFUNCTION()
	void SetScore(int32 score);

	/** Widget class to spawn for the heads up display. */
	UPROPERTY(EditAnywhere, Category = "Widgets")
	TSubclassOf<class UUserWidget> PlayerHUDWidgetClass;

private:
	UDroidshooterPlayerHUD* PlayerHUDWidget;
};
