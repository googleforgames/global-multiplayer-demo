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
#include "Blueprint/UserWidget.h"
#include "DroidshooterPlayerHUD.generated.h"

UCLASS(Abstract)
class DROIDSHOOTER_API UDroidshooterPlayerHUD : public UUserWidget
{
	GENERATED_BODY()

public:
	//UDroidshooterPlayerHUD(const FObjectInitializer& ObjectInitializer);

public:
	/** Update HUD with current health. */
	void SetHealth(float CurrentHealth, float MaxHealth);

	/** Update HUD with current power. */
	void SetPower(float CurrentPower, float MaxPower);

	/** Update HUD with current score. */
	void SetScore(int32 score);

	/** Widget to display current health. */
	UPROPERTY(EditAnywhere, BlueprintReadWrite, meta = (BindWidget))
	class UProgressBar* HealthBar;

	/** Widget to use to display current power. */
	UPROPERTY(EditAnywhere, BlueprintReadWrite, meta = (BindWidget))
	class UProgressBar* PowerBar;

	/** Widget to display current score. */
	UPROPERTY(EditAnywhere, BlueprintReadWrite, meta = (BindWidget))
	class UTextBlock* ValueScore;

};