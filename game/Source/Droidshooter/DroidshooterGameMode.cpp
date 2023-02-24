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

#include "DroidshooterGameMode.h"
#include "DroidshooterPlayerPawn.h"
#include "DroidshooterGameStateBase.h"
#include "DroidshooterPlayerState.h"
#include "Droidshooter.h"
#include "EngineUtils.h"
#include "GameFramework/PlayerStart.h"
#include "AgonesComponent.h"
#include "Classes.h"
#include <random>

ADroidshooterGameMode::ADroidshooterGameMode()
{
	// Causes the editor to hang. Loading classname during runtime is much better, check respawn function.
	/*static ConstructorHelpers::FClassFinder<APawn> PlayerPawnBlueprint(TEXT("/Game/Player/DS_PlayerPawnBP.DS_PlayerPawnBP_C"));
	if (PlayerPawnBlueprint.Class != NULL) {
		DefaultPawnClass = PlayerPawnBlueprint.Class;
	}*/

	AgonesSDK = CreateDefaultSubobject<UAgonesComponent>(TEXT("AgonesSDK"));
}

void ADroidshooterGameMode::InitGame(const FString& MapName, const FString& Options, FString& ErrorMessage)
{
    Super::InitGame(MapName, Options, ErrorMessage);
    UE_LOG(LogDroidshooter, Log, TEXT("Game is running: %s %s"), *MapName, *Options);

	if (GetWorld()->IsNetMode(NM_DedicatedServer)) {
		UE_LOG(LogDroidshooter, Log, TEXT("Server Started for map: %s"), *MapName);
	}

	for (TActorIterator<APlayerStart> It(GetWorld()); It; ++It)
	{
		FreePlayerStarts.Add(*It);
		UE_LOG(LogDroidshooter, Log, TEXT("Found player start: %s"), *(*It)->GetName());
	}
}

void ADroidshooterGameMode::PreLogin(const FString& Options, const FString& Address, const FUniqueNetIdRepl& UniqueId, FString& ErrorMessage)
{
	if (FreePlayerStarts.Num() == 0)
	{
		ErrorMessage = TEXT("Server full");
	}

	Super::PreLogin(Options, Address, UniqueId, ErrorMessage);
}

FString ADroidshooterGameMode::InitNewPlayer(APlayerController* NewPlayerController, const FUniqueNetIdRepl& UniqueId, const FString& Options, const FString& Portal)
{
	if (FreePlayerStarts.Num() == 0)
	{
		UE_LOG(LogDroidshooter, Log, TEXT("No free player starts in InitNewPlayer"));
		return FString(TEXT("No free player starts"));
	}

	NewPlayerController->StartSpot = FreePlayerStarts.Pop();
	UE_LOG(LogDroidshooter, Log, TEXT("Using player start %s for %s"),
		*NewPlayerController->StartSpot->GetName(), *NewPlayerController->GetName());
	return Super::InitNewPlayer(NewPlayerController, UniqueId, Options, Portal);
}

void ADroidshooterGameMode::Logout(AController* Exiting)
{
	UE_LOG(LogDroidshooter, Log, TEXT("Player is disconnecting!"));
	Super::Logout(Exiting);

	for (TActorIterator<APlayerStart> It(GetWorld()); It; ++It)
	{
		if (FreePlayerStarts.Contains(*It)) {
			UE_LOG(LogDroidshooter, Log, TEXT("Playerstart in Iterator: %s"), *(*It)->GetName());
		}
		else {
			UE_LOG(LogDroidshooter, Log, TEXT("Playerstart NOT in Iterator: %s. Adding again."), *(*It)->GetName());
			FreePlayerStarts.Add(*It);
		}
	}

}

void ADroidshooterGameMode::Respawn(AController* Controller)
{
	UE_LOG(LogDroidshooter, Log, TEXT("Respawning!"));

	if (Controller) {
		if (GetLocalRole() == ROLE_Authority)
		{
			std::random_device rd; // obtain a random number from hardware
			std::mt19937 gen(rd()); // seed the generator
			std::uniform_int_distribution<> distr(-10000, 10000); // define the range

			FString TheClassPath = "Class'/Game/Player/DS_PlayerPawnBP.DS_PlayerPawnBP_C'";
			const TCHAR* TheClass = *TheClassPath;
			UClass* PlayerPawnBlueprintClass = LoadObject<UClass>(nullptr, TheClass);

			if (PlayerPawnBlueprintClass == NULL)
				return;

			FVector Location = FVector(distr(gen), distr(gen), 0);
			if (ADroidshooterPlayerPawn* Pawn = GetWorld()->SpawnActor<ADroidshooterPlayerPawn>(PlayerPawnBlueprintClass, Location, FRotator::ZeroRotator)) {
				Controller->Possess(Pawn);
				ADroidshooterPlayerState* PlayerState = Cast< ADroidshooterPlayerState>(Pawn->GetPlayerState());

				// Reset health back to normal
				PlayerState->UpdateHealth(25.f);
			}

		}
	}
}

void ADroidshooterGameMode::PlayerHit() {
	if (ADroidshooterGameStateBase* GS = GetGameState<ADroidshooterGameStateBase>()) {
		UE_LOG(LogDroidshooter, Log, TEXT("Player was hit (in DroidshooterGameMode)"));
		GS->PlayerHit();
	}
}