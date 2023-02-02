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
#include "GameFramework/Actor.h"
#include "DroidshooterMapRoom.generated.h"

UCLASS()
class DROIDSHOOTER_API ADroidshooterMapRoom : public AActor
{
	GENERATED_BODY()

public:
	ADroidshooterMapRoom();

#if WITH_EDITOR
	/** Check properties that change to see if we need to rebuild. */
	virtual void PostEditChangeProperty(FPropertyChangedEvent& PropertyChangedEvent) override;
#endif

	/** Build or rebuild the room if needed. */
	virtual void OnConstruction(const FTransform& Transform) override;

	/** Static mesh to use for walls. */
	UPROPERTY(EditAnywhere)
	class UInstancedStaticMeshComponent* Walls;

	/** Static mesh to use for edges where two walls meet. */
	UPROPERTY(EditAnywhere)
	class UInstancedStaticMeshComponent* Edges;

	/** Static mesh to use for corners where three edges meet. */
	UPROPERTY(EditAnywhere)
	class UInstancedStaticMeshComponent* Corners;

	/** Static mesh to use for walls where tubes exit room. */
	UPROPERTY(EditAnywhere)
	class UInstancedStaticMeshComponent* TubeWalls;

	/** Static mesh to use for tubes extending from tube walls. */
	UPROPERTY(EditAnywhere)
	class UInstancedStaticMeshComponent* Tubes;

	/** Size of grid to use when placing meshes. */
	UPROPERTY(EditAnywhere, meta = (ClampMin = 0, RebuildMapRoom))
	float GridSize;

	/** Size of the room to build. This will round up to the next odd number. */
	UPROPERTY(EditAnywhere, meta = (ClampMin = 1, ClampMax = 25, RebuildMapRoom))
	uint32 RoomSize;

	/** How thick the walls are, used for box collision alignment. */
	UPROPERTY(EditAnywhere, meta = (ClampMin = 0, RebuildMapRoom))
	float WallThickness;

	/** How much to offset the edge collision by. */
	UPROPERTY(EditAnywhere, meta = (RebuildMapRoom))
	float EdgeCollisionOffset;

	/** How many faces to use while building the tube collision boxes. */
	UPROPERTY(EditAnywhere, meta = (ClampMin = 3, ClampMax = 32, RebuildMapRoom))
	uint32 TubeCollisionFaces;

	/** Radius from center of tube to center of collision boxes. */
	UPROPERTY(EditAnywhere, meta = (ClampMin = 0, RebuildMapRoom))
	float TubeCollisionRadius;

	/** How thick the collision boxes are. */
	UPROPERTY(EditAnywhere, meta = (ClampMin = 0, RebuildMapRoom))
	float TubeCollisionThickness;

	/** How many tubes to extend off the center positive X wall, if any (0 to disable). */
	UPROPERTY(EditAnywhere, meta = (ClampMax = 1000, RebuildMapRoom))
	uint32 PositiveXTubeSize;

	/** How many tubes to extend off the center negative X wall, if any (0 to disable). */
	UPROPERTY(EditAnywhere, meta = (ClampMax = 1000, RebuildMapRoom))
	uint32 NegativeXTubeSize;

	/** How many tubes to extend off the center positive Y wall, if any (0 to disable). */
	UPROPERTY(EditAnywhere, meta = (ClampMax = 1000, RebuildMapRoom))
	uint32 PositiveYTubeSize;

	/** How many tubes to extend off the center negative Y wall, if any (0 to disable). */
	UPROPERTY(EditAnywhere, meta = (ClampMax = 1000, RebuildMapRoom))
	uint32 NegativeYTubeSize;

	/** How many tubes to extend off the center positive Z wall, if any (0 to disable). */
	UPROPERTY(EditAnywhere, meta = (ClampMax = 1000, RebuildMapRoom))
	uint32 PositiveZTubeSize;

	/** How many tubes to extend off the center negative Z wall, if any (0 to disable). */
	UPROPERTY(EditAnywhere, meta = (ClampMax = 1000, RebuildMapRoom))
	uint32 NegativeZTubeSize;

private:
	/** Whether we need to rebuild or not. */
	int32 bRebuild:1;

	/** Distance from the center of the room to walls. */
	int32 WallOffset;

	/** Add section of tubes with lights and collision boxes. */
	void AddTubeInstances(uint32 TubeSize, const FRotator& Rotation);

	/** Helper function to add new components. */
	template<class T>
	T* AddComponent(const FTransform& Transform);

	/** Helper function to add collision boxes. */
	void AddCollisionBox(const FVector& Extent, const FRotator& Rotation,
		const FVector& Translation, const FRotator& FaceRotation = FRotator::ZeroRotator);

	/** Helper function to add point lights. */
	void AddPointLight(float Intensity, float Radius,
		const FRotator& Rotation, const FVector& Translation);
};