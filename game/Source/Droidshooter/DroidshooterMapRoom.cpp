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

#include "DroidshooterMapRoom.h"
#include "Droidshooter.h"
#include "Components/BoxComponent.h"
#include "Components/InstancedStaticMeshComponent.h"
#include "Components/PointLightComponent.h"
#include "Components/SceneComponent.h"

ADroidshooterMapRoom::ADroidshooterMapRoom()
{
	PrimaryActorTick.bCanEverTick = false;
	GridSize = 1000.f;
	RoomSize = 3;
	WallThickness = 50.f;
	EdgeCollisionOffset = -250.f;
	TubeCollisionFaces = 8;
	TubeCollisionRadius = 425.f;
	TubeCollisionThickness = 200.f;
	bRebuild = true;

	USceneComponent* SceneComponent = CreateDefaultSubobject<USceneComponent>(TEXT("SceneComponent"));
	SetRootComponent(SceneComponent);

	Walls = CreateDefaultSubobject<UInstancedStaticMeshComponent>(TEXT("Walls"));
	Walls->SetupAttachment(SceneComponent);

	Edges = CreateDefaultSubobject<UInstancedStaticMeshComponent>(TEXT("Edges"));
	Edges->SetupAttachment(SceneComponent);

	Corners = CreateDefaultSubobject<UInstancedStaticMeshComponent>(TEXT("Corners"));
	Corners->SetupAttachment(SceneComponent);

	TubeWalls = CreateDefaultSubobject<UInstancedStaticMeshComponent>(TEXT("TubeWalls"));
	TubeWalls->SetupAttachment(SceneComponent);

	Tubes = CreateDefaultSubobject<UInstancedStaticMeshComponent>(TEXT("Tubes"));
	Tubes->SetupAttachment(SceneComponent);
}

/** Helper function to keep calls in OnConstruction concise. */
static FORCEINLINE void AddInstance(UInstancedStaticMeshComponent* Component,
	const FRotator& Rotation, const FVector& Translation)
{
	Component->AddInstance(FTransform(Rotation, Rotation.RotateVector(Translation)));
}

/**
 * Rotations used while adding mesh instances. These assume the mesh
 * has been created with a base orientation of positive X.
 */
static const FRotator PositiveX(0.f, 0.f, 0.f);
static const FRotator PositiveX90(0.f, 0.f, 90.f);
static const FRotator PositiveX180(0.f, 0.f, 180.f);
static const FRotator PositiveX270(0.f, 0.f, 270.f);
static const FRotator NegativeX(0.f, 180.f, 0.f);
static const FRotator NegativeX90(0.f, 180.f, 90.f);
static const FRotator NegativeX180(0.f, 180.f, 180.f);
static const FRotator NegativeX270(0.f, 180.f, 270.f);

static const FRotator PositiveY(0.f, 90.f, 0.f);
static const FRotator PositiveY180(0.f, 90.f, 180.f);
static const FRotator NegativeY(0.f, 270.f, 0.f);
static const FRotator NegativeY180(0.f, 270.f, 180.f);

static const FRotator PositiveZ(90.f, 0.f, 0.f);
static const FRotator NegativeZ(-90.f, 0.f, 0.f);

static const FRotator PositivePitch45(45.f, 0.f, 0.f);
static const FRotator NegativePitch45(-45.f, 0.f, 0.f);
static const FRotator PositiveYaw45(0.f, 45.f, 0.f);
static const FRotator NegativeYaw45(0.f, -45.f, 0.f);

#if WITH_EDITOR
void ADroidshooterMapRoom::PostEditChangeProperty(FPropertyChangedEvent& PropertyChangedEvent)
{
	if (PropertyChangedEvent.Property && PropertyChangedEvent.Property->HasMetaData("RebuildMapRoom"))
		bRebuild = true;

	Super::PostEditChangeProperty(PropertyChangedEvent);
}
#endif

void ADroidshooterMapRoom::OnConstruction(const FTransform& Transform)
{
	Super::OnConstruction(Transform);

	// Only rebuild if needed.
	if (!bRebuild)
		return;

	bRebuild = false;

	UE_LOG(LogDroidshooter, Log,
		TEXT("ADroidshooterMapRoom::OnConstruction Building Room Size %d (this=%x)"),
		RoomSize, this);

	Walls->ClearInstances();
	Edges->ClearInstances();
	Corners->ClearInstances();
	TubeWalls->ClearInstances();
	Tubes->ClearInstances();

	TArray<UPointLightComponent*> Lights;
	GetComponents<UPointLightComponent>(Lights);
	for (UPointLightComponent* Light : Lights)
		Light->DestroyComponent();

	TArray<UBoxComponent*> Boxes;
	GetComponents<UBoxComponent>(Boxes);
	for (UBoxComponent* Box : Boxes)
		Box->DestroyComponent();

	// Implicit floor with integer division, which makes all room sizes end up being odd.
	int32 HalfSize = RoomSize / 2;
	WallOffset = (HalfSize + 1) * GridSize;
	FVector Translation(WallOffset, 0.f, 0.f);

	for (int32 a = -HalfSize; a <= HalfSize; a++)
	{
		Translation.Y = GridSize * a;

		for (int32 b = -HalfSize; b <= HalfSize; b++)
		{
			Translation.Z = GridSize * b;

			// Build walls, placing a tube wall in the center if needed.
			auto WallType = [&](uint32 TubeSize)
			{
				return (a == 0 && b == 0 && TubeSize > 0) ? TubeWalls : Walls;
			};

			AddInstance(WallType(PositiveXTubeSize), PositiveX, Translation);
			AddInstance(WallType(NegativeXTubeSize), NegativeX, Translation);
			AddInstance(WallType(PositiveYTubeSize), PositiveY, Translation);
			AddInstance(WallType(NegativeYTubeSize), NegativeY, Translation);
			AddInstance(WallType(PositiveZTubeSize), PositiveZ, Translation);
			AddInstance(WallType(NegativeZTubeSize), NegativeZ, Translation);
		}

		// Build edges.
		Translation.Z = WallOffset;
		AddInstance(Edges, PositiveX, Translation);
		AddInstance(Edges, PositiveX90, Translation);
		AddInstance(Edges, PositiveX180, Translation);
		AddInstance(Edges, PositiveX270, Translation);
		AddInstance(Edges, NegativeX, Translation);
		AddInstance(Edges, NegativeX90, Translation);
		AddInstance(Edges, NegativeX180, Translation);
		AddInstance(Edges, NegativeX270, Translation);
		AddInstance(Edges, PositiveY, Translation);
		AddInstance(Edges, PositiveY180, Translation);
		AddInstance(Edges, NegativeY, Translation);
		AddInstance(Edges, NegativeY180, Translation);
	}

	// Build corners.
	Translation.Y = WallOffset;
	AddInstance(Corners, PositiveX, Translation);
	AddInstance(Corners, PositiveX90, Translation);
	AddInstance(Corners, PositiveX180, Translation);
	AddInstance(Corners, PositiveX270, Translation);
	AddInstance(Corners, NegativeX, Translation);
	AddInstance(Corners, NegativeX90, Translation);
	AddInstance(Corners, NegativeX180, Translation);
	AddInstance(Corners, NegativeX270, Translation);

	// Build tubes and add wall and tube collisions.
	AddTubeInstances(PositiveXTubeSize, PositiveX);
	AddTubeInstances(NegativeXTubeSize, NegativeX);
	AddTubeInstances(PositiveYTubeSize, PositiveY);
	AddTubeInstances(NegativeYTubeSize, NegativeY);
	AddTubeInstances(PositiveZTubeSize, PositiveZ);
	AddTubeInstances(NegativeZTubeSize, NegativeZ);

	// Build edge collision boxes.
	Translation.X = WallOffset + EdgeCollisionOffset;
	Translation.Y = 0;
	Translation.Z = Translation.X;
	FVector Extent(WallThickness / 2.f, WallOffset, GridSize / 2.f);
	AddCollisionBox(Extent, PositiveX, Translation, PositivePitch45);
	AddCollisionBox(Extent, PositiveX90, Translation, PositiveYaw45);
	AddCollisionBox(Extent, PositiveX180, Translation, NegativePitch45);
	AddCollisionBox(Extent, PositiveX270, Translation, NegativeYaw45);
	AddCollisionBox(Extent, NegativeX, Translation, PositivePitch45);
	AddCollisionBox(Extent, NegativeX90, Translation, PositiveYaw45);
	AddCollisionBox(Extent, NegativeX180, Translation, NegativePitch45);
	AddCollisionBox(Extent, NegativeX270, Translation, NegativeYaw45);
	AddCollisionBox(Extent, PositiveY, Translation, PositivePitch45);
	AddCollisionBox(Extent, PositiveY180, Translation, NegativePitch45);
	AddCollisionBox(Extent, NegativeY, Translation, PositivePitch45);
	AddCollisionBox(Extent, NegativeY180, Translation, NegativePitch45);

	// Add large light in center of room.
	AddPointLight(2.f, WallOffset * 2, PositiveX, FVector::ZeroVector);
}

void ADroidshooterMapRoom::AddTubeInstances(uint32 TubeSize, const FRotator& Rotation)
{
	FVector Translation(WallOffset, 0, 0);
	FVector Extent(WallThickness / 2.f, WallOffset, WallOffset);

	if (TubeSize == 0)
	{
		// No tubes, so add one collision for the entire wall.
		AddCollisionBox(Extent, Rotation, Translation);
		return;
	}

	// Setup wall collision with four sections, leaving a hole in the middle for the tube.
	Translation.Y = (WallOffset / 2.f) - (GridSize / 4.f);
	Translation.Z = (WallOffset / 2.f) + (GridSize / 4.f);
	Extent.Y = Translation.Z;
	Extent.Z = Translation.Y;

	AddCollisionBox(Extent, Rotation, Translation);
	AddCollisionBox(Extent, Rotation + PositiveX90, Translation);
	AddCollisionBox(Extent, Rotation + PositiveX180, Translation);
	AddCollisionBox(Extent, Rotation + PositiveX270, Translation);

	// Add light at tube entrance.
	Translation.Y = 0;
	Translation.Z = 0;
	AddPointLight(1.f, GridSize, Rotation, Translation);

	// Start at 1 because the first tube is the tube wall added in OnConstruction.
	for (uint32 a = 1; a < TubeSize; a++)
	{
		Translation.X = WallOffset + GridSize * a;
		AddInstance(Tubes, Rotation, Translation);
		AddPointLight(1.f, GridSize, Rotation, Translation);
	}

	// Add tube collision boxes.
	Extent.X = ((TubeSize * GridSize) / 2.f) - ((GridSize - WallThickness) / 4.f);
	Extent.Y = GridSize / 2.f;
	Extent.Z = TubeCollisionThickness / 2.f;
	Translation.X = WallOffset - (WallThickness / 2.f) + Extent.X;

	for (uint32 a = 0; a < TubeCollisionFaces; a++)
	{
		float BoxAngle = (float(a) / TubeCollisionFaces) * (PI * 2.f);
		Translation.Y = FMath::Sin(BoxAngle) * TubeCollisionRadius;
		Translation.Z = FMath::Cos(BoxAngle) * TubeCollisionRadius;

		AddCollisionBox(Extent, Rotation, Translation,
			FRotator(0.f, 0.f, BoxAngle * (180.f / PI)));
	}
}

template<class T>
T* ADroidshooterMapRoom::AddComponent(const FTransform& Transform)
{
	T* Component = NewObject<T>(this);
	Component->AttachToComponent(RootComponent, FAttachmentTransformRules::KeepRelativeTransform);
	Component->RegisterComponent();
	AddInstanceComponent(Component);
	Component->SetRelativeTransform(Transform);
	return Component;
}

void ADroidshooterMapRoom::AddCollisionBox(const FVector& Extent, const FRotator& Rotation,
	const FVector& Translation, const FRotator& FaceRotation)
{
	UBoxComponent* Box = AddComponent<UBoxComponent>(
		FTransform(Rotation + FaceRotation, Rotation.RotateVector(Translation)));
	Box->SetCollisionProfileName(TEXT("BlockAll"));
	Box->SetBoxExtent(Extent);
}

void ADroidshooterMapRoom::AddPointLight(float Intensity, float Radius,
	const FRotator& Rotation, const FVector& Translation)
{
	UPointLightComponent* Light = AddComponent<UPointLightComponent>(
		FTransform(Rotation, Rotation.RotateVector(Translation)));
	Light->Intensity = Intensity;
	Light->SetAttenuationRadius(Radius);
	Light->SetSoftSourceRadius(Radius);
	Light->SetCastShadows(false);
	Light->bUseInverseSquaredFalloff = false;
	Light->LightFalloffExponent = 1;
}