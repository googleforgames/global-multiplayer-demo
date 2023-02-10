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
#include "GameFramework/PlayerController.h"
#include "DroidshooterPlayerController.generated.h"

UCLASS()
class DROIDSHOOTER_API ADroidshooterPlayerController : public APlayerController
{
	GENERATED_BODY()

public:
	/** Setup input actions and context mappings for player. */
	virtual void SetupInputComponent() override;

	/** Mapping context used for pawn control. */
	UPROPERTY()
	class UInputMappingContext* PawnMappingContext;

	/** Action to update location. */
	UPROPERTY()
	class UInputAction* MoveAction;

	/** Action to update rotation. */
	UPROPERTY()
	class UInputAction* RotateAction;

	/** Action to toggle free fly mode. */
	UPROPERTY()
	class UInputAction* FreeFlyAction;

	/** Action to update spring arm length. */
	UPROPERTY()
	class UInputAction* SpringArmLengthAction;

	/** Action to start and stop shooting. */
	UPROPERTY()
	class UInputAction* ShootAction;
};