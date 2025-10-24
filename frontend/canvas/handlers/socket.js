export function handleOnMessage(event, sceneRunner, userId) {
  const data = JSON.parse(event.data);

  switch (data.type) {
    case "snapshot":
      handleSnapshot(data, sceneRunner);
      break;
    case "updatedplayerlist":
      handleUpdatedPlayerList(data, sceneRunner, userId);
      break;
  }
}
export function handleSnapshot(data, sceneRunner) {
  sceneRunner.handlePlayerSnapshot(data);
}
export function handleDisconnect() {}
export function handleUpdatedPlayerList(data, sceneRunner, userId) {
  data.playerIds.forEach((id) => {
    if (id != userId) {
      sceneRunner.playerJoin({
        id: id,
      });
    }
  });
}
export function handleOnOpen() {}
export function handleOnClose() {}
