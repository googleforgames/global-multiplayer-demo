# Droidshooter

This is the Unreal Engine 5.1.0 project for the DroidShooter game!

## Local Builds

## Dedicated Game Server Docker Image

```shell
docker build . -t droidshooter-server
```

To run, you will need a copy of Agones [local SDK](https://agones.dev/site/docs/guides/client-sdks/local/) toolkit.

In one shell, start the local testing SDK server:

```shell
./sdk-server.linux.amd64 --local
```

Then in another shell:

```shell
docker run --rm --network=host droidshooter-server:latest
```

(In non *nix OS's where Docker is not native, you may need to do a slightly different network configuration to get the
same experience)

You now have a dedicated game server running on 127.0.0.1:7777

## Client building

See the top level [README.md](../README.md#game-client) game client section for full launch instructions.
