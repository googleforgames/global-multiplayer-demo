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
#include "GameFramework/GameModeBase.h"
#include "AgonesComponent.h"
#include "DroidshooterGameMode.generated.h"

UCLASS()
class DROIDSHOOTER_API ADroidshooterGameMode : public AGameModeBase
{
	GENERATED_BODY()

public:
	ADroidshooterGameMode();
	virtual void InitGame(const FString& MapName, const FString& Options, FString& ErrorMessage) override;
	virtual void PreLogin(const FString& Options, const FString& Address, const FUniqueNetIdRepl& UniqueId, FString& ErrorMessage) override;
	virtual FString InitNewPlayer(APlayerController* NewPlayerController, const FUniqueNetIdRepl& UniqueId, const FString& Options, const FString& Portal = TEXT("")) override;
	virtual void Logout(AController* Exiting) override;
	void Respawn(AController* Controller);
	void PlayerHit();


	UFUNCTION(BlueprintCallable)
	void DumpStats(FString token, const FString gameId, const int kills, const int deaths);

	UPROPERTY(EditAnywhere, BlueprintReadWrite)
	UAgonesComponent* AgonesSDK;

	/** Frontend api endpoint */
	UPROPERTY(EditAnywhere, BlueprintReadWrite)
	FString FrontendApi;

	/** Frontend api endpoint access key - server only */
	UPROPERTY(EditAnywhere, BlueprintReadWrite)
	FString ApiKey;
private:
	TArray<class APlayerStart*> FreePlayerStarts;
};