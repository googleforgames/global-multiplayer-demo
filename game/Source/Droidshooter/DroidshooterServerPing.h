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

#include <map>
#include <vector>
#include <algorithm>
#include "CoreMinimal.h"
#include "Icmp.h"

/**
 * 
 */
class DroidshooterServerPing
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

	void SetServersToValidate(uint16_t num);
	bool AllServersValidated();
	std::map<float, FString> GetPingedServers();
	void ClearPingedServers();

private:
	std::map<float, FString> pingResponses;
	uint32_t serversToValidate;
	static bool cmp(std::pair<float, FString>& a, std::pair<float, FString>& b);
};
