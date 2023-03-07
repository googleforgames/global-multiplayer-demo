// Copyright 2023 Google Inc. All Rights Reserved.Licensed under the Apache License, Version 2.0 (the "License");you may not use this file except in compliance with the License.You may obtain a copy of the License at    http://www.apache.org/licenses/LICENSE-2.0Unless required by applicable law or agreed to in writing, softwaredistributed under the License is distributed on an "AS IS" BASIS,WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.See the License for the specific language governing permissions andlimitations under the License.

#include "DroidshooterServerPing.h"
#include "Droidshooter.h"

DroidshooterServerPing::DroidshooterServerPing(): serversToValidate(0)
{
}

DroidshooterServerPing::~DroidshooterServerPing()
{
}

void DroidshooterServerPing::CheckIfServerIsOnline(FString ServerPublicIP, FString ServerPort)
{
	/* First bind our OnServerCheckFinished function to PingResult.
	* When we get any reply from UDP server PingResult will be called
	* and when PingResult is called the binded method (OnServerCheckFinished) is also called. */
	PingResult.BindRaw(this, &DroidshooterServerPing::OnServerCheckFinished);

	/* Our UDP server public ip and port we have to ping.
	* Port should be the exact same port we defined on UDP server node.js file on EC2 server */
	const FString Address = FString::Printf(TEXT("%s:%s"), *ServerPublicIP, *ServerPort);

	// Finally just ping.
	FUDPPing::UDPEcho(Address, 5.f, PingResult);
}

void DroidshooterServerPing::OnServerCheckFinished(FIcmpEchoResult Result)
{
	// Unbind the function. Its no longer required.
	PingResult.Unbind();

	// Simply set a status.
	FString PingStatus = "Ping Failed";

	// Do your stuff based on the result.
	switch (Result.Status)
	{
	case EIcmpResponseStatus::Success:
		PingStatus = "Success";
		break;
	case EIcmpResponseStatus::Timeout:
		PingStatus = "Timeout";
		break;
	case EIcmpResponseStatus::Unreachable:
		PingStatus = "Unreachable";
		break;
	case EIcmpResponseStatus::Unresolvable:
		PingStatus = "Unresolvable";
		break;
	case EIcmpResponseStatus::InternalError:
		PingStatus = "Internal Error";
		break;
	default:
		PingStatus = "Unknown Error";
		break;
	}

	if (Result.Status == EIcmpResponseStatus::Success) {
		pingResponses.insert({Result.Time, Result.ResolvedAddress });
	}
	
	serversToValidate--;

	// Simple log
	UE_LOG(LogDroidshooter, Log, TEXT("Ping status: %s @ %s in %.2fms"), *PingStatus, *Result.ResolvedAddress, Result.Time);
}

void DroidshooterServerPing::SetServersToValidate(uint16_t num) {
	serversToValidate = num;
}

bool DroidshooterServerPing::AllServersValidated() {
	return serversToValidate == 0;
}

std::map<float, FString> DroidshooterServerPing::GetPingedServers() {
	// Declare vector of pairs
	std::vector<std::pair<float, FString>> tmp_pair;

	for (auto& it : pingResponses) {
		tmp_pair.push_back(it);
	}

	std::sort(tmp_pair.begin(), tmp_pair.end(), DroidshooterServerPing::cmp);

	pingResponses.clear();

	for (auto& it : tmp_pair) {
		pingResponses.insert({ it.first, it.second });
	}

	return pingResponses;
}

void DroidshooterServerPing::ClearPingedServers() {
	pingResponses.clear();
}

bool DroidshooterServerPing::cmp(std::pair<float, FString>& a, std::pair<float, FString>& b)
{
	return a.first < b.first;
}
