// Copyright 2023 Google Inc. All Rights Reserved.Licensed under the Apache License, Version 2.0 (the "License");you may not use this file except in compliance with the License.You may obtain a copy of the License at    http://www.apache.org/licenses/LICENSE-2.0Unless required by applicable law or agreed to in writing, softwaredistributed under the License is distributed on an "AS IS" BASIS,WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.See the License for the specific language governing permissions andlimitations under the License.

#pragma once

#include <map>
#include <vector>
#include "CoreMinimal.h"
#include "Icmp.h"

/**
 * 
 */
class DROIDSHOOTER_API DroidshooterServerPing
{
public:
	DroidshooterServerPing();
	~DroidshooterServerPing();

	/* Delegate to bind our custom function when we receive a reply from our EC2 UDP server */
	FIcmpEchoResultDelegate PingResult;

	/* Sends a UDP echo to the given server and waits for a reply */
	UFUNCTION(BlueprintCallable)
	void CheckIfServerIsOnline(FString ServerPublicIP, FString ServerPort);

	/* Delegate called when we get a reply from our EC2 UDP server */
	void OnServerCheckFinished(FIcmpEchoResult Result);

private:
	void QueuePing(FString ip);
	void DequeuePing(FString ip);

	std::map<float, FString> pingResponses;
	std::vector<FString> pingQueue;
};
