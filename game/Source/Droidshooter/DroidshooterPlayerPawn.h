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
#include "GameFramework/Pawn.h"
#include "DroidshooterMainHUD.h"
#include "DroidshooterPlayerPawn.generated.h"

UCLASS()
class DROIDSHOOTER_API ADroidshooterPlayerPawn : public APawn
{
	GENERATED_BODY()

public:

	ADroidshooterPlayerPawn();

	/** Setup properties that should be replicated from the server to clients. */
	virtual void GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const override;

	/** Bind input actions from player controller. */
	virtual void SetupPlayerInputComponent(class UInputComponent* PlayerInputComponent) override;

	/** Setup pawn when game starts. */
	virtual void BeginPlay() override;

	/** Cleanup pawn when game ends. */
	virtual void EndPlay(const EEndPlayReason::Type EndPlayReason) override;

	/** Perform pawn updates that need to happen every frame. */
	virtual void Tick(float DeltaSeconds) override;

private:

	/** Static mesh to use for root component and collisions. */
	UPROPERTY(EditAnywhere)
	class UStaticMeshComponent* Collision;

	/** Static mesh for the body, attached to the root. */
	UPROPERTY(EditAnywhere)
	class UStaticMeshComponent* Body;

	/** Static mesh for the panels, attached to the body. */
	UPROPERTY(EditAnywhere)
	class UStaticMeshComponent* Panels;

	/** Static mesh for the thrusters, attached to the body. */
	UPROPERTY(EditAnywhere)
	class UStaticMeshComponent* Thrusters;

	/*
	* Springarm and Camera
	*/

	/** Spring arm to hold camera, attached to the root. */
	UPROPERTY(EditAnywhere)
	class USpringArmComponent* SpringArm;

	/** Camera attached to spring arm to provide pawn's view. */
	UPROPERTY(EditAnywhere)
	class UCameraComponent* Camera;

	/** Scale to apply to spring arm length input. */
	UPROPERTY(EditAnywhere)
	float SpringArmLengthScale;

	/** Minimum spring arm length. */
	UPROPERTY(EditAnywhere)
	float SpringArmLengthMin;

	/** Maximum spring arm length. */
	UPROPERTY(EditAnywhere)
	float SpringArmLengthMax;

	/** Handle input to update spring arm length. */
	void UpdateSpringArmLength(const struct FInputActionValue& ActionValue);

	/*
	* Movement
	*/

	/** Use floating pawn movement to smooth out motion. */
	UPROPERTY(EditAnywhere)
	class UFloatingPawnMovement* Movement;

	/** Scale to apply to location input. */
	UPROPERTY(EditAnywhere)
	float MoveScale;

	/** Scale to apply to rotation input. */
	UPROPERTY(EditAnywhere)
	float RotateScale;

	/** Whether to use free flying mode. Caution: might cause motion sickness! */
	UPROPERTY(EditAnywhere)
	bool bFreeFly;

	/** Handle input to update location. */
	void Move(const struct FInputActionValue& ActionValue);

	/** Handle input to update rotation. */
	void Rotate(const struct FInputActionValue& ActionValue);

	/** Handle input to toggle free flying. */
	void ToggleFreeFly();

	/** Update server with latest transform from the client. */
	UFUNCTION(Server, Unreliable)
	void UpdateServerTransform(FTransform Transform);

	/** Update client transform when server needs to send a correction. */
	UFUNCTION(Client, Unreliable)
	void UpdateClientTransform(FTransform Transform);

	/** How often to check speed with average translation for each interval. */
	UPROPERTY(EditAnywhere)
	float SpeedCheckInterval;

	/** Current sum of translations used for average in next speed check. */
	FVector SpeedCheckTranslationSum;

	/** Current count of translations used for average in next speed check. */
	uint32 SpeedCheckTranslationCount;

	/** Last average translation used in speed check. */
	FVector SpeedCheckLastTranslation;

	/** Last time speed was checked. */
	float SpeedCheckLastTime;

	/** Max number of consecutive moves with hits to allow from client. */
	UPROPERTY(EditAnywhere)
	uint32 MaxMovesWithHits;

	/** How many consecutive moves with hits we've seen from client. */
	uint32 MovesWithHits;

	/*
	* Pawn Animation
	*/

	/** Amplitude to scale Z movement by, 0 to disable Z movement. */
	UPROPERTY(EditAnywhere)
	float ZMovementAmplitude;

	/** Frequency to adjust Z movement by. */
	UPROPERTY(EditAnywhere)
	float ZMovementFrequency;

	/** Offset to apply to Z movement after after calculating amplitude. */
	UPROPERTY(EditAnywhere)
	float ZMovementOffset;

	/** The current input to apply to tilt. */
	float TiltInput;

	/** Max tilt to apply to body and head while turning. */
	UPROPERTY(EditAnywhere)
	float TiltMax;

	/** Scale to apply to movement input for tilt. */
	UPROPERTY(EditAnywhere)
	float TiltMoveScale;

	/** Scale to apply to rotation input for tilt. */
	UPROPERTY(EditAnywhere)
	float TiltRotateScale;

	/** Scale to apply to tilt when resetting to center. */
	UPROPERTY(EditAnywhere)
	float TiltResetScale;

	/** Perform tilting and hovering animation for pawn. */
	void UpdatePawnAnimation();

	/*
	* Shooting
	*/

	/** Whether we are currently shooting. */
	UPROPERTY(Replicated)
	bool bShooting;

	/** How often we can shoot. */
	UPROPERTY(EditAnywhere)
	float ShootingInterval;

	/** Where to spawn the ShotClass relative to the Body. */
	UPROPERTY(EditAnywhere)
	FVector ShootingOffset;

	/** Class to spawn when shooting. */
	UPROPERTY(EditAnywhere)
	TSubclassOf<class ADroidshooterShot> ShotClass;

	/** Handle input to start and stop shooting. */
	void Shoot(const struct FInputActionValue& ActionValue);

	/** Update server with latest shooting state from the client. */
	UFUNCTION(Server, Reliable)
	void UpdateServerShooting(bool bNewShooting);

	/** Try spawning a shot actor moving in the direction of where the camera is looking. */
	void TryShooting();

	/** Last time we shot. */
	float ShootingLastTime;

	/*
	* HUD
	*/

	ADroidshooterMainHUD* MyHUD;

private:

	/*
	* Power
	*/

	/** Maximum amount of power to allow for player. */
	UPROPERTY(EditAnywhere)
	float MaxPower;

	/** Current power of player. */
	float Power;

	/** How much power to regenerate every second. */
	UPROPERTY(EditAnywhere)
	float PowerRegenerateRate;

	/** Regenerate power during tick. */
	void RegeneratePower();
};