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

#include "DroidshooterPlayerPawn.h"
#include "Droidshooter.h"
#include "DroidshooterPlayerController.h"
#include "DroidshooterGameMode.h"
#include "DroidshooterMainHUD.h"
#include "DroidshooterPlayerHUD.h"
#include "DroidshooterShot.h"
#include "Blueprint/UserWidget.h"
#include "Camera/CameraComponent.h"
#include "Components/StaticMeshComponent.h"
#include "EnhancedInputComponent.h"
#include "EnhancedInputSubsystems.h"
#include "GameFramework/FloatingPawnMovement.h"
#include "GameFramework/SpringArmComponent.h"
#include "Net/UnrealNetwork.h"

ADroidshooterPlayerPawn::ADroidshooterPlayerPawn()
{
	Collision = CreateDefaultSubobject<UStaticMeshComponent>(TEXT("Collision"));
	SetRootComponent(Collision);
	Collision->SetVisibleFlag(false);

	Body = CreateDefaultSubobject<UStaticMeshComponent>(TEXT("Body"));
	Body->SetupAttachment(Collision);

	Panels = CreateDefaultSubobject<UStaticMeshComponent>(TEXT("Panels"));
	Panels->SetupAttachment(Body);

	Thrusters = CreateDefaultSubobject<UStaticMeshComponent>(TEXT("Thrusters"));
	Thrusters->SetupAttachment(Body);

	// Springarm and Camera
	SpringArm = CreateDefaultSubobject<USpringArmComponent>(TEXT("SpringArm"));
	SpringArm->SetupAttachment(Collision);
	SpringArm->SetRelativeLocation(FVector(40.f, 0.f, 50.f));
	SpringArm->SetRelativeRotation(FRotator(-15.f, 0.f, 0.f));
	SpringArm->TargetArmLength = 600.f;

	SpringArmLengthScale = 2000.f;
	SpringArmLengthMin = 0.f;
	SpringArmLengthMax = 1000.f;

	Camera = CreateDefaultSubobject<UCameraComponent>(TEXT("Camera"));
	Camera->SetupAttachment(SpringArm, USpringArmComponent::SocketName);
	Camera->SetRelativeRotation(FRotator(15.f, 0.f, 0.f));
	Camera->PostProcessSettings.bOverride_MotionBlurAmount = true;
	Camera->PostProcessSettings.MotionBlurAmount = 0.1f;

	// Movement
	Movement = CreateDefaultSubobject<UFloatingPawnMovement>(TEXT("Movement"));
	Movement->MaxSpeed = 5000.f;
	Movement->Acceleration = 5000.f;
	Movement->Deceleration = 10000.f;

	MoveScale = 1.f;
	RotateScale = 50.f;
	bFreeFly = false;
	SpeedCheckInterval = 0.5f;
	SpeedCheckTranslationSum = FVector::ZeroVector;
	SpeedCheckTranslationCount = 0;
	SpeedCheckLastTranslation = FVector::ZeroVector;
	SpeedCheckLastTime = 0.f;
	MaxMovesWithHits = 30;
	MovesWithHits = 0;

	// Pawn Animation
	ZMovementFrequency = 2.f;
	ZMovementAmplitude = 5.f;
	ZMovementOffset = 0.f;

	TiltInput = 0.f;
	TiltMax = 15.f;
	TiltMoveScale = 0.6f;
	TiltRotateScale = 0.4f;
	TiltResetScale = 0.3f;

	// Shooting
	bShooting = false;
	ShootingInterval = 0.2f;
	ShootingOffset = FVector(200.f, 0.f, 0.f);
	ShotClass = ADroidshooterShot::StaticClass();
	ShootingLastTime = 0.f;

	// Power
	MaxPower = 25.f;
	Power = MaxPower;
	PowerRegenerateRate = 1.f;

	// Allow ticking for the pawn.
	PrimaryActorTick.bCanEverTick = true;

	// We should match this to DroidshooterShot (life span * speed) so they are destroyed at the same distance.
	NetCullDistanceSquared = 1600000000.f;

	// Force pawn to always spawn, even if it is colliding with another object.
	SpawnCollisionHandlingMethod = ESpawnActorCollisionHandlingMethod::AlwaysSpawn;

	MyHUD = nullptr;
}

void ADroidshooterPlayerPawn::GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const
{
	Super::GetLifetimeReplicatedProps(OutLifetimeProps);
	DOREPLIFETIME_CONDITION(ADroidshooterPlayerPawn, bShooting, COND_SimulatedOnly);
}

void ADroidshooterPlayerPawn::SetupPlayerInputComponent(UInputComponent* PlayerInputComponent)
{
	Super::SetupPlayerInputComponent(PlayerInputComponent);

	UEnhancedInputComponent* EIC = Cast<UEnhancedInputComponent>(PlayerInputComponent);
	ADroidshooterPlayerController* FPC = GetController<ADroidshooterPlayerController>();
	check(EIC && FPC);
	EIC->BindAction(FPC->MoveAction, ETriggerEvent::Triggered, this, &ADroidshooterPlayerPawn::Move);
	EIC->BindAction(FPC->RotateAction, ETriggerEvent::Triggered, this, &ADroidshooterPlayerPawn::Rotate);
	EIC->BindAction(FPC->FreeFlyAction, ETriggerEvent::Started, this, &ADroidshooterPlayerPawn::ToggleFreeFly);
	EIC->BindAction(FPC->SpringArmLengthAction, ETriggerEvent::Triggered, this,
		&ADroidshooterPlayerPawn::UpdateSpringArmLength);
	EIC->BindAction(FPC->ShootAction, ETriggerEvent::Started, this, &ADroidshooterPlayerPawn::Shoot);
	EIC->BindAction(FPC->ShootAction, ETriggerEvent::Completed, this, &ADroidshooterPlayerPawn::Shoot);

	ULocalPlayer* LocalPlayer = FPC->GetLocalPlayer();
	check(LocalPlayer);
	UEnhancedInputLocalPlayerSubsystem* Subsystem =
		LocalPlayer->GetSubsystem<UEnhancedInputLocalPlayerSubsystem>();
	check(Subsystem);
	Subsystem->ClearAllMappings();
	Subsystem->AddMappingContext(FPC->PawnMappingContext, 0);
}

void ADroidshooterPlayerPawn::BeginPlay()
{
	Super::BeginPlay();

	if (IsLocallyControlled())
	{
		ADroidshooterPlayerController* FPC = GetController<ADroidshooterPlayerController>();
		check(FPC);

		MyHUD = FPC->GetHUD<ADroidshooterMainHUD>();
	}
}

void ADroidshooterPlayerPawn::EndPlay(const EEndPlayReason::Type EndPlayReason)
{
	// For now we don't wanna destroy the HUD.
	/*if (MyHUD)
	{
		MyHUD->EndPlay(EndPlayReason);
		MyHUD = nullptr;
	}*/

	Super::EndPlay(EndPlayReason);
}

void ADroidshooterPlayerPawn::Tick(float DeltaSeconds)
{
	Super::Tick(DeltaSeconds);

	RegeneratePower();
	TryShooting();

	// Don't animate if we're the server.
	if (GetNetMode() != NM_DedicatedServer)
	{
		UpdatePawnAnimation();
	}

	// Replicate movement to server if we're the client controlling the pawn.
	if (GetLocalRole() == ROLE_AutonomousProxy)
	{
		UpdateServerTransform(Collision->GetRelativeTransform());
	}
}

/*
* Camera and Springarm
*/

void ADroidshooterPlayerPawn::UpdateSpringArmLength(const FInputActionValue& ActionValue)
{
	SpringArm->TargetArmLength += ActionValue[0] * GetWorld()->GetDeltaSeconds() * SpringArmLengthScale;
	SpringArm->TargetArmLength = FMath::Clamp(SpringArm->TargetArmLength,
		SpringArmLengthMin, SpringArmLengthMax);
}

/*
* Movement
*/

void ADroidshooterPlayerPawn::Move(const FInputActionValue& ActionValue)
{
	FVector Input = ActionValue.Get<FInputActionValue::Axis3D>();
	// UFloatingPawnMovement handles scaling this input based on the DeltaTime for this frame.
	AddMovementInput(GetActorRotation().RotateVector(Input), MoveScale);
	TiltInput += Input.Y * TiltMoveScale * MoveScale;
}

void ADroidshooterPlayerPawn::Rotate(const FInputActionValue& ActionValue)
{
	FRotator Input(ActionValue[0], ActionValue[1], ActionValue[2]);
	Input *= GetWorld()->GetDeltaSeconds() * RotateScale;
	TiltInput += Input.Yaw * TiltRotateScale;

	if (bFreeFly) {
		AddActorLocalRotation(Input);
	}
	else {
		Input += GetActorRotation();
		Input.Pitch = FMath::ClampAngle(Input.Pitch, -89.9f, 89.9f);
		Input.Roll = 0;
		SetActorRotation(Input);
	}
}

void ADroidshooterPlayerPawn::ToggleFreeFly()
{
	bFreeFly = !bFreeFly;
}

void ADroidshooterPlayerPawn::UpdateServerTransform_Implementation(FTransform Transform)
{
	// Make sure the client does not try to move faster than the game allows. We can't check
	// on each update using the server delta time since the server may tick at different rates
	// than the client, and the server might process multiple updates in one tick. Instead, we
	// calculate an average position every SpeedCheckInterval and check the speed using that.
	float Now = GetWorld()->GetRealTimeSeconds();
	if (SpeedCheckLastTime == 0)
	{
		SpeedCheckLastTranslation = Transform.GetTranslation();
		SpeedCheckLastTime = Now;
		SpeedCheckTranslationSum = FVector::ZeroVector;
		SpeedCheckTranslationCount = 0;
	}
	else
	{
		SpeedCheckTranslationSum += Transform.GetTranslation();
		SpeedCheckTranslationCount++;

		if (Now - SpeedCheckLastTime > SpeedCheckInterval)
		{
			FVector SpeedCheckTranslation = SpeedCheckTranslationSum / SpeedCheckTranslationCount;
			float Distance = FVector::Distance(SpeedCheckLastTranslation, SpeedCheckTranslation);
			float Speed = Distance / (Now - SpeedCheckLastTime);
			//UE_LOG(LogDroidshooter, Log, TEXT("Client speed update: %s %.3f %d"), *Controller->GetName(), Speed, SpeedCheckCount);

			SpeedCheckLastTime = Now;
			SpeedCheckTranslationSum = FVector::ZeroVector;
			SpeedCheckTranslationCount = 0;

			// Allow 10% more than MaxSpeed to account for time and translation variation.
			if (Speed > Movement->MaxSpeed * 1.1f)
			{
				// Moving too fast, ignore update and move client back to last translation.
				UE_LOG(LogDroidshooter, Log, TEXT("Player moving too fast: %s %.3f"), *Controller->GetName(), Speed);
				UpdateClientTransform(FTransform(Collision->GetRelativeRotation(), SpeedCheckLastTranslation));
				return;
			}

			SpeedCheckLastTranslation = SpeedCheckTranslation;
		}
	}

	// Move client with a sweep to see if we hit anything. We seem to get hits on the server even when
	// the client sends valid moves, especially while sliding against objects. We'll always have a valid
	// move on the server since the sweep will correct the server side. We expect the client to eventually
	// send a transform that moves cleanly, but if we go too long (MaxMovesWithHits), send a correction
	// back to the client. This will cause a stutter on the client so we want to keep it minimal.
	FTransform OldTransform = Collision->GetRelativeTransform();
	FHitResult HitResult;
	Collision->SetRelativeTransform(Transform, true, &HitResult);
	if (HitResult.bBlockingHit) {
		//float ExpectedDistance = FVector::Distance(HitResult.TraceStart, HitResult.TraceEnd);
		//UE_LOG(LogDroidshooter, Log, TEXT("Player hit object: %s %d (%.3f-%.3f=%.3f)"), *Controller->GetName(),
		//	MovesWithHits, ExpectedDistance, HitResult.Distance, ExpectedDistance - HitResult.Distance);
		MovesWithHits++;
	}
	else {
		//UE_LOG(LogDroidshooter, Log, TEXT("Player move ok: %s"), *Controller->GetName());
		MovesWithHits = 0;
	}

	if (MovesWithHits > MaxMovesWithHits) {
		UE_LOG(LogDroidshooter, Log, TEXT("Correcting player transform: %s"), *Controller->GetName());
		UpdateClientTransform(Collision->GetRelativeTransform());
	}
}

void ADroidshooterPlayerPawn::UpdateClientTransform_Implementation(FTransform Transform)
{
	Collision->SetRelativeTransform(Transform);
}

/*
* Pawn Animation
*/

void ADroidshooterPlayerPawn::UpdatePawnAnimation()
{
	// Add Z Movement.
	if (ZMovementAmplitude)
	{
		float ZMovement = FMath::Sin(GetWorld()->GetTimeSeconds() * ZMovementFrequency) * ZMovementAmplitude;
		Body->SetRelativeLocation(FVector(0.f, 0.f, ZMovement + ZMovementOffset));
	}

	// Add body and head tilting.
	FRotator Rotation = Body->GetRelativeRotation();

	if (TiltInput != 0.f)
	{
		Rotation.Roll = FMath::Clamp(Rotation.Roll + TiltInput, -TiltMax, TiltMax);
		TiltInput = 0.f;
	}

	// Always try to tilt back towards the center.
	if (Rotation.Roll > 0.f)
	{
		Rotation.Roll -= TiltResetScale;
		if (Rotation.Roll < 0.f)
			Rotation.Roll = 0.f;
	}
	else if (Rotation.Roll < 0.f)
	{
		Rotation.Roll += TiltResetScale;
		if (Rotation.Roll > 0.f)
			Rotation.Roll = 0.f;
	}

	Body->SetRelativeRotation(Rotation);
	//Thrusters->SetRelativeRotation(FRotator(0.f, Rotation.Roll, 0.f));
}

/*
* Shooting
*/

void ADroidshooterPlayerPawn::Shoot(const FInputActionValue& ActionValue)
{
	bShooting = ActionValue[0] > 0.f;
	UpdateServerShooting(bShooting);
}

void ADroidshooterPlayerPawn::UpdateServerShooting_Implementation(bool bNewShooting)
{
	bShooting = bNewShooting;
}

void ADroidshooterPlayerPawn::TryShooting()
{
	float Now = GetWorld()->GetRealTimeSeconds();
	float PowerDelta = Cast<ADroidshooterShot>(ShotClass->GetDefaultObject())->PowerDelta;

	// We spawn shot actors independently on the server and all clients. This way we only need to replicate
	// the shooting state changes, and not each spawned shot actor and related movement updates.
	if (!bShooting || Now - ShootingLastTime < ShootingInterval || Power + PowerDelta <= 0)
	{
		return;
	}

	FRotator ShotRotation = Body->GetComponentRotation();
	FVector ShotStart = Body->GetComponentLocation() + ShotRotation.RotateVector(ShootingOffset);
	FActorSpawnParameters ActorSpawnParams;
	ActorSpawnParams.Owner = this;
	ADroidshooterShot* Shot = GetWorld()->SpawnActor<ADroidshooterShot>(ShotClass, ShotStart, ShotRotation, ActorSpawnParams);
	if (Shot)
	{
		ShootingLastTime = Now;
		Shot->SetInstigator(this);

		// Consume used power for shot and update HUD power bar.
		Power += PowerDelta;
		if (MyHUD)
		{
			MyHUD->SetPower(Power, MaxPower);
		}

		UE_LOG(LogDroidshooter, Log, TEXT("Shot spawned %s %s %s"), *GetName(),
			IsNetMode(NM_Client) ? TEXT("Client") : TEXT("Server"),
			Controller ? TEXT("Controlled") : TEXT("Simulated"));
	}
}


/*
* Power
*/

void ADroidshooterPlayerPawn::RegeneratePower()
{
	Power = FMath::Clamp(Power + (PowerRegenerateRate * GetWorld()->GetDeltaSeconds()), 0.f, MaxPower);
	if (MyHUD)
	{
		MyHUD->SetPower(Power, MaxPower);
	}
}
