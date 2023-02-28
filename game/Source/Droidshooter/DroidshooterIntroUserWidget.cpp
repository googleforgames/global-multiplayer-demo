// Copyright 2023 Google Inc. All Rights Reserved.Licensed under the Apache License, Version 2.0 (the "License");you may not use this file except in compliance with the License.You may obtain a copy of the License at    http://www.apache.org/licenses/LICENSE-2.0Unless required by applicable law or agreed to in writing, softwaredistributed under the License is distributed on an "AS IS" BASIS,WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.See the License for the specific language governing permissions andlimitations under the License.


#include "DroidshooterIntroUserWidget.h"
#include "HttpModule.h"
#include "Interfaces/IHttpRequest.h"
#include "Interfaces/IHttpResponse.h"
#include "Components/TextBlock.h"
#include "Components/EditableTextBox.h"
#include "Components/Button.h"
#include "Droidshooter.h"
#include "Json.h"

void UDroidshooterIntroUserWidget::NativePreConstruct() 
{
    Super::NativePreConstruct();
    // B_Auth->OnClicked.AddDynamic(this, &UDroidshooterIntroUserWidget::AuthenticateCall);
}

void UDroidshooterIntroUserWidget::FindPreferredGameServerLocation(const FString& accessToken) 
{
    // Query game server list
    // Ping each server
    // Create timer to check if all values have been updated
    // Save best region
    // Query FetchGameServer 
}

void UDroidshooterIntroUserWidget::FetchGameServer(const FString& accessToken)
{
    UE_LOG(LogDroidshooter, Log, TEXT("Fetch server/ip call with token: %s"), *accessToken);

    FString uriBase = "http://localhost:8081";
    FString uriPlay = uriBase + TEXT("/play");

    FHttpModule& httpModule = FHttpModule::Get();
    TSharedRef<IHttpRequest, ESPMode::ThreadSafe> pRequest = httpModule.CreateRequest();

    pRequest->SetHeader("Authorization", "Bearer " + accessToken);
    pRequest->SetVerb(TEXT("GET"));
    pRequest->SetURL(uriPlay);

    // Set the callback, which will execute when the HTTP call is complete
    pRequest->OnProcessRequestComplete().BindLambda(
        // Here, we "capture" the 'this' pointer (the "&"), so our lambda can call this
        // class's methods in the callback.
        [&](
            FHttpRequestPtr pRequest,
            FHttpResponsePtr pResponse,
            bool connectedSuccessfully) mutable {

                if (connectedSuccessfully) {
                    std::function<void(const TSharedPtr<FJsonObject>&)> f = [=](const TSharedPtr<FJsonObject>& JsonResponseObject) {
                        this->ProcessGameserverResponse(JsonResponseObject);
                    };
                    ProcessGenericJsonResponse(pResponse->GetContentAsString(), f);
                }
                else {
                    switch (pRequest->GetStatus()) {
                    case EHttpRequestStatus::Failed_ConnectionError:
                        UE_LOG(LogDroidshooter, Log, TEXT("Connection failed."));
                    default:
                        UE_LOG(LogDroidshooter, Log, TEXT("Request failed."));
                    }
                }
        });

    // Finally, submit the request for processing
    pRequest->ProcessRequest();
}

void UDroidshooterIntroUserWidget::AuthenticateCall(const FString& accessToken)
{
    UE_LOG(LogDroidshooter, Log, TEXT("Auth call with token: %s"), *accessToken);

	FString uriBase = "http://localhost:8081";
	FString uriProfile = uriBase + TEXT("/profile");

	FHttpModule& httpModule = FHttpModule::Get();

    // Create an http request
        // The request will execute asynchronously, and call us back on the Lambda below
    TSharedRef<IHttpRequest, ESPMode::ThreadSafe> pRequest = httpModule.CreateRequest();

    // This is where we set the HTTP method (GET, POST, etc)
    // The Space-Track.org REST API exposes a "POST" method we can use to get
    // general perturbations data about objects orbiting Earth
    pRequest->SetVerb(TEXT("GET"));

    // Here we set auth JWT token
    pRequest->SetHeader("Authorization", "Bearer " + accessToken);

    // We'll need to tell the server what type of content to expect in the POST data
    // pRequest->SetHeader(TEXT("Content-Type"), TEXT("application/x-www-form-urlencoded"));

    // FString RequestContent = TEXT("data=") + SomeValueVariable;
    // Set the POST content, which contains:
    // * Identity/password credentials for authentication
    // * A REST API query
    //   This allows us to login and get the results is a single call
    //   Otherwise, we'd need to manage session cookies across multiple calls.
    // pRequest->SetContentAsString(RequestContent);

    // Set the http URL
    pRequest->SetURL(uriProfile);

    // Set the callback, which will execute when the HTTP call is complete
    pRequest->OnProcessRequestComplete().BindLambda(
        // Here, we "capture" the 'this' pointer (the "&"), so our lambda can call this
        // class's methods in the callback.
        [&](
            FHttpRequestPtr pRequest,
            FHttpResponsePtr pResponse,
            bool connectedSuccessfully) mutable {

                if (connectedSuccessfully) {

                    // We should have a JSON response - attempt to process it.
                    std::function<void(const TSharedPtr<FJsonObject>&)> f = [=](const TSharedPtr<FJsonObject>& JsonResponseObject) {
                        this->ProcessProfileResponse(JsonResponseObject);
                    };
                        
                    ProcessGenericJsonResponse(pResponse->GetContentAsString(), f);
                }
                else {
                    switch (pRequest->GetStatus()) {
                    case EHttpRequestStatus::Failed_ConnectionError:
                        UE_LOG(LogDroidshooter, Log, TEXT("Connection failed."));
                    default:
                        UE_LOG(LogDroidshooter, Log, TEXT("Request failed."));
                    }
                }
        });

    // Finally, submit the request for processing
    pRequest->ProcessRequest();
}

void UDroidshooterIntroUserWidget::ProcessGenericJsonResponse(const FString& ResponseContent, std::function<void(const TSharedPtr<FJsonObject>&)>& func)
{
    TSharedRef<TJsonReader<TCHAR>> JsonReader = TJsonReaderFactory<TCHAR>::Create(ResponseContent);
    TSharedPtr<FJsonObject> JsonObject;

    if (FJsonSerializer::Deserialize(JsonReader, JsonObject))
    {
        func(JsonObject);
    }
}

void UDroidshooterIntroUserWidget::ProcessProfileResponse(const TSharedPtr<FJsonObject>& JsonResponseObject) 
{
    UE_LOG(LogDroidshooter, Log, TEXT("Processing json response for profile"));
    if (JsonResponseObject)
    {
        FString Name = JsonResponseObject->GetStringField(TEXT("player_name"));

        if (Name.Len() != 0) {
            //UE_LOG(LogDroidshooter, Log, TEXT("Logged in player is: %s"), "yes");
            UE_LOG(LogDroidshooter, Log, TEXT("Logged in player is: %s"), *Name);

            NameTextBlock->SetText(FText::FromString(Name));
        }
        else {
            UE_LOG(LogDroidshooter, Log, TEXT("Unable to fetch player data. Timed out token?"));
        }

    }
}

void UDroidshooterIntroUserWidget::ProcessGameserverResponse(const TSharedPtr<FJsonObject>& JsonResponseObject)
{
    if (JsonResponseObject)
    {
        FString IP = JsonResponseObject->GetStringField(TEXT("IP"));
        FString Port = JsonResponseObject->GetStringField(TEXT("Port"));

        if (IP.Len() != 0 && Port.Len() != 0) {
            // Set IP - Port variables!
            ServerIPValue = IP;
            ServerPortValue = Port;

            ServerIPBox->SetText(FText::FromString(IP));
            ServerPortBox->SetText(FText::FromString(Port));

            UE_LOG(LogDroidshooter, Log, TEXT("Found game server at: %s %s"), *IP, *Port);

        }
        else {
            UE_LOG(LogDroidshooter, Log, TEXT("Unable to server player data. Timed out token?"));
        }
    }
}
