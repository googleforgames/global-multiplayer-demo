// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#pragma once

#include "CoreMinimal.h"
#include "TimerManager.h" 
#include "Blueprint/UserWidget.h"
#include "DroidshooterServerPing.h"
#include "DroidshooterIntroUserWidget.generated.h"

/**
 * 
 */
UCLASS()
class DROIDSHOOTER_API UDroidshooterIntroUserWidget : public UUserWidget
{
	GENERATED_BODY()

public:

	void NativePreConstruct();

	UFUNCTION(BlueprintCallable)
	void AuthenticateCall(const FString& frontendApi, const FString& accessToken);

	UFUNCTION(BlueprintCallable)
	void FetchGameServer(const FString& frontendApi, const FString& accessToken, const FString preferredRegion, const FString ping);

	UFUNCTION(BlueprintCallable)
	void FindPreferredGameServerLocation(const FString& frontendApi, const FString& accessToken);

	void ProcessProfileResponse(const FString& ResponseContent);
	void ProcessGameserverResponse(const FString& ResponseContent);
	void ProcessServersToPingResponse(const FString& ResponseContent);
	void ProcessGenericJsonResponse(const FString& ResponseContent, std::function<void(const TSharedPtr<FJsonObject>&)>& func);

	void AllServersValidated();

	// Server IP/Port editboxes
	UPROPERTY(EditAnywhere, BlueprintReadWrite, meta = (BindWidget))
	class UEditableTextBox* ServerIPBox;
	UPROPERTY(EditAnywhere, BlueprintReadWrite, meta = (BindWidget))
	class UEditableTextBox* ServerPortBox;

	/** Widget to display current user's name. */
	UPROPERTY(EditAnywhere, BlueprintReadWrite, meta = (BindWidget))
	class UTextBlock* NameTextBlock;

	/** Saving token for further queries */
	UPROPERTY(EditAnywhere, BlueprintReadWrite)
	FString GlobalAccessToken;

	/** Saving token for further queries */
	UPROPERTY(EditAnywhere, BlueprintReadWrite)
	FString FrontendApi;

	/** Saving token for further queries */
	UPROPERTY(EditAnywhere, BlueprintReadWrite)
	FString ServerIPValue;

	/** Saving token for further queries */
	UPROPERTY(EditAnywhere, BlueprintReadWrite)
	FString ServerPortValue;

private:

	FTimerHandle MemberTimerHandle;
	DroidshooterServerPing ServerPinger;

};
