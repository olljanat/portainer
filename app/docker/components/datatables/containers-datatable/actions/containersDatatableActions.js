angular.module('portainer.docker').component('containersDatatableActions', {
  templateUrl: 'app/docker/components/datatables/containers-datatable/actions/containersDatatableActions.html',
  controller: 'ContainersDatatableActionsController',
  bindings: {
    selectedItems: '=',
    selectedItemCount: '=',
    noStoppedItemsSelected: '=',
    noRunningItemsSelected: '=',
    noPausedItemsSelected: '=',
    showStartAction: '<',
    showStopAction: '<',
    showKillAction: '<',
    showRestartAction: '<',
    showPauseAction: '<',
    showResumeAction: '<',
    showRemoveAction: '<',
    showAddAction: '<'
  }
});
